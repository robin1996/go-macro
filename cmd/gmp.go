package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/robin1996/go-macro/macrofile"
	"github.com/go-vgo/robotgo"
)

func playStartSeq() {
	fmt.Println("Playing macro in...")
	time.Sleep(time.Second)
	fmt.Println("3...")
	time.Sleep(time.Second)
	fmt.Println("2...")
	time.Sleep(time.Second)
	fmt.Println("1...")
	time.Sleep(time.Second)
	fmt.Println("GO!")
}

func playAction(step macrofile.Step) {
	switch step.Type {
	case macrofile.LeftClick:
		fmt.Println("Left click!!")
		robotgo.MoveClick(step.Pos.X, step.Pos.Y, "left")
	case macrofile.RightClick:
		fmt.Println("Right click!!")
		robotgo.MoveClick(step.Pos.X, step.Pos.Y, "right")
	case macrofile.Test:
		fmt.Println("Test!!")
		colour := robotgo.GetPixelColor(step.Pos.X, step.Pos.Y)
		if colour == step.Colour {
			fmt.Println("-- PASS --")
		} else {
			fmt.Println("-- FAIL --")
		}
	}
	time.Sleep(step.Duration)
}

func main() {
	// Flag Setup
	cdPtr := flag.Bool("ShowCountDown", true, "Set this to false to skip playing a count down before starting the macro!\n'gmp.exe -ShowCountDown=flase'")
	flag.Parse()

	steps := macrofile.ReadMacro()

	// Play count down
	if *cdPtr {
		playStartSeq()
	}

	for _, step := range steps {
		playAction(step)
	}
}
