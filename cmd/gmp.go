package main

import (
	"flag"
	"fmt"
	"time"
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

func main() {
	// Flag Setup
	cdPtr := flag.Bool("ShowCountDown", true, "Set this to false to skip playing a count down before starting the macro!\n'gmp.exe -ShowCountDown=flase'")
	flag.Parse()

	// Play count down
	if *cdPtr {
		playStartSeq()
	}
}
