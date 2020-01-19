package main

import (
	"fmt"
	"github.com/hybridgroup/mjpeg"
	"gocv.io/x/gocv"
	"log"
	"net/http"
)

var (
	deviceID int
	err      error
	webcam   *gocv.VideoCapture
	stream   *mjpeg.Stream
)

func main() {
	// open webcam
	webcam, err = gocv.VideoCaptureDevice(0)
	if err != nil {
		fmt.Printf("error opening video capture device: %v\n", deviceID)
		return
	}
	defer webcam.Close()

	// create the mjpeg stream
	stream = mjpeg.NewStream()

	// start capturing
	go capture()

	fmt.Println("Capturing....")

	// start http server
	http.Handle("/", stream)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func capture() {
	img := gocv.NewMat()
	defer img.Close()

	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("cannot read device %d\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		buf, _ := gocv.IMEncode(".jpg", img)
		stream.UpdateJPEG(buf)
	}
}
