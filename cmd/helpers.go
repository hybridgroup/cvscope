package cmd

import (
	"fmt"

	"gocv.io/x/gocv"
)

// type CVCommander interface {
// 	Init()
// 	Validate()
// 	WindowTitle() string
// 	Process(input gocv.Mat, output gocv.Mat)
// }

const (
	zKey  = 122
	xKey  = 120
	aKey  = 97
	sKey  = 115
	gKey  = 103
	pKey  = 112
	wKey  = 119
	space = 32
	esc   = 27
)

func codeFragmentHeader(lang string) {
	fmt.Println("===============================")
	fmt.Printf("%s code for current filter:\n", lang)
	fmt.Println("-------------------------------")
}

func prevShape(currentShape int) int {
	currentShape--
	if currentShape < 0 {
		currentShape = 2
	}

	return currentShape
}

func nextShape(currentShape int) int {
	return (currentShape + 1) % 3
}

func getCurrentMorphShape(currentShape int) gocv.MorphShape {
	return gocv.MorphShape(currentShape)
}

func getCurrentMorphShapeDescription(currentShape int) string {
	switch currentShape {
	case 0:
		return "MorphRect"
	case 1:
		return "MorphCross"
	case 2:
		return "MorphEllipse"
	}

	return "Unknown"
}

func getCurrentBorder(currentBorder int) gocv.BorderType {
	return gocv.BorderType(currentBorder)
}

func getCurrentBorderDescription(currentBorder int) string {
	switch currentBorder {
	case 0:
		return "BorderConstant"
	case 1:
		return "BorderReplicate"
	case 2:
		return "BorderReflect"
	case 4:
		return "BorderReflect101"
	}

	return "Unknown"
}

func prevBorder(currentBorder int) int {
	currentBorder--
	if currentBorder < 0 {
		currentBorder = 4
	}
	if currentBorder == 3 {
		currentBorder = 2
	}

	return currentBorder
}

func nextBorder(currentBorder int) int {
	currentBorder = (currentBorder + 1) % 5
	if currentBorder == 3 {
		currentBorder = 4
	}

	return currentBorder
}

// ksize has to be either 0 or an odd number
func ensureOdd(tracker *gocv.Trackbar) int {
	size := tracker.GetPos()

	if size%2 == 1 || size == 0 {
		return size
	}

	tracker.SetPos(size - 1)
	return size - 1
}

func handlePause(text string) {
	pause = !pause
	if pause {
		text = "**PAUSED** " + text
	}
	window.SetWindowTitle(text)
}

func writeFile(cmdName string, img gocv.Mat) {
	gocv.IMWrite(cmdName+".jpg", img)
}
