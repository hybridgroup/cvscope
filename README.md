# CVscope

CVscope is a visual tool for exploring video image filters that are part of OpenCV.

CVscope itself is written in the Go programming language using GoCV.

In order to use CVscope you must first have installed GoCV.

## How to run

Just running the main program will display a list of the various valid CVscope commands:

    go run main.go

Here is a CVscope command that opens a window using the default camera, and then displays the video using the `Blur` filter:

    go run main.go blur

![CVscope](https://raw.githubusercontent.com/hybridgroup/cvscope/master/images/cvscope.png)

While the CVscope program is running, you can generate the line of GoCV code that matches the current image filter settings. By pressing the `g` key, the code is output to the command line window where you started CVscope running. For example, when running the above command, pressing `g` generates the following code:

    ===============================
    Go code for current filter:
    -------------------------------
    gocv.Blur(src, &dest, image.Pt(12, 12))

## How to build

CVscope does not yet support cross platform builds, so you must build the program on the desired target operating system.

### macOS/Linux

    go build -o cvscope main.go

### Windows

    go build -o cvscope.exe main.go

