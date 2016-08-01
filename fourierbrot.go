package main

import (
	"fmt"
	"github.com/cryptix/wav"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/mjibson/go-dsp/fft"
	"image"
	"image/color"
	"io"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: fourierbrot <file.wav>\n")
		os.Exit(1)
	}

	wavReader := openWav(os.Args[1])

	fmt.Println(wavReader)

	frames := int(wavReader.GetSampleCount() / 8192)
	for i := 0; i < frames; i++ {
		frame := readFrame(wavReader)
		fftData := fourierTransform(frame)
		display(i, fftData)
	}
}

func readFrame(reader *wav.Reader) []float64 {
	var frame = make([]float64, 4096)
	for i := range frame {
		_, err := reader.ReadSample()
		if err == io.EOF {
			break
		}
		checkErr(err)
		sample, err := reader.ReadSample()
		if err == io.EOF {
			break
		}
		checkErr(err)

		var signedSample int16
		if uint16(sample)&0x8000 == 0 {
			signedSample = int16(sample)
		} else {
			signedSample = int16(sample&0x7FFF - 0x8000)
		}
		frame[i] = float64(signedSample) / 0x4000
	}
	return frame
}

func fourierTransform(frame []float64) []complex128 {
	return fft.FFTReal(frame)
}

func display(j int, fftData []complex128) {
	dest := image.NewRGBA(image.Rect(0, 0, 400, 400))
	gc := draw2dimg.NewGraphicContext(dest)

	gc.SetStrokeColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
	gc.SetLineWidth(2)
	gc.MoveTo(200, 200)

	for i := range fftData {
		x := 200 + (2 * real(fftData[i]))
		y := 200 + (2 * imag(fftData[i]))
		gc.LineTo(x, y)
	}
	gc.Close()
	gc.Stroke()

	draw2dimg.SaveToPngFile(fmt.Sprintf("%04d.png", j), dest)
}

func openWav(filename string) *wav.Reader {
	testInfo, err := os.Stat(filename)
	checkErr(err)

	testWav, err := os.Open(filename)
	checkErr(err)

	wavReader, err := wav.NewReader(testWav, testInfo.Size())
	checkErr(err)

	return wavReader
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
