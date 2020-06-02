package cmd

import (
	"fmt"
	"image"

	"github.com/spf13/cobra"
	"gocv.io/x/cvscope/scope"
	"gocv.io/x/gocv"
)

func init() {
	rootCmd.AddCommand(dilateCmd)
}

var currentDilateShape int

var dilateCmd = &cobra.Command{
	Use:   "dilate",
	Short: "Dilate video images",
	Long: `Dilate video images.

Key commands:
  Use 'z' and 'x' keys to page through structuring element shapes.
  Press 'esc' to exit.
  Press 'space' to pause/resume filtering.
  Press 'g' to generate Go code based on the current filter.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleDilateCmd()
	},
}

func handleDilateCmd() {
	video, err := scope.OpenVideoCapture(videoSource)
	if err != nil {
		fmt.Printf("Error opening video: %v\n", err)
		return
	}
	defer video.Close()

	window = gocv.NewWindow(dilateWindowTitle())
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

	fmt.Printf("Start reading video: %v\n", videoSource)

	for {
		if ok := video.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", videoSource)
			return
		}
		if img.Empty() {
			continue
		}

		// Dilate image proccessing filter
		kernel := gocv.GetStructuringElement(getCurrentMorphShape(currentDilateShape), image.Pt(trackerX.GetPos(), trackerY.GetPos()))
		gocv.Dilate(img, &processed, kernel)
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
		case zKey:
			currentDilateShape = prevShape(currentDilateShape)
			window.SetWindowTitle(dilateWindowTitle())
		case xKey:
			currentDilateShape = nextShape(currentDilateShape)
			window.SetWindowTitle(dilateWindowTitle())
		case gKey:
			dilateGoCodeFragment(getCurrentMorphShapeDescription(currentDilateShape), trackerX.GetPos(), trackerY.GetPos())
		case pKey:
			dilatePythonCodeFragment(currentDilateShape, trackerX.GetPos(), trackerY.GetPos())
		case space:
			handlePause(dilateWindowTitle())
		case wKey:
			writeFile("dilate", processed)
		case esc:
			return
		}
	}
}

func dilateWindowTitle() string {
	return "Dilate - " + getCurrentMorphShapeDescription(currentDilateShape) + " - CVscope"
}

func dilateGoCodeFragment(morphType string, x, y int) {
	codeFragmentHeader("Go")
	fmt.Printf("\nkernel := gocv.GetStructuringElement(gocv.%s, image.Pt(%d, %d))\n", morphType, x, y)
	fmt.Printf("gocv.Dilate(src, &dest, kernel)\n\n")
}

func dilatePythonCodeFragment(morphType, x, y int) {
	codeFragmentHeader("Python")
	fmt.Println("Not implemented.")
}
