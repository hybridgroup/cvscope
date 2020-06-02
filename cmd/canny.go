package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gocv.io/x/cvscope/scope"
	"gocv.io/x/gocv"
)

func init() {
	rootCmd.AddCommand(cannyCmd)
}

var currentCannyShape int

var cannyCmd = &cobra.Command{
	Use:   "canny",
	Short: "canny video images",
	Long: `canny video images.

Key commands:
  Use 'z' and 'x' keys to page through structuring element shapes.
  Press 'esc' to exit.
  Press 'space' to pause/resume filtering.
  Press 'g' to generate Go code based on the current filter.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleCannyCmd()
	},
}

func handleCannyCmd() {
	video, err := scope.OpenVideoCapture(videoSource)
	if err != nil {
		fmt.Printf("Error opening video: %v\n", err)
		return
	}
	defer video.Close()

	window = gocv.NewWindow(cannyWindowTitle())
	defer window.Close()

	tracker1 := window.CreateTrackbar("t1", 100)
	tracker1.SetPos(50)

	tracker2 := window.CreateTrackbar("t2", 100)
	tracker2.SetPos(50)

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

		gray := gocv.NewMat()
		defer gray.Close()

		// only works on grayscale images
		gocv.CvtColor(img, &gray, gocv.ColorBGRAToGray)

		// canny image proccessing filter
		gocv.Canny(gray, &processed, float32(tracker1.GetPos()), float32(tracker2.GetPos()))

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
			cannyGoCodeFragment(float32(tracker1.GetPos()), float32(tracker2.GetPos()))
		case pKey:
			cannyPythonCodeFragment(float32(tracker1.GetPos()), float32(tracker2.GetPos()))
		case space:
			handlePause(cannyWindowTitle())
		case wKey:
			writeFile("canny", processed)
		case esc:
			return
		}
	}
}

func cannyWindowTitle() string {
	return "Canny - CVscope"
}

func cannyGoCodeFragment(t1, t2 float32) {
	codeFragmentHeader("Go")
	fmt.Printf("\ngocv.Canny(src, &dest, %1.f, %1.f))\n", t1, t2)
}

func cannyPythonCodeFragment(t1, t2 float32) {
	codeFragmentHeader("Python")
	fmt.Println("Not implemented.")
}
