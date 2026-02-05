package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/blackjack/webcam"
)

const BufferSize = 12 

func main() { 
	cam, err := webcam.Open("/dev/video0")
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}
	defer cam.Close()

	var format webcam.PixelFormat
	var framesize webcam.FrameSize

	for f, s := range cam.GetSupportedFormats() {
		fmt.Printf("format: %s\n", s)
		format = f
		break
	}

	for _, s := range cam.GetSupportedFrameSizes(format) {
		framesize = s
		fmt.Printf("framesize: %s\n", framesize.GetString())
		break
	}

	_, _, _, err = cam.SetImageFormat(format, framesize.MaxWidth, framesize.MaxHeight)
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	err = cam.StartStreaming()
	if err != nil {
		log.Fatalf("error: %s\n", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go readWebCam(signalChan, cam)

	<-signalChan
	time.Sleep(time.Second)
}

func readWebCam(signalChan chan os.Signal, cam *webcam.Webcam) {

	tui := InitTui()
	defer tui.CloseTui() 

	frameChan := make(chan *[]byte, BufferSize)

	go writeToFile(frameChan, tui)
	go processImage(frameChan, tui)

	for {
		select {
		case <-signalChan:
			close(frameChan)
			log.Printf("Closing everything now...")
			os.Exit(0)
		default:
			err := cam.WaitForFrame(5)
			if err != nil { 
				log.Fatalf("error: %s\n", err)
			}

			frame, err := cam.ReadFrame()
			if err != nil { 
				log.Fatalf("error: %s\n", err)
			}

			frameChan <- &frame
		}
	}
}

func processImage(frameChn chan *[]byte, tui *Tui) {
	var i int 
	for frame := range frameChn { 
		reader := bytes.NewReader(*frame)
		_, err := jpeg.Decode(reader)
		if err != nil { 
			log.Fatalf("error: %s\n", err)
		}

		tui.UpdateProcessText(fmt.Sprintf("process frame number: %d\n", i)) 
		i += 1 
	}

	log.Printf("closed image processing worker\n") 
}

func writeToFile(frameChn chan *[]byte, tui *Tui) { 
	var i int 
	for frame := range frameChn { 
		err := os.WriteFile(fmt.Sprintf("frames/%d.jpg", i), *frame, 0644)
		if err != nil { 
			log.Fatalf("error: %s\n", err)
		}
		
		tui.UpdateWriterText(fmt.Sprintf("Wrote: frames/%d.jpg", i))

		i += 1 
	}

	log.Printf("closed file writing worker\n")
}
