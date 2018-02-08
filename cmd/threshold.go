package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

func init() {
	rootCmd.AddCommand(thresholdCmd)
}

var currentThreshold int

var thresholdCmd = &cobra.Command{
	Use:   "threshold",
	Short: "Apply threshold filter to video images",
	Long:  `Apply threshold filter to video images`,
	Run: func(cmd *cobra.Command, args []string) {
		handleThresholdCmd()
	},
}

func handleThresholdCmd() {
	webcam, err := gocv.VideoCaptureDevice(int(deviceID))
	if err != nil {
		fmt.Printf("Error opening video capture device: %v\n", deviceID)
		return
	}
	defer webcam.Close()

	window := gocv.NewWindow(thresholdWindowTitle())
	defer window.Close()

	tracker := window.CreateTrackbar("threshold", 255)
	tracker.SetPos(128)

	img := gocv.NewMat()
	defer img.Close()

	processed := gocv.NewMat()
	defer processed.Close()

	pause := false
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

		// Threshold image proccessing filter
		gocv.Threshold(img, processed, float32(tracker.GetPos()), 255.0, getCurrentThreshold())

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
			prevThreshold()
			window.SetWindowTitle(thresholdWindowTitle())
		case 120:
			// 'x'
			nextThreshold()
			window.SetWindowTitle(thresholdWindowTitle())
		case 103:
			// 'g'
			thresholdGoCodeFragment(tracker.GetPos(), getCurrentThresholdDescription())
		case 112:
			// 'p'
			thresholdPythonCodeFragment(tracker.GetPos(), currentThreshold)
		case 32:
			// 'space'
			pause = !pause
			text := thresholdWindowTitle()
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

func getCurrentThreshold() gocv.ThresholdType {
	return gocv.ThresholdType(currentThreshold)
}

func getCurrentThresholdDescription() string {
	switch currentThreshold {
	case 0:
		return "ThresholdBinary"
	case 1:
		return "ThresholdBinaryInv"
	case 2:
		return "ThresholdTrunc"
	case 3:
		return "ThresholdToZero"
	case 4:
		return "ThresholdToZeroInv"
	}

	return "Unknown"
}

func prevThreshold() {
	currentThreshold--
	if currentThreshold < 0 {
		currentThreshold = 4
	}
}

func nextThreshold() {
	currentThreshold = (currentThreshold + 1) % 5
}

func thresholdWindowTitle() string {
	return getCurrentThresholdDescription() + " - Threshold - CV Toolkit"
}

func thresholdGoCodeFragment(pos int, threshType string) {
	codeFragmentHeader("Go")
	fmt.Printf("\ngocv.Threshold(src, dest, %.1f, 255.0, gocv.%s)\n\n", float32(pos), threshType)
}

func thresholdPythonCodeFragment(pos int, threshType int) {
	codeFragmentHeader("Python")
	fmt.Printf("\nretval, dest = cv.threshold(src, %.1f, 255.0, %d)\n\n", float32(pos), threshType)
}
