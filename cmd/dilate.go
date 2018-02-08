package cmd

import (
	"fmt"
	"image"

	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

func init() {
	rootCmd.AddCommand(dilateCmd)
}

var currentDilateShape int

var dilateCmd = &cobra.Command{
	Use:   "dilate",
	Short: "Dilate video images",
	Long:  `Dilate video images`,
	Run: func(cmd *cobra.Command, args []string) {
		handleDilateCmd()
	},
}

func handleDilateCmd() {
	webcam, err := gocv.VideoCaptureDevice(deviceID)
	if err != nil {
		fmt.Printf("Error opening video capture device: %v\n", deviceID)
		return
	}
	defer webcam.Close()

	window := gocv.NewWindow(dilateWindowTitle())
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

		// Dilate image proccessing filter
		kernel := gocv.GetStructuringElement(getCurrentMorphShape(currentDilateShape), image.Pt(trackerX.GetPos(), trackerY.GetPos()))
		gocv.Dilate(img, processed, kernel)
		kernel.Close()

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
			currentDilateShape = prevShape(currentDilateShape)
			window.SetWindowTitle(dilateWindowTitle())
		case 120:
			// 'x'
			currentDilateShape = nextShape(currentDilateShape)
			window.SetWindowTitle(dilateWindowTitle())
		case 103:
			// 'g'
			dilateGoCodeFragment(getCurrentMorphShapeDescription(currentDilateShape), trackerX.GetPos(), trackerY.GetPos())
		case 112:
			// 'p'
			dilatePythonCodeFragment(currentDilateShape, trackerX.GetPos(), trackerY.GetPos())
		case 32:
			// 'space'
			pause = !pause
			text := dilateWindowTitle()
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

func dilateWindowTitle() string {
	return getCurrentMorphShapeDescription(currentDilateShape) + " - Dilate - CVscope"
}

func dilateGoCodeFragment(morphType string, x, y int) {
	codeFragmentHeader("Go")
	fmt.Printf("\nkernel := gocv.GetStructuringElement(gocv.%s, image.Pt(%d, %d))\n", morphType, x, y)
	fmt.Printf("gocv.Dilate(src, dest, kernel)\n\n")
}

func dilatePythonCodeFragment(morphType, x, y int) {
	codeFragmentHeader("Python")
	fmt.Println("Not implemented.")
}
