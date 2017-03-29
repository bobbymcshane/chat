package main

import (
	"github.com/nsf/termbox-go"
)

type orientation int

const (
	unspecified = iota
	horizontal  = iota
	vertical    = iota
)

type pane struct {
}

type container struct {
	orientation int // unspecified, horizontal, vertical
	containers  []*container
}

func sliceView(width, height int, cells [][]termbox.Cell) [][]termbox.Cell {
	buffer := make([][]termbox.Cell, height)
	for i := 0; i < height; i++ {
		buffer[i] = cells[i][width:]
	}
	return buffer
}

// draw the container in the set slice of cells
func drawView(container *container, cells [][]termbox.Cell) {
	viewWidth := len(cells[0])
	viewHeight := len(cells)
	containerWidth := 0
	remainder := 0
	if numContainers := len(container.containers); numContainers > 1 {
		// if we have two containers, we want one divider
		numDividers := numContainers - 1
		correctedWidth := viewWidth - numDividers
		containerWidth = correctedWidth / numContainers
		remainder = correctedWidth % numContainers
		// TODO: account for remainder
	} else {
		// TODO: draw buffer
		return
	}

	dividerX := 0
	for _, c := range container.containers {
		cWidth := containerWidth
		if remainder > 0 {
			// add one to the container width so we don't get the remainder all in the last container
			cWidth++
			remainder--
		}

		if dividerX > 0 {
			for y := 0; y < viewHeight; y++ {
				cells[y][dividerX] = termbox.Cell{VERTICAL_LINE, termbox.ColorWhite, termbox.ColorBlack}
			}
			// add one because we drew a divider
			dividerX++
		}

		drawView(c, sliceView(dividerX, viewHeight, cells))
		dividerX += cWidth
	}
}

// returns buffer cells[height][width]
func getBuffer() [][]termbox.Cell {
	width, height := termbox.Size()
	cells := termbox.CellBuffer()
	buffer := make([][]termbox.Cell, height)
	// Loop over the rows, slicing each row from the front of the remaining cells slice
	for i := range buffer {
		if width > len(cells) {
			panic(i)
		}
		buffer[i], cells = cells[:width], cells[width:]
	}

	return buffer
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	var containers []*container
	for i := 0; i < 3; i++ {
		containers = append(containers, &container{})
	}

	layout := container{vertical, containers}
	// TODO: draw windows here
	termbox.Flush()
	drawView(&layout, getBuffer())

	termbox.Flush()
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Ch {
			case 'q':
				return
			case 'n':
				layout.containers = append(layout.containers, &container{})
				termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
				drawView(&layout, getBuffer())
				termbox.Flush()
				break
			}
		case termbox.EventResize:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			drawView(&layout, getBuffer())
			termbox.Flush()
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
