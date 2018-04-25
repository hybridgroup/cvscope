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
var gaussianKsizeX, gaussianKsizeY, gaussianSigmaX, gaussianSigmaY *gocv.Trackbar
var gaussianKX, gaussianKY int
var gaussianSX, gaussianSY float64

var gaussianBlurCmd = &cobra.Command{
	Use:   "gaussian",
	Short: "Apply Gaussian blur to video images",
	Long: `Apply Gaussian blur to video images.

Commands:
  Use 'z' and 'x' keys to page through border calculation types.
  Press 'esc' to exit.
  Press 'space' to pause/resume filtering.
  Press 'g' to generate Go code based on the current filter.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleGaussianBlurCmd()
	},
}

func handleGaussianBlurCmd() {
	video, err := openVideoSource()
	if err != nil {
		fmt.Printf("Error opening video: %v\n", err)
		return
	}
	defer video.Close()

	window := gocv.NewWindow(gaussianBlurWindowTitle())
	defer window.Close()

	gaussianKsizeX = window.CreateTrackbar("ksize X", 25)
	gaussianKsizeX.SetPos(0)

	gaussianKsizeY = window.CreateTrackbar("ksize Y", 25)
	gaussianKsizeY.SetPos(0)

	gaussianSigmaX = window.CreateTrackbar("sigma X", 60)
	gaussianSigmaX.SetPos(30)

	gaussianSigmaY = window.CreateTrackbar("sigma Y", 60)
	gaussianSigmaY.SetPos(0)

	img := gocv.NewMat()
	defer img.Close()

	processed := gocv.NewMat()
	defer processed.Close()

	pause := false
	fmt.Printf("Start reading camera device: %v\n", deviceID)
MainLoop:
	for {
		if ok := video.Read(&img); !ok {
			fmt.Printf("Error cannot read device %d\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// make sure we do not have any invalid values
		validateGaussianBlurTrackers()

		// GaussianBlur image proccessing filter
		gocv.GaussianBlur(img, &processed, image.Pt(gaussianKX, gaussianKY), gaussianSX, gaussianSY, getCurrentBorder(currentGaussianBlurBorder))

		// Display the processed image?
		if pause {
			window.IMShow(img)
		} else {
			window.IMShow(processed)
		}

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
			gaussianBlurGoCodeFragment(gaussianKX, gaussianKY, gaussianSX, gaussianSY, getCurrentBorderDescription(currentGaussianBlurBorder))
		case 112:
			// 'p'
			gaussianBlurPythonCodeFragment(gaussianKX, gaussianKY, gaussianSX, gaussianSY, currentGaussianBlurBorder)
		case 32:
			// 'space'
			pause = !pause
			text := gaussianBlurWindowTitle()
			if pause {
				text = "**PAUSED** " + text
			}
			window.SetWindowTitle(text)
		case 27:
			// 'ESC'
			break MainLoop
		}
	}
}

// either ksize or sigmax have to be non-zero
func validateGaussianBlurTrackers() {
	if gaussianSigmaX.GetPos() == 0 {
		if gaussianKsizeX.GetPos() == 0 {
			gaussianKsizeX.SetPos(1)
		}
		if gaussianKsizeY.GetPos() == 0 {
			gaussianKsizeY.SetPos(1)
		}
	}

	gaussianKX = ensureOdd(gaussianKsizeX)
	gaussianKY = ensureOdd(gaussianKsizeY)
	gaussianSX = float64(gaussianSigmaX.GetPos())
	gaussianSY = float64(gaussianSigmaY.GetPos())
}

func gaussianBlurWindowTitle() string {
	return getCurrentBorderDescription(currentGaussianBlurBorder) + " - Gaussian Blur - CV Toolkit"
}

func gaussianBlurGoCodeFragment(x, y int, sx, sy float64, borderType string) {
	codeFragmentHeader("Go")
	fmt.Printf("\ngocv.GaussianBlur(src, &dest, image.Pt(%d, %d), %1.f, %1.f, gocv.%s)\n\n",
		x, y, sx, sy, borderType)
}

func gaussianBlurPythonCodeFragment(x, y int, sx, sy float64, borderType int) {
	codeFragmentHeader("Python")
	fmt.Println("Not implemented.")
}
