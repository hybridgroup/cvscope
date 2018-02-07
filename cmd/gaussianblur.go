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
var ksize, sigmax, sigmay *gocv.Trackbar

var gaussianBlurCmd = &cobra.Command{
	Use:   "gaussianblur",
	Short: "Apply Gaussian blur to video images",
	Long:  `Apply Gaussian blur to video images`,
	Run: func(cmd *cobra.Command, args []string) {
		handleGaussianBlurCmd()
	},
}

func handleGaussianBlurCmd() {
	webcam, err := gocv.VideoCaptureDevice(int(deviceID))
	if err != nil {
		fmt.Printf("Error opening video capture device: %v\n", deviceID)
		return
	}
	defer webcam.Close()

	window := gocv.NewWindow(gaussianBlurWindowTitle())
	defer window.Close()

	ksize = window.CreateTrackbar("ksize", 25)
	ksize.SetPos(0)

	sigmax = window.CreateTrackbar("sigma X", 60)
	sigmax.SetPos(30)

	sigmay = window.CreateTrackbar("sigma Y", 60)
	sigmay.SetPos(0)

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
		gocv.GaussianBlur(img, processed, image.Pt(ksize.GetPos(), ksize.GetPos()), float64(sigmax.GetPos()),
			float64(sigmay.GetPos()), getCurrentBorder(currentGaussianBlurBorder))

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
			gaussianBlurGoCodeFragment(ksize.GetPos(), ksize.GetPos(), float64(sigmax.GetPos()), float64(sigmay.GetPos()), getCurrentBorderDescription(currentGaussianBlurBorder))
		case 112:
			// 'p'

		case 27:
			// 'ESC'
			break MainLoop
		}
	}
}

// either ksize or sigmax have to be non-zero
func validateTrackers() {
	if ensureOdd(ksize) == 0 && sigmax.GetPos() == 0 {
		ksize.SetPos(1)
	}
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

func gaussianBlurPythonCodeFragment(pos int, threshType int) {
	codeFragmentHeader("Python")
	fmt.Println("Upgrade to pro")
}
