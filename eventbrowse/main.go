package main

import (
	"fmt"
	"log"

	ui "github.com/gizak/termui"
	"github.com/gizak/termui/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	l := widgets.NewList()

	l.Rows = []string{
		"[0] github.com/gizak/termui",
		"[1] [你好，世界]",
		"[2] [こんにちは世界]",
		"[3] output",
		"[4] output.go",
		"[5] random_out.go",
		"[6] dashboard.go",
		"[7] foo",
		"[8] bar",
		"[9] baz",
	}

	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	l.TextStyle.Fg = ui.ColorWhite
	l.SelectedRowStyle.Bg = ui.ColorYellow

	l.Border = false
	width, height := ui.TerminalDimensions()
	l.SetRect(0, 0, width, height)

	ui.Render(l)

	previousKey := ""
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "<Enter>":
			fmt.Printf("selected: %v\n", l.Rows[l.SelectedRow])
			return
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			l.ScrollDown()
		case "k", "<Up>":
			l.ScrollUp()
		case "<C-d>":
			l.ScrollHalfPageDown()
		case "<C-u>":
			l.ScrollHalfPageUp()
		case "<C-f>":
			l.ScrollPageDown()
		case "<C-b>":
			l.ScrollPageUp()
		case "g":
			if previousKey == "g" {
				l.ScrollTop()
			}
		case "<Home>":
			l.ScrollTop()
		case "G", "<End>":
			l.ScrollBottom()
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		ui.Render(l)
	}
}
