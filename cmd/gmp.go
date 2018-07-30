package main

import (
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
	playStartSeq()
}
