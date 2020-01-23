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
)

var (
	deviceID int
	err      error
	cam      *gocv.VideoCapture
	stream   *mjpeg.Stream
)

const MinimumArea = 3000

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// open cam
	cam, err = gocv.VideoCaptureDevice(0)
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

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		os.Getenv("user"): os.Getenv("password"),
	}))

	authorized.GET("/", func(c *gin.Context) {
		stream.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	_ = r.Run(":8080")
}

func capture() {
	window := gocv.NewWindow("Motion Window")
	defer window.Close()
	log.Info("Capturing....1")
	img := gocv.NewMat()
	defer img.Close()
	log.Info("Capturing....2")
	imgDelta := gocv.NewMat()
	defer imgDelta.Close()
	log.Info("Capturing....3")
	imgThresh := gocv.NewMat()
	defer imgThresh.Close()
	log.Info("Capturing....4")
	mog2 := gocv.NewBackgroundSubtractorMOG2()
	defer mog2.Close()
	log.Info("Capturing....5")
	status := "Ready"
	log.Info("Capturing....6")
	fmt.Printf("Start reading device: %v\n", deviceID)
	for {
		log.Info("Capturing....7")
		if ok := cam.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}
		log.Info("Capturing....8")
		status = "Ready"
		statusColor := color.RGBA{0, 255, 0, 0}

		// first phase of cleaning up image, obtain foreground only
		mog2.Apply(img, &imgDelta)
		log.Info("Capturing....9")
		// remaining cleanup of the image to use for finding contours.
		// first use threshold
		gocv.Threshold(imgDelta, &imgThresh, 25, 255, gocv.ThresholdBinary)
		log.Info("Capturing....10")
		// then dilate
		kernel := gocv.GetStructuringElement(gocv.MorphRect, image.Pt(3, 3))
		defer kernel.Close()
		gocv.Dilate(imgThresh, &imgThresh, kernel)
		log.Info("Capturing....11")
		// now find contours
		contours := gocv.FindContours(imgThresh, gocv.RetrievalExternal, gocv.ChainApproxSimple)
		for i, c := range contours {
			area := gocv.ContourArea(c)
			if area < MinimumArea {
				continue
			}
			log.Info("Capturing....12")
			status = "Motion detected"
			statusColor = color.RGBA{255, 0, 0, 0}
			gocv.DrawContours(&img, contours, i, statusColor, 2)

			rect := gocv.BoundingRect(c)
			gocv.Rectangle(&img, rect, color.RGBA{0, 0, 255, 0}, 2)
		}
		log.Info("Capturing....13")
		gocv.PutText(&img, status, image.Pt(10, 20), gocv.FontHersheyPlain, 1.2, statusColor, 2)
		log.Info("Capturing....")
		buf, _ := gocv.IMEncode(".jpg", img)
		stream.UpdateJPEG(buf)
	}
}
