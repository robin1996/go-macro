// Package mouse provides low level global hook for mouse input.
package mouse

import (
	"context"
	"syscall"
	"unsafe"

	"github.com/moutend/go-hook"
)

const (
	WM_LBUTTONDOWN = 0x0201
    WM_LBUTTONUP = 0x0202
    WM_MOUSEMOVE = 0x0200
    WM_MOUSEWHEEL = 0x020A
    WM_RBUTTONDOWN = 0x0204
    WM_RBUTTONUP = 0x0205
)

const (
	LeftClick = iota
	RightClick
)

// MOUSEHOOKSTRUCT corresponds to MOUSEHOOKSTRUCT structure.
// For more information, see the documentation on MSDN.
//
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms644968(v=vs.85).aspx
type MOUSEHOOKSTRUCT struct {
	hook.POINT
	MouseData   uint32
	Flags       uint32
	Time        uint32
	DWExtraInfo uint32
}

type MouseMessage struct {
	Button int
	hook.POINT
}

func notify(ctx context.Context, ch chan<- MouseMessage) {
	if ctx == nil {
		panic("hook/mouse: nil context")
	}
	if ch == nil {
		panic("hook/mouse: Notify using nil channel")
	}

	const WH_MOUSE_LL = 14
	var lResult hook.HHOOK
	hookProcedure := func(code, wParam, lParam uint64) uintptr {
		if (code >= 0) && (wParam == WM_LBUTTONDOWN) {
			m := *(*MOUSEHOOKSTRUCT)(unsafe.Pointer(uintptr(lParam)))
			mm := MouseMessage{LeftClick, m.POINT}
			ch <- mm
		}
		return uintptr(hook.CallNextHookEx(0, code, wParam, lParam))
	}

	go func() {
		lResult = hook.SetWindowsHookExW(
			WH_MOUSE_LL,
			hook.HOOKPROC(syscall.NewCallback(hookProcedure)),
			0,
			0)
		if lResult == 0 {
			panic("failed to set hook procedure")
		}
		var msg *hook.MSG
		hook.GetMessageW(&msg, 0, 0, 0)
		panic("hook finished")
	}()

	<-ctx.Done()
	if !hook.UnhookWindowsHookEx(lResult) {
		panic("failed to unhook")
	}
	return
}

// Notify causes package mouse to relay all keyboard events to ch.
func Notify(ctx context.Context, ch chan<- MouseMessage) {
	notify(ctx, ch)
}
