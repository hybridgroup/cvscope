package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gocv.io/x/cvscope/scope"
	"gocv.io/x/gocv"
)

func init() {
	rootCmd.AddCommand(medianBlurCmd)
}

var medianBlurCmd = &cobra.Command{
	Use:   "medianblur",
	Short: "MedianBlur video images",
	Long: `Blur video images using a median filter.

Key commands:
  Press 'esc' to exit.
  Press 'space' to pause/resume filtering.
  Press 'g' to generate Go code based on the current filter.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleMedianBlurCmd()
	},
}

func handleMedianBlurCmd() {
	video, err := scope.OpenVideoCapture(videoSource)
	if err != nil {
		fmt.Printf("Error opening video: %v\n", err)
		return
	}
	defer video.Close()

	window = gocv.NewWindow(medianBlurWindowTitle())
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	processed := gocv.NewMat()
	defer processed.Close()

	tracker := window.CreateTrackbar("ksize", 25)
	tracker.SetMin(0)
	tracker.SetPos(5)

	fmt.Printf("Start reading video: %v\n", videoSource)

	for {
		if ok := video.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", videoSource)
			return
		}
		if img.Empty() {
			continue
		}

		// MedianBlur image processing filter
		ksize := ensureOdd(tracker)
		gocv.MedianBlur(img, &processed, ksize)

		// Display the processed image?
		if pause {
			window.IMShow(img)
		} else {
			window.IMShow(processed)
		}

		// Check to see if the user has pressed any keys on the keyboard
		key := window.WaitKey(1)
		switch key {
		case gKey:
			medianBlurGoCodeFragment(tracker.GetPos())
		case pKey:
			medianBlurPythonCodeFragment(tracker.GetPos())
		case space:
			handlePause(medianBlurWindowTitle())
		case wKey:
			writeFile("medianblur", processed)
		case esc:
			return
		}
	}
}

func medianBlurWindowTitle() string {
	return "MedianBlur - CVscope"
}

func medianBlurGoCodeFragment(x int) {
	codeFragmentHeader("Go")
	fmt.Printf("gocv.MedianBlur(src, &dest, %d)\n\n", x)
}

func medianBlurPythonCodeFragment(x int) {
	codeFragmentHeader("Python")
	fmt.Println("Not implemented.")
}
