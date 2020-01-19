package main

import (
	"fmt"
	"github.com/hybridgroup/mjpeg"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gocv.io/x/gocv"
	"image/color"
	"log"
	"net/http"
)

var (
	deviceID int
	err      error
	cam   *gocv.VideoCapture
	stream   *mjpeg.Stream
)

func main() {
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

	fmt.Println("Capturing....")

	// start http server
	http.Handle("/", stream)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func capture() {
	img := gocv.NewMat()
	defer img.Close()

	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	classifier.Load("./data/haarcascade_frontalface_alt2.xml")

	defer classifier.Close()

	for {
		if ok := cam.Read(&img); !ok {
			fmt.Printf("cannot read device %d\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// color for the rect when faces detected
		boxColor := color.RGBA{0, 255, 0, 0}

		// detect faces
		rects := classifier.DetectMultiScale(img)

		// draw a rectangle around each face on the original image,
		// along with text identifying as "Human"
		for _, r := range rects {
			gocv.Rectangle(&img, r, boxColor, 3)
		}

		buf, _ := gocv.IMEncode(".jpg", img)
		stream.UpdateJPEG(buf)
	}
}
