package main

import (
	"fmt"
	"github.com/hybridgroup/mjpeg"
	"gocv.io/x/gocv"
	"image"
	"image/color"
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

	xmlFile := "./data/haarcascade_frontalface_default.xml"
	// load classifier to recognize faces
	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	classifier.Load(xmlFile)

	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("cannot read device %d\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// detect faces
		rects := classifier.DetectMultiScale(img)

		// color for the rect when faces detected
		blue := color.RGBA{0, 0, 255, 0}

		// draw a rectangle around each face on the original image,
		// along with text identifing as "Human"
		for _, r := range rects {
			gocv.Rectangle(&img, r, blue, 3)

			size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 2, 2)
			pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
			gocv.PutText(&img, "Human", pt, gocv.FontHersheyPlain, 2, blue, 2)
		}

		buf, _ := gocv.IMEncode(".jpg", img)
		stream.UpdateJPEG(buf)
	}
}
