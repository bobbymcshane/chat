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
	parent   ContainerNavigator
	children []Container
}

func NewLayout() *Layout {
	return &Layout{}
}

func (layout *Layout) GetLayout() *Layout {
	return layout
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

func NewVerticalLayout() Container {
	return &VerticalLayout{NewLayout()}
}

func (layout *VerticalLayout) Draw(cells [][]termbox.Cell) {
	panic("unimplemented")
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

func NewHorizontalLayout() Container {
	return &HorizontalLayout{NewLayout()}
}

func (layout *HorizontalLayout) Draw(cells [][]termbox.Cell) {
	panic("unimplemented")
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
