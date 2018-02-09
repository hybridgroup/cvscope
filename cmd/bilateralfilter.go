package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

func init() {
	rootCmd.AddCommand(bilateralFilterCmd)
}

var bilateralFilterCmd = &cobra.Command{
	Use:   "bilateral",
	Short: "Apply bilateral filter to video images",
	Long: `Apply bilateral filter to video images.

Commands:
  Press 'esc' to exit.
  Press 'space' to pause/resume filtering.
  Press 'g' to generate Go code based on the current filter.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleBilateralFilterCmd()
	},
}

func handleBilateralFilterCmd() {
	video, err := openVideoSource()
	if err != nil {
		fmt.Printf("Error opening video: %v\n", err)
		return
	}
	defer video.Close()

	window := gocv.NewWindow(bilateralFilterWindowTitle())
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	processed := gocv.NewMat()
	defer processed.Close()

	diameter := window.CreateTrackbar("diameter", 5)
	diameter.SetMin(1)
	diameter.SetPos(2)

	sigmaColor := window.CreateTrackbar("sigma color", 255)
	sigmaColor.SetPos(60)

	sigmaSpace := window.CreateTrackbar("sigma space", 255)
	sigmaSpace.SetPos(0)

	pause := false
	fmt.Printf("Start reading camera device: %v\n", deviceID)
MainLoop:
	for {
		if ok := video.Read(img); !ok {
			fmt.Printf("Error cannot read device %d\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// BilateralFilter image processing filter
		gocv.BilateralFilter(img, processed, diameter.GetPos(), float64(sigmaColor.GetPos()), float64(sigmaSpace.GetPos()))

		// Display the processed image?
		if pause {
			window.IMShow(img)
		} else {
			window.IMShow(processed)
		}

		// Check to see if the user has pressed any keys on the keyboard
		key := window.WaitKey(1)
		switch key {
		case 103:
			// 'g'
			bilateralFilterGoCodeFragment(diameter.GetPos(), float64(sigmaColor.GetPos()), float64(sigmaSpace.GetPos()))
		case 112:
			// 'p'
			bilateralFilterPythonCodeFragment(diameter.GetPos(), float64(sigmaColor.GetPos()), float64(sigmaSpace.GetPos()))
		case 32:
			// 'space'
			pause = !pause
			text := bilateralFilterWindowTitle()
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

func bilateralFilterWindowTitle() string {
	return "BilateralFilter - CVscope"
}

func bilateralFilterGoCodeFragment(d int, sC, sS float64) {
	codeFragmentHeader("Go")
	fmt.Printf("gocv.BilateralFilter(src, dest, %d, %1.f, %1.f)\n\n", d, sC, sS)
}

func bilateralFilterPythonCodeFragment(d int, sC, sS float64) {
	codeFragmentHeader("Python")
	fmt.Println("Not implemented.")
}
