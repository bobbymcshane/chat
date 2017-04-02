package main

import (
	"github.com/nsf/termbox-go"
	//"log"
)

type Container interface {
	// return layout implementation struct
	GetLayout() *Layout

	Parent() Container
	SetParent(c Container)

	Children() []Container
	SetChildren(c []Container)

	// add c to Container
	Append(c Container)
	// remove c from Container
	Remove(c Container)
	Delete()

	// rendering
	Draw([][]termbox.Cell)
}

type ContainerNavigator interface {
	// navigation
	Container

	IsFocused() bool
	Focus()
	UnFocus()

	Above() Container
	AboveContainer(c Container) Container

	Below() Container
	BelowContainer(c Container) Container

	Left() Container
	LeftContainer(c Container) Container

	Right() Container
	RightContainer(c Container) Container

	In() Container
	Out() Container
}

// General Layout. Implements ContainerNavigator
type Layout struct {
	focused  bool
	parent   ContainerNavigator
	children []Container
}

func NewLayout() *Layout {
	return &Layout{}
}

func (layout *Layout) GetLayout() *Layout {
	return layout
}

func (layout *Layout) IsFocused() bool {
	return layout.focused
}

func (layout *Layout) Focus() {
	layout.focused = true
}

func (layout *Layout) UnFocus() {
	layout.focused = false
}

// TODO: how do I make this take an arbitrary number of containers?
func (layout *Layout) Append(c Container) {
	layout.children = append(layout.children, c)
}

func (layout *Layout) Remove(c Container) {
	children := layout.Children()
	for i, child := range children {
		if c.GetLayout() == child.GetLayout() {
			layout.SetChildren(append(children[:i], children[i+1:]...))
			return
		}
	}
}

func (layout *Layout) Delete() {
	if parent := layout.Parent(); parent != nil {
		parent.Remove(layout)
	}
}

func (layout *Layout) Parent() Container {
	return layout.parent
}

func (layout *Layout) SetParent(c Container) {
	layout.parent = c.(ContainerNavigator)
}

func (layout *Layout) Children() []Container {
	return layout.children
}

func (layout *Layout) SetChildren(children []Container) {
	layout.children = children
}

func (layout *Layout) Draw(cells [][]termbox.Cell) {
	panic("unimplemented")
}

func (layout *Layout) Above() Container {
	if parent := layout.Parent(); parent != nil {
		return parent.(ContainerNavigator).AboveContainer(layout)
	}
	return nil
}

func (layout *Layout) AboveContainer(c Container) Container {
	return nil
}

func (layout *Layout) Below() Container {
	if parent := layout.Parent(); parent != nil {
		return parent.(ContainerNavigator).BelowContainer(layout)
	}
	return nil
}

func (layout *Layout) BelowContainer(c Container) Container {
	return nil
}

func (layout *Layout) Left() Container {
	if parent := layout.Parent(); parent != nil {
		return parent.(ContainerNavigator).LeftContainer(layout)
	}
	return nil
}

func (layout *Layout) LeftContainer(c Container) Container {
	return nil
}

func (layout *Layout) Right() Container {
	if parent := layout.Parent(); parent != nil {
		return parent.(ContainerNavigator).RightContainer(layout)
	}
	return nil
}

func (layout *Layout) RightContainer(c Container) Container {
	return nil
}

func (layout *Layout) In() Container {
	if layout.children != nil {
		return layout.children[0]
	} else {
		return nil
	}
}

func (layout *Layout) Out() Container {
	return layout.parent
}

// Vertical Layout. Implements ContainerNavigator
type VerticalLayout struct {
	*Layout
}

func NewVerticalLayout() *VerticalLayout {
	return &VerticalLayout{NewLayout()}
}

func sliceView(startX, startY, width, height int, cells [][]termbox.Cell) [][]termbox.Cell {
	//fmt.Printf("Splitting view [%v x %v] [%v x %v]\n", startX, startY, width, height)
	buffer := make([][]termbox.Cell, height)
	for i := 0; i < height; i++ {
		buffer[i] = cells[startY+i][startX : startX+width]
	}
	return buffer
}

func (layout *VerticalLayout) Draw(cells [][]termbox.Cell) {
	viewWidth := len(cells[0])
	viewHeight := len(cells)
	//fmt.Printf("Drawing vertical view [%v x %v]\n", viewWidth, viewHeight)
	containerWidth := viewWidth
	remainder := 0
	if numContainers := len(layout.Children()); numContainers > 1 {
		// if we have two containers, we want one divider
		numDividers := numContainers - 1
		correctedWidth := viewWidth - numDividers
		containerWidth = correctedWidth / numContainers
		remainder = correctedWidth % numContainers
		// TODO: account for remainder
	} else if numContainers == 1 {
		layout.Children()[0].Draw(cells)
		return
	} else if numContainers == 0 {
		// TODO: draw buffer
		if layout.IsFocused() && viewHeight >= 1 {
			cells[0][0] = termbox.Cell{'*', termbox.ColorWhite, termbox.ColorBlack}
		} else {
			for y := range cells {
				for x, _ := range cells[y] {
					cells[y][x] = termbox.Cell{'v', termbox.ColorWhite, termbox.ColorBlack}
				}
			}
		}
		return
	}

	dividerX := 0
	for _, c := range layout.Children() {
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

		c.Draw(sliceView(dividerX, 0, cWidth, viewHeight, cells))
		dividerX += cWidth
	}
}

func (layout *VerticalLayout) AboveContainer(c Container) Container {
	// there are no containers above any children in a vertical layout
	// look for a container above this one
	return layout.Above()
}

func (layout *VerticalLayout) BelowContainer(c Container) Container {
	// there are no containers below any children in a vertical layout
	// look for a container below this one
	return layout.Below()
}

func (layout *VerticalLayout) LeftContainer(c Container) Container {
	for i, con := range layout.Children() {
		if con.GetLayout() == c.GetLayout() {
			if i > 0 {
				return layout.Children()[i-1]
			}
			// nothing left of c... TODO ask our parent
			return layout.Left()
		}
	}
	return nil
}

func (layout *VerticalLayout) RightContainer(c Container) Container {
	for i, con := range layout.Children() {
		if con.GetLayout() == c.GetLayout() {
			if i+1 < len(layout.Children()) {
				return layout.Children()[i+1]
			}
			// nothing right of c... TODO ask our parent
			return layout.Right()
		}
	}
	return nil
}

// Horizontal Layout. Implements ContainerNavigator
type HorizontalLayout struct {
	*Layout
}

func NewHorizontalLayout() *HorizontalLayout {
	return &HorizontalLayout{NewLayout()}
}

func (layout *HorizontalLayout) Draw(cells [][]termbox.Cell) {
	viewHeight := len(cells)
	viewWidth := len(cells[0])
	//fmt.Printf("Drawing horizontal view [%v x %v]\n", viewWidth, viewHeight)
	containerHeight := viewHeight
	remainder := 0
	if numContainers := len(layout.Children()); numContainers > 1 {
		// if we have two containers, we want one divider
		numDividers := numContainers - 1
		correctedHeight := viewHeight - numDividers
		containerHeight = correctedHeight / numContainers
		remainder = correctedHeight % numContainers
		// TODO: account for remainder
	} else if numContainers == 1 {
		layout.Children()[0].Draw(cells)
		return
	} else if numContainers == 0 {
		// TODO: draw buffer
		if layout.IsFocused() && viewHeight >= 1 {
			cells[0][0] = termbox.Cell{'*', termbox.ColorWhite, termbox.ColorBlack}
		} else {
			for y := range cells {
				for x := range cells[y] {
					cells[y][x] = termbox.Cell{'h', termbox.ColorWhite, termbox.ColorBlack}
				}
			}
		}
		return
	}

	dividerY := 0
	for _, c := range layout.Children() {
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
		c.Draw(sliceView(0, dividerY, viewWidth, cHeight, cells))
		dividerY += cHeight
	}
}

func (layout *HorizontalLayout) AboveContainer(c Container) Container {
	for i, con := range layout.Children() {
		if con.GetLayout() == c.GetLayout() {
			if i > 0 {
				return layout.Children()[i-1]
			}
			// nothing above c... TODO ask our parent
			return layout.Above()
		}
	}
	return nil
}

func (layout *HorizontalLayout) BelowContainer(c Container) Container {
	for i, con := range layout.Children() {
		if con.GetLayout() == c.GetLayout() {
			if i+1 < len(layout.Children()) {
				return layout.Children()[i+1]
			}
			// nothing below c... TODO ask our parent
			return layout.Below()
		}
	}
	return nil
}

func (layout *HorizontalLayout) LeftContainer(c Container) Container {
	// there are no containers left of any children in a horizontal layout
	// look for a container left of this one
	return layout.Left()
}

func (layout *HorizontalLayout) RightContainer(c Container) Container {
	// there are no containers right of any children in a horizontal layout
	// look for a container right of this one
	return layout.Right()
}
