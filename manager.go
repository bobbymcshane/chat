package main

import (
	"github.com/nsf/termbox-go"
)

type pane struct {
}

type container struct {
	width, height int
	panes         *[]pane
}

func drawView(panes []pane) {
	twidth, theight := termbox.Size()
	numDividers := 0
	paneWidth := 0
	if numPanes := len(panes); numPanes > 1 {
		// if we have two panes, we want one divider
		numDividers = numPanes - 1
		paneWidth = twidth / numPanes
	}

	for divider := 1; divider <= numDividers; divider++ {
		x := divider * paneWidth
		for y := 0; y < theight; y++ {
			termbox.SetCell(x, y, '|', termbox.ColorWhite, termbox.ColorBlack)
		}
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	panes := []pane{{}, {}}
	// TODO: draw windows here
	drawView(panes)

	termbox.Flush()
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Ch {
			case 'q':
				return
			case 'n':
				panes = append(panes, pane{})
				termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
				drawView(panes)
				termbox.Flush()
			}
		case termbox.EventResize:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			drawView(panes)
			termbox.Flush()
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
