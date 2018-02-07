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
	Short: "erode",
	Long:  `erode`,
	Run: func(cmd *cobra.Command, args []string) {
		handleErodeCmd()
	},
}

func handleErodeCmd() {
	webcam, err := gocv.VideoCaptureDevice(deviceID)
	if err != nil {
		fmt.Printf("Error opening video capture device: %v\n", deviceID)
		return
	}
	defer webcam.Close()

	window := gocv.NewWindow(erodeWindowTitle())
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	processed := gocv.NewMat()
	defer processed.Close()

	tracker := window.CreateTrackbar("ksize", 25)
	tracker.SetMin(1)
	tracker.SetPos(12)

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

		// Erode image proccessing filter
		kernel := gocv.GetStructuringElement(getCurrentMorphShape(currentErodeShape), image.Pt(tracker.GetPos(), tracker.GetPos()))
		defer kernel.Close()
		gocv.Erode(img, processed, kernel)

		// Display the processed image
		window.IMShow(processed)

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
			erodeGoCodeFragment(getCurrentMorphShapeDescription(currentErodeShape), tracker.GetPos(), tracker.GetPos())
		case 112:
			// 'p'
			erodePythonCodeFragment(tracker.GetPos(), currentErodeShape)
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

func erodePythonCodeFragment(pos int, morphType int) {
	codeFragmentHeader("Python")
	fmt.Println("Upgrade to pro")
}
