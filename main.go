package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hybridgroup/mjpeg"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"gocv.io/x/gocv"
	"image"
	"image/color"
	"os"
	"raspi-gocv/vault"
)

var (
	deviceID int
	err      error
	cam      *gocv.VideoCapture
	stream   *mjpeg.Stream
)

var (
	username string
	password string
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	deviceID := os.Getenv("deviceID")
	// open cam
	cam, err = gocv.OpenVideoCapture(deviceID)
	if err != nil {
		log.Info("error opening video capture device: %v\n", deviceID)
		return
	}
	defer cam.Close()

	// create the mjpeg stream
	stream = mjpeg.NewStream()

	// start capturing
	go capture()

	log.Info("Capturing....")

	username, password = vault.ReadSecret("secret/data/demo/userinfo")

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		username: password,
	}))

	authorized.GET("/", func(c *gin.Context) {
		stream.ServeHTTP(c.Writer, c.Request)
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	_ = r.Run(":8080")
}

func capture() {
	xmlFile := "data/fullbody_recognition_model.xml"
	// prepare image matrix
	img := gocv.NewMat()
	defer img.Close()

	// color for the rect when faces detected
	blue := color.RGBA{0, 0, 255, 0}

	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load(xmlFile) {
		log.Warn("Error reading cascade file: %v\n", xmlFile)
		return
	}

	log.Info("Start reading device: %v\n", deviceID)
	for {
		if ok := cam.Read(&img); !ok {
			log.Info("Device closed: %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// detect faces
		rects := classifier.DetectMultiScale(img)

		// draw a rectangle around each face on the original image,
		// along with text identifying as "Human"
		for _, r := range rects {
			gocv.Rectangle(&img, r, blue, 3)

			size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
			pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
			gocv.PutText(&img, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)
		}

		buf, _ := gocv.IMEncode(".jpg", img)
		stream.UpdateJPEG(buf)
	}
}
