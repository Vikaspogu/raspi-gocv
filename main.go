package main

import (
	"fmt"
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

const (MinimumArea = 3000)

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
		fmt.Printf("error opening video capture device: %v\n", deviceID)
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
	img := gocv.NewMat()
	defer img.Close()

	imgDelta := gocv.NewMat()
	defer imgDelta.Close()

	imgThresh := gocv.NewMat()
	defer imgThresh.Close()

	mog2 := gocv.NewBackgroundSubtractorMOG2()
	defer mog2.Close()

	status := "Ready"

	fmt.Printf("Start reading device: %v\n", deviceID)
	for {
		if ok := cam.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		status = "Ready"
		statusColor := color.RGBA{0, 255, 0, 0}

		// first phase of cleaning up image, obtain foreground only
		mog2.Apply(img, &imgDelta)

		// remaining cleanup of the image to use for finding contours.
		// first use threshold
		gocv.Threshold(imgDelta, &imgThresh, 25, 255, gocv.ThresholdBinary)

		// then dilate
		kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
		defer kernel.Close()
		gocv.Dilate(imgThresh, &imgThresh, kernel)

		// now find contours
		contours := gocv.FindContours(imgThresh, gocv.RetrievalExternal, gocv.ChainApproxSimple)
		for i, c := range contours {
			area := gocv.ContourArea(c)
			if area < MinimumArea {
				continue
			}

			status = "Motion detected"
			statusColor = color.RGBA{255, 0, 0, 0}
			gocv.DrawContours(&img, contours, i, statusColor, 2)

			rect := gocv.BoundingRect(c)
			gocv.Rectangle(&img, rect, color.RGBA{0, 0, 255, 0}, 2)
		}

		gocv.PutText(&img, status, image.Pt(10, 20), gocv.FontHersheyPlain, 1.2, statusColor, 2)

		buf, _ := gocv.IMEncode(".jpg", img)
		stream.UpdateJPEG(buf)
	}
}
