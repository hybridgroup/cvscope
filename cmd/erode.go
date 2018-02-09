package cmd

import (
	"fmt"
	"image"

	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

func init() {
	rootCmd.AddCommand(erodeCmd)
}

var currentErodeShape int

var erodeCmd = &cobra.Command{
	Use:   "erode",
	Short: "Erode video images",
	Long:  `Erode video images`,
	Run: func(cmd *cobra.Command, args []string) {
		handleErodeCmd()
	},
}

func handleErodeCmd() {
	video, err := openVideoSource()
	if err != nil {
		fmt.Printf("Error opening video: %v\n", err)
	}
	defer video.Close()

	window := gocv.NewWindow(erodeWindowTitle())
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	processed := gocv.NewMat()
	defer processed.Close()

	trackerX := window.CreateTrackbar("ksize X", 25)
	trackerX.SetMin(1)
	trackerX.SetPos(12)

	trackerY := window.CreateTrackbar("ksize Y", 25)
	trackerY.SetMin(1)
	trackerY.SetPos(12)

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

		// Erode image proccessing filter
		kernel := gocv.GetStructuringElement(getCurrentMorphShape(currentErodeShape), image.Pt(trackerX.GetPos(), trackerY.GetPos()))
		gocv.Erode(img, processed, kernel)
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
			currentErodeShape = prevShape(currentErodeShape)
			window.SetWindowTitle(erodeWindowTitle())
		case 120:
			// 'x'
			currentErodeShape = nextShape(currentErodeShape)
			window.SetWindowTitle(erodeWindowTitle())
		case 103:
			// 'g'
			erodeGoCodeFragment(getCurrentMorphShapeDescription(currentErodeShape), trackerX.GetPos(), trackerY.GetPos())
		case 112:
			// 'p'
			erodePythonCodeFragment(currentErodeShape, trackerX.GetPos(), trackerY.GetPos())
		case 32:
			// 'space'
			pause = !pause
			text := erodeWindowTitle()
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

func erodeWindowTitle() string {
	return getCurrentMorphShapeDescription(currentErodeShape) + " - Erode - CV Toolkit"
}

func erodeGoCodeFragment(morphType string, x, y int) {
	codeFragmentHeader("Go")
	fmt.Printf("\nkernel := gocv.GetStructuringElement(gocv.%s, image.Pt(%d, %d))\n", morphType, x, y)
	fmt.Printf("gocv.Erode(src, dest, kernel)\n\n")
}

func erodePythonCodeFragment(morphType, x, y int) {
	codeFragmentHeader("Python")
	fmt.Println("Not implemented.")
}
