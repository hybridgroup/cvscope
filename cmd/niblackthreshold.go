package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gocv.io/x/cvscope/scope"
	"gocv.io/x/gocv"
	"gocv.io/x/gocv/contrib"
)

func init() {
	rootCmd.AddCommand(niblackThresholdCmd)
}

var (
	currentBinarizationMethod int
	//currentAThreshold          int
	//blockSize                  int
	kValue, rValue     float32
	kTracker, rTracker *gocv.Trackbar
)

var niblackThresholdCmd = &cobra.Command{
	Use:   "niblack",
	Short: "Apply Niblack threshold to video images",
	Long: `Apply Niblack threshold to video images.

Key commands:
  Use 'z' and 'x' keys to page through Niblack threshold calculation types.
  Use 'a' and 's' keys to page through threshold calculation types.
  Press 'esc' to exit.
  Press 'w' to write JPG file.
  Press 'space' to pause/resume filtering.
  Press 'g' to generate Go code based on the current filter.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleNiblackThresholdCmd()
	},
}

func handleNiblackThresholdCmd() {
	video, err := scope.OpenVideoCapture(videoSource)
	if err != nil {
		fmt.Printf("Error opening video: %v\n", err)
		return
	}
	defer video.Close()

	window = gocv.NewWindow(niblackThresholdWindowTitle())
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	processed := gocv.NewMat()
	defer processed.Close()

	gray := gocv.NewMat()
	defer gray.Close()

	blockSizeTracker = window.CreateTrackbar("block size", 255)
	blockSizeTracker.SetPos(3)

	kTracker = window.CreateTrackbar("k", 10)
	kTracker.SetPos(5)

	rTracker = window.CreateTrackbar("r", 512)
	rTracker.SetPos(128)

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
		validateNiblackThresholdTrackers()

		// only works on grayscale images
		gocv.CvtColor(img, &gray, gocv.ColorBGRAToGray)

		// NiblackThreshold image processing filter
		contrib.NiblackThreshold(gray, &processed, 255.0, getCurrentAThreshold(), blockSize, kValue, getCurrentBinarizationMethod(), rValue)

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
			prevBinarizationMethod()
			window.SetWindowTitle(niblackThresholdWindowTitle())
		case xKey:
			nextBinarizationMethod()
			window.SetWindowTitle(niblackThresholdWindowTitle())
		case aKey:
			prevAThreshold()
			window.SetWindowTitle(niblackThresholdWindowTitle())
		case sKey:
			nextAThreshold()
			window.SetWindowTitle(niblackThresholdWindowTitle())
		case gKey:
			niblackThresholdGoCodeFragment(255.0, getCurrentAThresholdDescription(), blockSize, kValue, getCurrentBinarizationMethodDescription(), rValue)
		case pKey:
			niblackThresholdPythonCodeFragment(255.0, getCurrentBinarizationMethodDescription(), getCurrentAThresholdDescription(), blockSize, kValue)
		case space:
			handlePause(niblackThresholdWindowTitle())
		case wKey:
			writeFile("niblack", processed)
		case esc:
			return
		}
	}
}

// blocksize has to be odd. k value ranges from 0.0 to 1.0.
func validateNiblackThresholdTrackers() {
	blockSize = ensureOdd(blockSizeTracker)
	kValue = float32(kTracker.GetPos()) / 10.0
	rValue = float32(rTracker.GetPos())
}

func getCurrentBinarizationMethod() contrib.BinarizationMethod {
	return contrib.BinarizationMethod(currentBinarizationMethod)
}

func getCurrentBinarizationMethodDescription() string {
	switch currentBinarizationMethod {
	case 0:
		return "BinarizationNiblack"
	case 1:
		return "BinarizationSauvola"
	case 2:
		return "BinarizationWolf"
	case 3:
		return "BinarizationNICK"
	}

	return "Unknown"
}

func prevBinarizationMethod() {
	currentBinarizationMethod--
	if currentBinarizationMethod < 0 {
		currentBinarizationMethod = 3
	}
}

func nextBinarizationMethod() {
	currentBinarizationMethod = (currentBinarizationMethod + 1) % 4
}

// func getCurrentAThreshold() gocv.ThresholdType {
// 	return gocv.ThresholdType(currentAThreshold)
// }

// func getCurrentAThresholdDescription() string {
// 	switch currentAThreshold {
// 	case 0:
// 		return "ThresholdBinary"
// 	case 1:
// 		return "ThresholdBinaryInv"
// 	}

// 	return "Unknown"
// }

// func prevAThreshold() {
// 	currentAThreshold--
// 	if currentAThreshold < 0 {
// 		currentAThreshold = 1
// 	}
// }

// func nextAThreshold() {
// 	currentAThreshold = (currentAThreshold + 1) % 2
// }

func niblackThresholdWindowTitle() string {
	return "niBlackThreshold - " + getCurrentBinarizationMethodDescription() + " - " + getCurrentAThresholdDescription() + " - CVscope"
}

func niblackThresholdGoCodeFragment(mv float32, at string, bs int, k float32, t string, r float32) {
	codeFragmentHeader("Go")
	fmt.Printf("contrib.NiblackThreshold(src, &dest, %.1f, gocv.%s, %d, %.1f, contrib.%s, %.1f)\n\n", mv, at, bs, k, t, r)
}

func niblackThresholdPythonCodeFragment(mv float32, at string, t string, bs int, c float32) {
	codeFragmentHeader("Python")
	fmt.Println("Not implemented.")
}
