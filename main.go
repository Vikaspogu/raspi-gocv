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
	"path/filepath"
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
	model := "data/res10_300x300_ssd_iter_140000.caffemodel"
	config := "data/deploy.prototxt"
	backend := gocv.NetBackendDefault
	target := gocv.NetTargetCPU
	img := gocv.NewMat()
	defer img.Close()

	// open DNN object tracking model
	net := gocv.ReadNet(model, config)
	if net.Empty() {
		log.Info("Error reading network model from : %v %v\n", model, config)
		return
	}
	defer net.Close()
	_ = net.SetPreferableBackend(backend)
	_ = net.SetPreferableTarget(target)

	var ratio float64
	var mean gocv.Scalar
	var swapRGB bool

	if filepath.Ext(model) == ".caffemodel" {
		ratio = 1.0
		mean = gocv.NewScalar(104, 177, 123, 0)
		swapRGB = false
	} else {
		ratio = 1.0 / 127.5
		mean = gocv.NewScalar(127.5, 127.5, 127.5, 0)
		swapRGB = true
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

		// convert image Mat to 300x300 blob that the object detector can analyze
		blob := gocv.BlobFromImage(img, ratio, image.Pt(300, 300), mean, swapRGB, false)

		// feed the blob into the detector
		net.SetInput(blob, "")

		// run a forward pass thru the network
		prob := net.Forward("")

		performDetection(&img, prob)

		prob.Close()
		blob.Close()

		buf, _ := gocv.IMEncode(".jpg", img)
		stream.UpdateJPEG(buf)
	}
}

// performDetection analyzes the results from the detector network,
// which produces an output blob with a shape 1x1xNx7
// where N is the number of detections, and each detection
// is a vector of float values
// [batchId, classId, confidence, left, top, right, bottom]
func performDetection(frame *gocv.Mat, results gocv.Mat) {
	for i := 0; i < results.Total(); i += 7 {
		confidence := results.GetFloatAt(0, i+2)
		if confidence > 0.5 {
			left := int(results.GetFloatAt(0, i+3) * float32(frame.Cols()))
			top := int(results.GetFloatAt(0, i+4) * float32(frame.Rows()))
			right := int(results.GetFloatAt(0, i+5) * float32(frame.Cols()))
			bottom := int(results.GetFloatAt(0, i+6) * float32(frame.Rows()))
			gocv.Rectangle(frame, image.Rect(left, top, right, bottom), color.RGBA{0, 255, 0, 0}, 2)
		}
	}
}
