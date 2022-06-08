package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"gocv.io/x/cvscope/scope"
	"gocv.io/x/gocv"
)

func init() {
	rootCmd.AddCommand(convertCmd)
}

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "convert video images",
	Long: `Convert video images using brightness and contrast.
	
Key commands:
  Press 'esc' to exit.
  Press 'space' to pause/resume filtering.
  Press 'g' to generate Go code based on the current filter.`,
	Run: func(cmd *cobra.Command, args []string) {
		handleConvertCmd()
	},
}

func handleConvertCmd() {
	video, err := scope.OpenVideoCapture(videoSource)
	if err != nil {
		fmt.Printf("Error opening video: %v\n", err)
		return
	}
	defer video.Close()

	window = gocv.NewWindow(convertWindowTitle())
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	processed := gocv.NewMat()
	defer processed.Close()

	trackerA := window.CreateTrackbar("alpha", 30)
	trackerA.SetMin(10)
	trackerA.SetPos(12)

	trackerB := window.CreateTrackbar("beta", 100)
	trackerB.SetMin(0)
	trackerB.SetPos(20)

	fmt.Printf("Start reading video: %v\n", videoSource)

	for {
		if ok := video.Read(&img); !ok {
			fmt.Printf("Device closed: %v\n", videoSource)
			return
		}
		if img.Empty() {
			continue
		}

		// convert image proccessing filter
		img.ConvertToWithParams(&processed, -1, float32(trackerA.GetPos())/10.0, float32(trackerB.GetPos()))

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
			convertGoCodeFragment(float32(trackerA.GetPos())/10.0, float32(trackerB.GetPos()))
		case pKey:
			convertPythonCodeFragment(trackerA.GetPos(), trackerB.GetPos())
		case space:
			handlePause(convertWindowTitle())
		case wKey:
			writeFile("convert", processed)
		case esc:
			return
		}
	}
}

func convertWindowTitle() string {
	return "Convert - CVscope"
}

func convertGoCodeFragment(a, b float32) {
	codeFragmentHeader("Go")
	fmt.Printf("src.ConvertTo(&dest, %.1f, %.1f)\n\n", a, b)
}

func convertPythonCodeFragment(x, y int) {
	codeFragmentHeader("Python")
	fmt.Println("Not implemented.")
}
