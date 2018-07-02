package cmd

import (
	"fmt"
	"image"

	"github.com/spf13/cobra"
	"gocv.io/x/gocv"
)

func init() {
	rootCmd.AddCommand(morphologyExCmd)
}

var currentMorphologyExShape int
var currentMorphOp int

var morphologyExCmd = &cobra.Command{
	Use:   "morph",
	Short: "Perform MorphologyEx operations on video images",
	Long: `Perform MorphologyEx operations on video images.

Commands:
  Use 'z' and 'x' keys to page through structuring element shapes.
  Use 'n' and 'm' keys to page through morphology operations.
  Press 'esc' to exit.
  Press 'space' to pause/resume filtering.
  Press 'g' to generate Go code based on the current filter.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleMorphologyExCmd()
	},
}

func handleMorphologyExCmd() {
	video, err := gocv.OpenVideoCapture(videoSource)
	if err != nil {
		fmt.Printf("Error opening video: %v\n", err)
		return
	}
	defer video.Close()

	window := gocv.NewWindow(morphologyExWindowTitle())
	defer window.Close()

	morphologyExTrackerX := window.CreateTrackbar("ksize X", 25)
	morphologyExTrackerX.SetMin(1)
	morphologyExTrackerX.SetPos(12)

	morphologyExTrackerY := window.CreateTrackbar("ksize Y", 25)
	morphologyExTrackerY.SetMin(1)
	morphologyExTrackerY.SetPos(12)

	img := gocv.NewMat()
	defer img.Close()

	processed := gocv.NewMat()
	defer processed.Close()

	pause := false
	fmt.Printf("Start reading video: %v\n", videoSource)
MainLoop:
	for {
		if ok := video.Read(&img); !ok {
			fmt.Printf("Start reading video: %v\n", videoSource)
			return
		}
		if img.Empty() {
			continue
		}

		// MorphologyEx image proccessing filter
		kernel := gocv.GetStructuringElement(getCurrentMorphShape(currentMorphologyExShape), image.Pt(morphologyExTrackerX.GetPos(), morphologyExTrackerY.GetPos()))
		gocv.MorphologyEx(img, &processed, gocv.MorphType(currentMorphOp), kernel)
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
			currentMorphologyExShape = prevShape(currentMorphologyExShape)
			window.SetWindowTitle(morphologyExWindowTitle())
		case 120:
			// 'x'
			currentMorphologyExShape = nextShape(currentMorphologyExShape)
			window.SetWindowTitle(morphologyExWindowTitle())
		case 110:
			// 'n'
			prevMorphOp()
			window.SetWindowTitle(morphologyExWindowTitle())
		case 109:
			// 'm'
			nextMorphOp()
			window.SetWindowTitle(morphologyExWindowTitle())
		case 103:
			// 'g'
			morphologyExGoCodeFragment(getCurrentMorphShapeDescription(currentMorphologyExShape), morphologyExTrackerX.GetPos(), morphologyExTrackerY.GetPos(), getCurrentMorphOpDescription(currentMorphOp))
		case 112:
			// 'p'
			morphologyExPythonCodeFragment(currentMorphologyExShape, morphologyExTrackerX.GetPos(), morphologyExTrackerY.GetPos(), currentMorphOp)
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

func getCurrentMorphOp() gocv.MorphType {
	return gocv.MorphType(currentMorphOp)
}

func getCurrentMorphOpDescription(op int) string {
	switch op {
	case 0:
		return "MorphErode"
	case 1:
		return "MorphDilate"
	case 2:
		return "MorphOpen"
	case 3:
		return "MorphClose"
	case 4:
		return "MorphGradient"
	case 5:
		return "MorphTophat"
	case 6:
		return "MorphBlackhat"
	}

	return "Unknown"
}

func prevMorphOp() {
	currentMorphOp--
	if currentMorphOp < 0 {
		currentMorphOp = 6
	}
}

func nextMorphOp() {
	currentMorphOp = (currentMorphOp + 1) % 7
}

func morphologyExWindowTitle() string {
	return getCurrentMorphOpDescription(currentMorphOp) + " - " +
		getCurrentMorphShapeDescription(currentMorphologyExShape) +
		" - MorphologyEx - CVscope"
}

func morphologyExGoCodeFragment(morphType string, x, y int, morphOp string) {
	codeFragmentHeader("Go")
	fmt.Printf("\nkernel := gocv.GetStructuringElement(gocv.%s, image.Pt(%d, %d))\n", morphType, x, y)
	fmt.Printf("gocv.MorphologyEx(src, &dest, gocv.%s, kernel)\n\n", morphOp)
}

func morphologyExPythonCodeFragment(morphType, x, y, morphOp int) {
	codeFragmentHeader("Python")
	fmt.Println("Not implemented.")
}
