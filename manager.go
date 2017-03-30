package main

import (
	//"fmt"
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

func sliceView(startX, startY, width, height int, cells [][]termbox.Cell) [][]termbox.Cell {
	//fmt.Printf("Splitting view [%v x %v] [%v x %v]\n", startX, startY, width, height)
	buffer := make([][]termbox.Cell, height)
	for i := 0; startY+i < height; i++ {
		buffer[i] = cells[startY+i][startX+width:]
	}
	return buffer
}

// draw the container in the set slice of cells
func drawHorizontalView(container *container, cells [][]termbox.Cell) {
	viewHeight := len(cells)
	viewWidth := len(cells[0])
	//fmt.Printf("Drawing horizontal view [%v x %v]\n", viewWidth, viewHeight)
	containerHeight := viewHeight
	remainder := 0
	if numContainers := len(container.containers); numContainers > 1 {
		// if we have two containers, we want one divider
		numDividers := numContainers - 1
		correctedHeight := viewHeight - numDividers
		containerHeight = correctedHeight / numContainers
		remainder = correctedHeight % numContainers
		// TODO: account for remainder
	} else {
		// TODO: draw buffer
		return
	}

	dividerY := 0
	for _, c := range container.containers {
		cHeight := containerHeight
		if remainder > 0 {
			// add one to the container width so we don't get the remainder all in the last container
			cHeight++
			remainder--
		}

		if dividerY > 0 {
			for x := 0; x < viewWidth; x++ {
				cells[dividerY][x] = termbox.Cell{HORIZONTAL_LINE, termbox.ColorWhite, termbox.ColorBlack}
			}
			// add one because we drew a divider
			dividerY++
		}
		c.draw(sliceView(0, dividerY, viewWidth, cHeight, cells))
		dividerY += cHeight
	}
}

// draw the container in the set slice of cells
func drawVerticalView(container *container, cells [][]termbox.Cell) {
	viewWidth := len(cells[0])
	viewHeight := len(cells)
	//fmt.Printf("Drawing vertical view [%v x %v]\n", viewWidth, viewHeight)
	containerWidth := viewWidth
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

		c.draw(sliceView(dividerX, 0, cWidth, viewHeight, cells))
		dividerX += cWidth
	}
}

func (container *container) draw(cells [][]termbox.Cell) {
	if container.orientation == horizontal {
		drawHorizontalView(container, cells)
	} else {
		drawVerticalView(container, cells)
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

func test() [][]termbox.Cell {
	x, y := 204, 52
	buffer := make([][]termbox.Cell, y)
	// Loop over the rows, slicing each row from the front of the remaining cells slice
	for i := range buffer {
		buffer[i] = make([]termbox.Cell, x)
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
		containers = append(containers, &container{horizontal, nil})
	}

	layout := container{vertical, containers}
	// TODO: draw windows here
	layout.draw(getBuffer())

	termbox.Flush()
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Ch {
			case 'q':
				return
			case 'n':
				layout.containers[0].containers = append(layout.containers[0].containers, &container{})
				termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
				layout.draw(getBuffer())
				termbox.Flush()
				break
			}
		case termbox.EventResize:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			layout.draw(getBuffer())
			termbox.Flush()
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
