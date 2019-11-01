package main

import (
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
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

// DrawEmpty draw empty list
func DrawEmpty() {
	// DrawList([]string{"", "", "[no lyrics](fg:red)"}, 0)
}

// DrawList draw lyric
func DrawList(rows []string, cline int) {
	w, h := ui.TerminalDimensions()
	l := widgets.NewList()

	if len(rows) > 0 {
		l.Title = rows[0]
	}

	l.PaddingTop = 2
	l.WrapText = false
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.SetRect(0, 0, w, h)

	idx := 1
	if cline > (h/2)-1 {
		idx = cline - (h/2 - h/8)
	}
	if idx > len(rows)-1 {
		idx = len(rows) - 1
	}
	l.Rows = rows[idx:]

	ui.Render(l)
}
