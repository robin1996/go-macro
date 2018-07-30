package main

import (
	"context"
	"fmt"
	"sync"
	"time"
	"unsafe"

	sys "golang.org/x/sys/windows"

	"github.com/go-vgo/robotgo"
	"github.com/moutend/go-hook"
	"github.com/robin1996/go-macro/macrofile"
	"github.com/robin1996/go-macro/mouse"
)

type HotKey struct {
	Id        int
	Modifiers int
	KeyCode   int
}

type MSG struct {
	HWND   uintptr
	UINT   uintptr
	WPARAM int16
	LPARAM int64
	DWORD  int32
	hook.POINT
}

type TestMessage struct {
	hook.POINT
	Colour string
}

func hotKeyEvents(startStopChan chan bool, testChan chan TestMessage) {
	recording := false
	user32 := sys.MustLoadDLL("user32")
	defer user32.Release()

	regHotKey := user32.MustFindProc("RegisterHotKey")
	peekMsg := user32.MustFindProc("PeekMessageW")

	hotKeys := map[int16]*HotKey{
		1: &HotKey{1, 0, 0x78}, // F9 -- See https://docs.microsoft.com/en-us/windows/desktop/inputdev/virtual-key-codes
		2: &HotKey{2, 0, 0x79}, // F10
	}

	for _, v := range hotKeys {
		r1, _, err := regHotKey.Call(
			0, uintptr(v.Id), uintptr(v.Modifiers), uintptr(v.KeyCode))
		if r1 != 1 {
			fmt.Println("Failed to register", v, ", error:", err)
		}
	}

	for {
		var msg = &MSG{}
		peekMsg.Call(uintptr(unsafe.Pointer(msg)), 0, 0, 0, 1)

		switch id := msg.WPARAM; id {
		case 1:
			startStopChan <- true
			recording = !recording
		case 2:
			if recording {
				x, y := robotgo.GetMousePos()
				tmsg := TestMessage{
					hook.POINT{int32(x), int32(y)},
					robotgo.GetPixelColor(x, y),
				}
				testChan <- tmsg
			}
		}
	}
}

func main() {
	startStopChan := make(chan bool, 2)
	testChan := make(chan TestMessage, 1)

	go hotKeyEvents(startStopChan, testChan)

	for {
		var wg sync.WaitGroup
		var isInterrupted bool
		var steps []macrofile.Step

		ctx, cancel := context.WithCancel(context.Background())
		mouseChan := make(chan mouse.ActionMessage, 1)

		fmt.Println("Press F9 to start/stop recording, F10 to add a test.")

		<-startStopChan

		go func() {
			wg.Add(1)
			mouse.Notify(ctx, mouseChan)
			wg.Done()
		}()

		fmt.Println("Starting recording...")

		for {
			startTime := time.Now()
			if isInterrupted {
				cancel()
				break
			}
			select {
			case <-startStopChan:
				isInterrupted = true
			case l := <-testChan:
				fmt.Println(l.Colour, l.POINT.X, l.POINT.Y)
				step := macrofile.Step{3, macrofile.Point{int(l.POINT.X), int(l.POINT.Y)}, l.Colour, time.Since(startTime)}
				steps = append(steps, step)
			case k := <-mouseChan:
				fmt.Println(k.Action, k.POINT.X, k.POINT.Y)
				step := macrofile.Step{k.Action, macrofile.Point{int(k.POINT.X), int(k.POINT.Y)}, "", time.Since(startTime)}
				steps = append(steps, step)
			}
		}
		wg.Wait()
		macrofile.WriteMacro(steps)
		fmt.Println("done")
	}
}
