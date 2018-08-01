package cmd

import (
	"fmt"
	"image"

	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

func init() {
	rootCmd.AddCommand(blurCmd)
}

var blurCmd = &cobra.Command{
	Use:   "blur",
	Short: "Blur video images",
	Long: `Blur video images using a normalized box filter.

Key commands:
  Press 'esc' to exit.
  Press 'space' to pause/resume filtering.
  Press 'g' to generate Go code based on the current filter.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleBlurCmd()
	},
}

func handleBlurCmd() {
	video, err := gocv.OpenVideoCapture(videoSource)
	if err != nil {
		fmt.Printf("Error opening video: %v\n", err)
		return
	}
	defer video.Close()

	window := gocv.NewWindow(blurWindowTitle())
	defer window.Close()

	trackerX := window.CreateTrackbar("ksize X", 25)
	trackerX.SetMin(1)
	trackerX.SetPos(12)

	trackerY := window.CreateTrackbar("ksize Y", 25)
	trackerY.SetMin(1)
	trackerY.SetPos(12)

	img := gocv.NewMat()
	defer img.Close()

	processed := gocv.NewMat()
	defer processed.Close()

	pause := false
	fmt.Printf("Start reading video: %v\n", videoSource)
MainLoop:
	for {
		if ok := video.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", videoSource)
			return
		}
		if img.Empty() {
			continue
		}

		// Blur image proccessing filter
		gocv.Blur(img, &processed, image.Pt(trackerX.GetPos(), trackerY.GetPos()))

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
			blurGoCodeFragment(trackerX.GetPos(), trackerY.GetPos())
		case 112:
			// 'p'
			blurPythonCodeFragment(trackerX.GetPos(), trackerY.GetPos())
		case 32:
			// 'space'
			pause = !pause
			text := blurWindowTitle()
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

func blurWindowTitle() string {
	return "Blur - CVscope"
}

func blurGoCodeFragment(x, y int) {
	codeFragmentHeader("Go")
	fmt.Printf("gocv.Blur(src, &dest, image.Pt(%d, %d))\n\n", x, y)
}

func blurPythonCodeFragment(x, y int) {
	codeFragmentHeader("Python")
	fmt.Println("Not implemented.")
}
