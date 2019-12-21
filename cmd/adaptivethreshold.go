package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gocv.io/x/cvscope/scope"
	"gocv.io/x/gocv"
)

func init() {
	rootCmd.AddCommand(adaptiveThresholdCmd)
}

var (
	currentAdaptiveThreshold   int
	currentAThreshold          int
	blockSize                  int
	cValue                     float32
	blockSizeTracker, cTracker *gocv.Trackbar
)

var adaptiveThresholdCmd = &cobra.Command{
	Use:   "adaptive",
	Short: "Apply adaptive threshold to video images",
	Long: `Apply adaptive threshold to video images.

Key commands:
  Use 'z' and 'x' keys to page through adaptive threshold calculation types.
  Use 'a' and 's' keys to page through threshold calculation types.
  Press 'esc' to exit.
  Press 'space' to pause/resume filtering.
  Press 'g' to generate Go code based on the current filter.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleAdaptiveThresholdCmd()
	},
}

func handleAdaptiveThresholdCmd() {
	video, err := scope.OpenVideoCapture(videoSource)
	if err != nil {
		fmt.Printf("Error opening video: %v\n", err)
		return
	}
	defer video.Close()

	window = gocv.NewWindow(adaptiveThresholdWindowTitle())
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	processed := gocv.NewMat()
	defer processed.Close()

	gray := gocv.NewMat()
	defer gray.Close()

	blockSizeTracker = window.CreateTrackbar("block size", 255)
	blockSizeTracker.SetPos(3)

	cTracker = window.CreateTrackbar("C", 512)
	cTracker.SetPos(256)

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
		validateAdaptiveThresholdTrackers()

		// only works on grayscale images
		gocv.CvtColor(img, &gray, gocv.ColorBGRAToGray)

		// AdaptiveThreshold image processing filter
		gocv.AdaptiveThreshold(gray, &processed, 255, getCurrentAdaptiveThreshold(), getCurrentAThreshold(), blockSize, cValue)

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
			prevAdaptiveThreshold()
			window.SetWindowTitle(adaptiveThresholdWindowTitle())
		case xKey:
			nextAdaptiveThreshold()
			window.SetWindowTitle(adaptiveThresholdWindowTitle())
		case aKey:
			prevAThreshold()
			window.SetWindowTitle(adaptiveThresholdWindowTitle())
		case sKey:
			nextAThreshold()
			window.SetWindowTitle(adaptiveThresholdWindowTitle())
		case gKey:
			adaptiveThresholdGoCodeFragment(255.0, getCurrentAdaptiveThresholdDescription(), getCurrentAThresholdDescription(), blockSize, cValue)
		case pKey:
			adaptiveThresholdPythonCodeFragment(255.0, getCurrentAdaptiveThresholdDescription(), getCurrentAThresholdDescription(), blockSize, cValue)
		case space:
			handlePause(adaptiveThresholdWindowTitle())
		case esc:
			return
		}
	}
}

// blocksize has to be odd. C value ranges from -256.0 to 256.0.
func validateAdaptiveThresholdTrackers() {
	blockSize = ensureOdd(blockSizeTracker)
	cValue = float32(cTracker.GetPos() - 256)
}

func getCurrentAdaptiveThreshold() gocv.AdaptiveThresholdType {
	return gocv.AdaptiveThresholdType(currentAdaptiveThreshold)
}

func getCurrentAdaptiveThresholdDescription() string {
	switch currentAdaptiveThreshold {
	case 0:
		return "AdaptiveThresholdMean"
	case 1:
		return "AdaptiveThresholdGaussian"
	}

	return "Unknown"
}

func prevAdaptiveThreshold() {
	currentAdaptiveThreshold--
	if currentAdaptiveThreshold < 0 {
		currentAdaptiveThreshold = 1
	}
}

func nextAdaptiveThreshold() {
	currentAdaptiveThreshold = (currentAdaptiveThreshold + 1) % 2
}

func getCurrentAThreshold() gocv.ThresholdType {
	return gocv.ThresholdType(currentAThreshold)
}

func getCurrentAThresholdDescription() string {
	switch currentAThreshold {
	case 0:
		return "ThresholdBinary"
	case 1:
		return "ThresholdBinaryInv"
	}

	return "Unknown"
}

func prevAThreshold() {
	currentAThreshold--
	if currentAThreshold < 0 {
		currentAThreshold = 1
	}
}

func nextAThreshold() {
	currentAThreshold = (currentAThreshold + 1) % 2
}

func adaptiveThresholdWindowTitle() string {
	return "AdaptiveThreshold - " + getCurrentAdaptiveThresholdDescription() + " - " + getCurrentAThresholdDescription() + " - CVscope"
}

func adaptiveThresholdGoCodeFragment(mv float32, at string, t string, bs int, c float32) {
	codeFragmentHeader("Go")
	fmt.Printf("gocv.AdaptiveThreshold(src, &dest, %1.f, gocv.%s, gocv.%s, %d, %1.f)\n\n", mv, at, t, bs, c)
}

func adaptiveThresholdPythonCodeFragment(mv float32, at string, t string, bs int, c float32) {
	codeFragmentHeader("Python")
	fmt.Println("Not implemented.")
}
