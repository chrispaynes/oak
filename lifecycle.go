// Package oak is a game engine...
package oak

import (
	"image"

	"bitbucket.org/oakmoundstudio/oak/dlog"
	"bitbucket.org/oakmoundstudio/oak/event"

	"golang.org/x/exp/shiny/screen"
)

var (
	winBuffer     screen.Buffer
	screenControl screen.Screen
	windowControl screen.Window

	windowRect     image.Rectangle
	windowUpdateCH = make(chan bool)

	osCh = make(chan func())

	lifecycleInit bool
)

//func init() {
//	runtime.LockOSThread()
//}
func lifecycleLoop(s screen.Screen) {
	if lifecycleInit {
		dlog.Error("Started lifecycle twice, aborting second call")
		return
	}
	lifecycleInit = true

	screenControl = s
	var err error

	// The window buffer represents the subsection of the world which is available to
	// be shown in a window.
	winBuffer, err = screenControl.NewBuffer(image.Point{ScreenWidth, ScreenHeight})
	if err != nil {
		dlog.Error(err)
		return
	}
	//defer winBuffer.Release()

	// The window controller handles incoming hardware or platform events and
	// publishes image data to the screen.
	changeWindow(ScreenWidth, ScreenHeight)
	//defer windowControl.Release()

	eb = event.GetBus()

	go keyHoldLoop()
	go inputLoop()
	go drawLoop()

	event.ResolvePending()
}

// runs f on the osLocked thread.
func osLockedFunc(f func()) {
	done := make(chan bool, 1)
	osCh <- func() {
		f()
		done <- true
	}
	<-done
}

func changeWindow(width, height int) {
	//windowFlag := windowControl != nil
	//if windowFlag {
	//	windowUpdateCH <- true
	//	windowControl.Publish()
	//	windowControl.Release()
	//}
	// The window controller handles incoming hardware or platform events and
	// publishes image data to the screen.
	wC, err := windowController(screenControl, width, height)
	if err != nil {
		dlog.Error(err)
		panic(err)
	}
	windowControl = wC
	windowRect = image.Rect(0, 0, width, height)

	//if windowFlag {
	//	eFilter = gesture.EventFilter{EventDeque: windowControl}
	//	windowUpdateCH <- true
	//}
}

// ChangeWindow sets the width and height of the game window. But it doesn't.
func ChangeWindow(width, height int) {
	//osLockedFunc(func() { changeWindow(width, height) })
	windowRect = image.Rect(0, 0, width, height)
}

// GetScreen returns the current screen as an rgba buffer
func GetScreen() *image.RGBA {
	return winBuffer.RGBA()
}
