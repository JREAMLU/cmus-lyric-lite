package main

import (
	"time"

	ui "github.com/gizak/termui/v3"
)

const (
	duration = 500 * time.Millisecond
)

func main() {
	err := ui.Init()
	if err != nil {
		panic(err)
	}

	defer ui.Close()

	cmus := NewCmus()

	tick := time.NewTicker(duration)

	uiEvents := ui.PollEvents()
	for {
		select {
		case <-tick.C:
			Listen(cmus)
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			case "?":
				tick.Stop()
			}
		}
	}

}
