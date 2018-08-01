# CVscope

![CVscope](https://raw.githubusercontent.com/hybridgroup/cvscope/master/images/cvscope-gaussian.png)

CVscope is a CLI tool for visually exploring video image filters that are part of OpenCV ([http://opencv.org](http://opencv.org)).

CVscope itself is written in the Go programming language using GoCV ([https://gocv.io](https://gocv.io)).

## How to install

In order to use CVscope you must first have installed OpenCV and GoCV. Please follow the instructions located at [https://gocv.io/getting-started/](https://gocv.io/getting-started/) for your operating system.

You can then install CVscope by running:

    go get -u gocv.io/x/cvscope

## How to run

Running the main program without any parameters:

    cvscope

will display a list of the various valid CVscope commands:

    CVscope is a command line tool for OpenCV that lets you visually explore various 
    filters and algorithms for computer vision.

    It can also generate Go code that them implements the needed commands using the GoCV
    programming library. CVscope is itself written using GoCV.

    Usage:
    cvscope [command]

    Available Commands:
    bilateral   Apply bilateral filter to video images
    blur        Blur video images
    dilate      Dilate video images
    ...

Here is a CVscope command that opens a window using the default camera, and then displays the video using the `Blur` filter:

    cvscope blur

![CVscope](https://raw.githubusercontent.com/hybridgroup/cvscope/master/images/cvscope.png)

While the CVscope program is running you can change the values for `ksize X` and `ksize Y` by adjusting the sliders, and the video will display the current image filter settings in real-time.

You can also generate the Go code that matches the current image filter settings. By pressing the `g` key, the code is output to the command line window where you started CVscope running. For example, when running the `blur` command, pressing `g` outputs the following:

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

### License

CVscope is licensed under the Apache 2.0 license. Copyright 2018 The Hybrid Group.
