package cmd

import (
	"fmt"
	"image"

	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

func init() {
	rootCmd.AddCommand(gaussianBlurCmd)
}

var currentGaussianBlurBorder int
var ksizeX, ksizeY, sigmaX, sigmaY *gocv.Trackbar
var kX, kY int
var sX, sY float64

var gaussianBlurCmd = &cobra.Command{
	Use:   "gaussianblur",
	Short: "Apply Gaussian blur to video images",
	Long:  `Apply Gaussian blur to video images`,
	Run: func(cmd *cobra.Command, args []string) {
		handleGaussianBlurCmd()
	},
}

func handleGaussianBlurCmd() {
	webcam, err := gocv.VideoCaptureDevice(deviceID)
	if err != nil {
		fmt.Printf("Error opening video capture device: %v\n", deviceID)
		return
	}
	defer webcam.Close()

	window := gocv.NewWindow(gaussianBlurWindowTitle())
	defer window.Close()

	ksizeX = window.CreateTrackbar("ksize X", 25)
	ksizeX.SetPos(0)

	ksizeY = window.CreateTrackbar("ksize Y", 25)
	ksizeY.SetPos(0)

	sigmaX = window.CreateTrackbar("sigma X", 60)
	sigmaX.SetPos(30)

	sigmaY = window.CreateTrackbar("sigma Y", 60)
	sigmaY.SetPos(0)

	img := gocv.NewMat()
	defer img.Close()

	processed := gocv.NewMat()
	defer processed.Close()

	fmt.Printf("Start reading camera device: %v\n", deviceID)
MainLoop:
	for {
		if ok := webcam.Read(img); !ok {
			fmt.Printf("Error cannot read device %d\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// make sure we do not have any invalid values
		validateTrackers()

		// GaussianBlur image proccessing filter
		gocv.GaussianBlur(img, processed, image.Pt(kX, kY), sX, sY, getCurrentBorder(currentGaussianBlurBorder))

		// Display the processed image
		window.IMShow(processed)

		// Check to see if the user has pressed any keys on the keyboard
		key := window.WaitKey(1)
		switch key {
		case 122:
			// 'z'
			currentGaussianBlurBorder = prevBorder(currentGaussianBlurBorder)
			window.SetWindowTitle(gaussianBlurWindowTitle())
		case 120:
			// 'x'
			currentGaussianBlurBorder = nextBorder(currentGaussianBlurBorder)
			window.SetWindowTitle(gaussianBlurWindowTitle())
		case 103:
			// 'g'
			gaussianBlurGoCodeFragment(kX, kY, sX, sY, getCurrentBorderDescription(currentGaussianBlurBorder))
		case 112:
			// 'p'
			gaussianBlurPythonCodeFragment(kX, kY, sX, sY, currentGaussianBlurBorder)
		case 27:
			// 'ESC'
			break MainLoop
		}
	}
}

// either ksize or sigmax have to be non-zero
func validateTrackers() {
	if sigmaX.GetPos() == 0 {
		if ksizeX.GetPos() == 0 {
			ksizeX.SetPos(1)
		}
		if ksizeY.GetPos() == 0 {
			ksizeY.SetPos(1)
		}
	}

	kX = ensureOdd(ksizeX)
	kY = ensureOdd(ksizeY)
	sX = float64(sigmaX.GetPos())
	sY = float64(sigmaY.GetPos())
}

// ksize has to be either 0 or an odd number
func ensureOdd(tracker *gocv.Trackbar) int {
	size := tracker.GetPos()

	if size%2 == 1 || size == 0 {
		return size
	}

	tracker.SetPos(size - 1)
	return size - 1
}

func gaussianBlurWindowTitle() string {
	return getCurrentBorderDescription(currentGaussianBlurBorder) + " - Gaussian Blur - CV Toolkit"
}

func gaussianBlurGoCodeFragment(x, y int, sx, sy float64, borderType string) {
	codeFragmentHeader("Go")
	fmt.Printf("\ngocv.GaussianBlur(src, dest, image.Pt(%d, %d), %1.f, %1.f, gocv.%s)\n\n",
		x, y, sx, sy, borderType)
}

func gaussianBlurPythonCodeFragment(x, y int, sx, sy float64, borderType int) {
	codeFragmentHeader("Python")
	fmt.Println("Not implemented.")
}
