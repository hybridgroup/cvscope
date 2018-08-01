package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

func init() {
	rootCmd.AddCommand(scharrCmd)
}

var currentScharrBorder int
var scharrDXTracker, scharrDYTracker, scharrScaleTracker, scharrDeltaTracker *gocv.Trackbar
var scharrDX, scharrDY int
var scharrScale, scharrDelta float64

var scharrCmd = &cobra.Command{
	Use:   "scharr",
	Short: "Apply Scharr to video images",
	Long: `Apply Scharr to video images.

Key commands:
  Use 'z' and 'x' keys to page through border calculation types.
  Press 'esc' to exit.
  Press 'space' to pause/resume filtering.
  Press 'g' to generate Go code based on the current filter.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleScharrCmd()
	},
}

func handleScharrCmd() {
	video, err := gocv.OpenVideoCapture(videoSource)
	if err != nil {
		fmt.Printf("Error opening video: %v\n", err)
		return
	}
	defer video.Close()

	window := gocv.NewWindow(scharrWindowTitle())
	defer window.Close()

	scharrDXTracker = window.CreateTrackbar("dx", 1)
	scharrDXTracker.SetPos(1)

	scharrDYTracker = window.CreateTrackbar("dy", 1)
	scharrDYTracker.SetPos(0)

	scharrScaleTracker = window.CreateTrackbar("scale", 60)
	scharrScaleTracker.SetPos(30)

	scharrDeltaTracker = window.CreateTrackbar("delta", 60)
	scharrDeltaTracker.SetPos(30)

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

		// make sure we do not have any invalid values
		validateScharrTrackers()

		// scharr image processing filter
		gocv.Scharr(img, &processed, gocv.MatTypeCV16S, scharrDX, scharrDY, scharrScale, scharrDelta, getCurrentBorder(currentScharrBorder))

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
			currentScharrBorder = prevBorder(currentScharrBorder)
			window.SetWindowTitle(scharrWindowTitle())
		case 120:
			// 'x'
			currentScharrBorder = nextBorder(currentScharrBorder)
			window.SetWindowTitle(scharrWindowTitle())
		case 103:
			// 'g'
			scharrGoCodeFragment(scharrDX, scharrDY, scharrScale, scharrDelta, getCurrentBorderDescription(currentScharrBorder))
		case 112:
			// 'p'
			scharrPythonCodeFragment(scharrDX, scharrDY, scharrScale, scharrDelta, currentScharrBorder)
		case 32:
			// 'space'
			pause = !pause
			text := scharrWindowTitle()
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
func validateScharrTrackers() {
	switch {
	case scharrDXTracker.GetPos() == 1:
		scharrDYTracker.SetPos(0)
	case scharrDYTracker.GetPos() == 1:
		scharrDXTracker.SetPos(0)
	case scharrDXTracker.GetPos() == 0 && scharrDYTracker.GetPos() == 0:
		scharrDYTracker.SetPos(1)
	}

	scharrDX = scharrDXTracker.GetPos()
	scharrDY = scharrDYTracker.GetPos()
	scharrScale = float64(scharrScaleTracker.GetPos())
	scharrDelta = float64(scharrDeltaTracker.GetPos())
}

func scharrWindowTitle() string {
	return "Scharr - " + getCurrentBorderDescription(currentScharrBorder) + " - CVscope"
}

func scharrGoCodeFragment(dx, dy int, sx, sy float64, borderType string) {
	codeFragmentHeader("Go")
	fmt.Printf("\ngocv.Scharr(src, &dest, gocv.MatTypeCV16S, %d, %d, %1.f, %1.f, gocv.%s)\n\n",
		dx, dy, sx, sy, borderType)
}

func scharrPythonCodeFragment(dx, dy int, sx, sy float64, borderType int) {
	codeFragmentHeader("Python")
	fmt.Println("Not implemented.")
}
