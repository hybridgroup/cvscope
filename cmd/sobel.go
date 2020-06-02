package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gocv.io/x/cvscope/scope"
	"gocv.io/x/gocv"
)

func init() {
	rootCmd.AddCommand(sobelCmd)
}

var currentSobelBorder int
var sobelDXTracker, sobelDYTracker, sobelKSizeTracker, sobelScaleTracker, sobelDeltaTracker *gocv.Trackbar
var sobelDX, sobelDY, sobelKSize int
var sobelScale, sobelDelta float64

var sobelCmd = &cobra.Command{
	Use:   "sobel",
	Short: "Apply Sobel to video images",
	Long: `Apply Sobel to video images.

Key commands:
  Use 'z' and 'x' keys to page through border calculation types.
  Press 'esc' to exit.
  Press 'space' to pause/resume filtering.
  Press 'g' to generate Go code based on the current filter.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleSobelCmd()
	},
}

func handleSobelCmd() {
	video, err := scope.OpenVideoCapture(videoSource)
	if err != nil {
		fmt.Printf("Error opening video: %v\n", err)
		return
	}
	defer video.Close()

	window = gocv.NewWindow(sobelWindowTitle())
	defer window.Close()

	sobelDXTracker = window.CreateTrackbar("dx", 1)
	sobelDXTracker.SetPos(1)

	sobelDYTracker = window.CreateTrackbar("dy", 1)
	sobelDYTracker.SetPos(0)

	sobelKSizeTracker = window.CreateTrackbar("ksize", 7)
	sobelKSizeTracker.SetPos(3)

	sobelScaleTracker = window.CreateTrackbar("scale", 60)
	sobelScaleTracker.SetPos(30)

	sobelDeltaTracker = window.CreateTrackbar("delta", 60)
	sobelDeltaTracker.SetPos(30)

	img := gocv.NewMat()
	defer img.Close()

	processed := gocv.NewMat()
	defer processed.Close()

	fmt.Printf("Start reading video: %v\n", videoSource)

	for {
		if ok := video.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", videoSource)
			return
		}
		if img.Empty() {
			continue
		}

		// make sure we do not have any invalid values
		validateSobelTrackers()

		// sobel image processing filter
		gocv.Sobel(img, &processed, gocv.MatTypeCV16S, sobelDX, sobelDY, sobelKSize, sobelScale, sobelDelta, getCurrentBorder(currentSobelBorder))

		// Display the processed image?
		if pause {
			window.IMShow(img)
		} else {
			window.IMShow(processed)
		}

		// Check to see if the user has pressed any keys on the keyboard
		key := window.WaitKey(1)
		switch key {
		case zKey:
			currentSobelBorder = prevBorder(currentSobelBorder)
			window.SetWindowTitle(sobelWindowTitle())
		case xKey:
			currentSobelBorder = nextBorder(currentSobelBorder)
			window.SetWindowTitle(sobelWindowTitle())
		case gKey:
			sobelGoCodeFragment(sobelDX, sobelDY, sobelKSize, sobelScale, sobelDelta, getCurrentBorderDescription(currentSobelBorder))
		case pKey:
			sobelPythonCodeFragment(sobelDX, sobelDY, sobelKSize, sobelScale, sobelDelta, currentSobelBorder)
		case space:
			handlePause(sobelWindowTitle())
		case wKey:
			writeFile("sobel", processed)
		case esc:
			return
		}
	}
}

// either ksize or sigmax have to be non-zero
func validateSobelTrackers() {
	switch {
	case sobelDXTracker.GetPos() == 1:
		sobelDYTracker.SetPos(0)
	case sobelDYTracker.GetPos() == 1:
		sobelDXTracker.SetPos(0)
	case sobelDXTracker.GetPos() == 0 && sobelDYTracker.GetPos() == 0:
		sobelDYTracker.SetPos(1)
	}

	sobelDX = sobelDXTracker.GetPos()
	sobelDY = sobelDYTracker.GetPos()
	sobelKSize = ensureOdd(sobelKSizeTracker)
	sobelScale = float64(sobelScaleTracker.GetPos())
	sobelDelta = float64(sobelDeltaTracker.GetPos())
}

func sobelWindowTitle() string {
	return "Sobel - " + getCurrentBorderDescription(currentSobelBorder) + " - CVscope"
}

func sobelGoCodeFragment(dx, dy, k int, sx, sy float64, borderType string) {
	codeFragmentHeader("Go")
	fmt.Printf("\ngocv.Sobel(src, &dest, gocv.MatTypeCV16S, %d, %d, %d, %1.f, %1.f, gocv.%s)\n\n",
		dx, dy, k, sx, sy, borderType)
}

func sobelPythonCodeFragment(dx, dy, k int, sx, sy float64, borderType int) {
	codeFragmentHeader("Python")
	fmt.Println("Not implemented.")
}
