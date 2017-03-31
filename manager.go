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

type direction int

const (
	up    = iota
	down  = iota
	left  = iota
	right = iota
	out   = iota
	in    = iota
)

type pane struct {
}

type container struct {
	orientation orientation // unspecified, horizontal, vertical
	focused     bool
	containers  []*container
	parent      *container
}

func (container *container) Focus(d direction) *container {
	if !container.focused {
		panic("container not focused")
	}

	toFocus := container
	switch d {
	case in:
		if len(container.containers) > 0 {
			toFocus = container.containers[0]
		}
	case out:
		if container.parent != nil {
			toFocus = container.parent
		}
	case down:
		for containerItr, parent := container, container.parent; parent != nil; containerItr, parent = parent, parent.parent {
			if containerItr == parent {
				panic("invalid node traversal")
			}
			if parent.orientation != horizontal {
				// all of theses containers are above or below us
				continue
			}
			for i, c := range parent.containers {
				if c == containerItr {
					// we have found our current container in the parent. pick the container to the right
					if i+1 < len(parent.containers) {
						toFocus = parent.containers[i+1]
						for ; len(toFocus.containers) > 0; toFocus = toFocus.containers[0] {
						}
						goto found
					}
				}
			}
		}
	case up:
		for containerItr, parent := container, container.parent; parent != nil; containerItr, parent = parent, parent.parent {
			if containerItr == parent {
				panic("invalid node traversal")
			}
			if parent.orientation != horizontal {
				// all of theses containers are above or below us
				continue
			}
			for i, c := range parent.containers {
				if c == containerItr {
					// we have found our current container in the parent. pick the container to the right
					if i > 0 {
						toFocus = parent.containers[i-1]
						for ; len(toFocus.containers) > 0; toFocus = toFocus.containers[len(toFocus.containers)-1] {
						}
						goto found
					}
				}
			}
		}
	case left:
		for containerItr, parent := container, container.parent; parent != nil; containerItr, parent = parent, parent.parent {
			if containerItr == parent {
				panic("invalid node traversal")
			}
			if parent.orientation != vertical {
				// all of theses containers are above or below us
				continue
			}
			for i, c := range parent.containers {
				if c == containerItr {
					// we have found our current container in the parent. pick the container to the right
					if i > 0 {
						toFocus = parent.containers[i-1]
						for ; len(toFocus.containers) > 0; toFocus = toFocus.containers[len(toFocus.containers)-1] {
						}
						goto found
					}
				}
			}
		}
	case right:
		for containerItr, parent := container, container.parent; parent != nil; containerItr, parent = parent, parent.parent {
			if containerItr == parent {
				panic("invalid node traversal")
			}
			if parent.orientation != vertical {
				// all of theses containers are above or below us
				continue
			}
			for i, c := range parent.containers {
				if c == containerItr {
					// we have found our current container in the parent. pick the container to the right
					if i+1 < len(parent.containers) {
						toFocus = parent.containers[i+1]
						for ; len(toFocus.containers) > 0; toFocus = toFocus.containers[0] {
						}
						goto found
					}
				}
			}
		}
	}

found:
	if toFocus != container {
		container.focused = false
		toFocus.focused = true
	}
	return toFocus
}

func sliceView(startX, startY, width, height int, cells [][]termbox.Cell) [][]termbox.Cell {
	//fmt.Printf("Splitting view [%v x %v] [%v x %v]\n", startX, startY, width, height)
	buffer := make([][]termbox.Cell, height)
	for i := 0; i < height; i++ {
		buffer[i] = cells[startY+i][startX : startX+width]
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
	} else if numContainers == 1 {
		container.containers[0].draw(cells)
		return
	} else if numContainers == 0 {
		// TODO: draw buffer
		if container.focused && viewHeight >= 1 {
			cells[0][0] = termbox.Cell{'*', termbox.ColorWhite, termbox.ColorBlack}
		} else {
			if container.orientation == horizontal {
				for y := range cells {
					for x := range cells[y] {
						cells[y][x] = termbox.Cell{'h', termbox.ColorWhite, termbox.ColorBlack}
					}
				}
			}
		}
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
	} else if numContainers == 1 {
		container.containers[0].draw(cells)
		return
	} else if numContainers == 0 {
		// TODO: draw buffer
		if container.focused && viewHeight >= 1 {
			cells[0][0] = termbox.Cell{'*', termbox.ColorWhite, termbox.ColorBlack}
		} else {
			if container.orientation == vertical {
				for y := range cells {
					for x, _ := range cells[y] {
						cells[y][x] = termbox.Cell{'v', termbox.ColorWhite, termbox.ColorBlack}
					}
				}
			}
		}
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

// create a new container inside of parent
func (parent *container) newContainer(o orientation) {
	numContainers := 1
	if len(parent.containers) == 0 {
		// make two containers if we are making the first container in a container
		numContainers++
	}
	for i := 0; i < numContainers; i++ {
		newContainer := &container{}
		newContainer.orientation = o
		newContainer.parent = parent
		parent.containers = append(parent.containers, newContainer)
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

	layout := container{orientation: vertical}
	for i := 0; i < 3; i++ {
		layout.newContainer(vertical)
	}
	// focus first container at lowest level
	focused := layout.containers[0]
	for ; len(focused.containers) > 0; focused = focused.containers[0] {
	}

	focused.focused = true

	// TODO: draw windows here
	layout.draw(getBuffer())

	termbox.Flush()
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Ch {
			case 'q':
				return
			case 'd':
				for {
					if focused.parent == nil {
						// deleted the last container. end program
						return
					}

					if len(focused.parent.containers) > 1 {
						break
					} else {
						// delete parent container if we are deleting the only container in the parent
						focused = focused.Focus(out)
					}
				}

				toDelete := focused
				for o := direction(up); toDelete == focused && o < in; o++ {
					focused = focused.Focus(o)
				}

				for i, c := range toDelete.parent.containers {
					if c == toDelete {
						if len(toDelete.parent.containers) > 1 {
							toDelete.parent.containers = append(toDelete.parent.containers[:i], toDelete.parent.containers[i+1:]...)
						} else {
						}
					}
				}
			case 'c':
				focused.newContainer(vertical)
			case 'i':
				focused = focused.Focus(in)
			case 'o':
				focused = focused.Focus(out)
			case 'h':
				focused = focused.Focus(left)
			case 'j':
				focused = focused.Focus(down)
			case 'k':
				focused = focused.Focus(up)
			case 'l':
				focused = focused.Focus(right)
			case 0:
				switch ev.Key {
				case termbox.KeyCtrlV:
					focused.orientation = vertical
				case termbox.KeyCtrlC:
					focused.orientation = horizontal
				}
			}
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			layout.draw(getBuffer())
			termbox.Flush()

		case termbox.EventResize:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			layout.draw(getBuffer())
			termbox.Flush()
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
