package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gocv.io/x/cvscope/scope"
	"gocv.io/x/gocv"
)

func init() {
	rootCmd.AddCommand(laplacianCmd)
}

var currentLaplacianBorder int
var laplacianSizeTracker, laplacianScaleTracker, laplacianDeltaTracker *gocv.Trackbar
var laplacianSize int
var laplacianScale, laplacianDelta float64

var laplacianCmd = &cobra.Command{
	Use:   "laplacian",
	Short: "Apply Laplacian to video images",
	Long: `Apply Laplacian to video images.

Key commands:
  Use 'z' and 'x' keys to page through border calculation types.
  Press 'esc' to exit.
  Press 'space' to pause/resume filtering.
  Press 'g' to generate Go code based on the current filter.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleLaplacianCmd()
	},
}

func handleLaplacianCmd() {
	video, err := scope.OpenVideoCapture(videoSource)
	if err != nil {
		fmt.Printf("Error opening video: %v\n", err)
		return
	}
	defer video.Close()

	window := gocv.NewWindow(laplacianWindowTitle())
	defer window.Close()

	laplacianSizeTracker = window.CreateTrackbar("size", 31)
	laplacianSizeTracker.SetPos(1)

	laplacianScaleTracker = window.CreateTrackbar("scale", 60)
	laplacianScaleTracker.SetPos(0)

	laplacianDeltaTracker = window.CreateTrackbar("delta", 60)
	laplacianDeltaTracker.SetPos(0)

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
		validateLaplacianTrackers()

		// Laplacian image proccessing filter
		gocv.Laplacian(img, &processed, gocv.MatTypeCV16S, laplacianSize, laplacianScale, laplacianDelta, getCurrentBorder(currentLaplacianBorder))

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
			currentLaplacianBorder = prevBorder(currentLaplacianBorder)
			window.SetWindowTitle(laplacianWindowTitle())
		case 120:
			// 'x'
			currentLaplacianBorder = nextBorder(currentLaplacianBorder)
			window.SetWindowTitle(laplacianWindowTitle())
		case 103:
			// 'g'
			laplacianGoCodeFragment(laplacianSize, laplacianScale, laplacianDelta, getCurrentBorderDescription(currentLaplacianBorder))
		case 112:
			// 'p'
			laplacianPythonCodeFragment(laplacianSize, laplacianScale, laplacianDelta, currentLaplacianBorder)
		case 32:
			// 'space'
			pause = !pause
			text := laplacianWindowTitle()
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
func validateLaplacianTrackers() {
	if laplacianSizeTracker.GetPos() == 0 {
		laplacianSizeTracker.SetPos(1)
	}

	laplacianSize = ensureOdd(laplacianSizeTracker)
	laplacianScale = float64(laplacianScaleTracker.GetPos())
	laplacianDelta = float64(laplacianDeltaTracker.GetPos())
}

func laplacianWindowTitle() string {
	return "Laplacian - " + getCurrentBorderDescription(currentLaplacianBorder) + " - CVscope"
}

func laplacianGoCodeFragment(x int, sx, sy float64, borderType string) {
	codeFragmentHeader("Go")
	fmt.Printf("\ngocv.Laplacian(src, &dest, gocv.MatTypeCV16S, %d, %1.f, %1.f, gocv.%s)\n\n",
		x, sx, sy, borderType)
}

func laplacianPythonCodeFragment(x int, sx, sy float64, borderType int) {
	codeFragmentHeader("Python")
	fmt.Println("Not implemented.")
}
