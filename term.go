package main

import (
	"time"

	ui "github.com/gizak/termui/v3"
)

const (
	duration = 500 * time.Millisecond
)

// Run run term
func Run() error {
	err := ui.Init()
	if err != nil {
		return err
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
				return nil
			case "?":
				tick.Stop()
			}
		}
	}
}

// Draw draw lyric
func Draw() {

}
