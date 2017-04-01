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

	// navigation
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

	// rendering
	Draw([][]termbox.Cell)
}

// General Layout
type Layout struct {
	parent   Container
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
	children := c.Children()
	for i, child := range children {
		if child == c {
			c.SetChildren(append(children[:i], children[i+1:]...))
		}
	}
}

func (layout *Layout) Delete() {
	layout.Parent().Remove(layout)
}

func (layout *Layout) Parent() Container {
	return layout.parent
}

func (layout *Layout) SetParent(c Container) {
	layout.parent = c
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
	return nil
}

func (layout *Layout) AboveContainer(c Container) Container {
	return nil
}

func (layout *Layout) Below() Container {
	return nil
}

func (layout *Layout) BelowContainer(c Container) Container {
	return nil
}

func (layout *Layout) Left() Container {
	return nil
}

func (layout *Layout) LeftContainer(c Container) Container {
	return nil
}

func (layout *Layout) Right() Container {
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

// Vertical Layout
type VerticalLayout struct {
	*Layout
}

func NewVerticalLayout() *VerticalLayout {
	return &VerticalLayout{NewLayout()}
}

func (layout *VerticalLayout) Draw(cells [][]termbox.Cell) {
	panic("unimplemented")
}

func (layout *VerticalLayout) Above() Container {
	return nil
}

func (layout *VerticalLayout) AboveContainer(c Container) Container {
	return nil
}

func (layout *VerticalLayout) Below() Container {
	return nil
}

func (layout *VerticalLayout) BelowContainer(c Container) Container {
	return nil
}

func (layout *VerticalLayout) Left() Container {
	if parent := layout.Parent(); parent != nil {
		return parent.LeftContainer(layout)
	}
	return nil
}

func (layout *VerticalLayout) LeftContainer(c Container) Container {
	for i, con := range layout.Children() {
		// TODO: maybe use reflection or something to figure out how to compare here?
		if con.GetLayout() == c.GetLayout() {
			if i > 0 {
				return layout.Children()[i-1]
			}
			// nothing to our left... TODO ask our parent
			break
		}
	}
	return nil
}

func (layout *VerticalLayout) Right() Container {
	if parent := layout.Parent(); parent != nil {
		return parent.RightContainer(layout)
	}
	return nil
}

func (layout *VerticalLayout) RightContainer(c Container) Container {
	for i, con := range layout.Children() {
		// TODO: maybe use reflection or something to figure out how to compare here?
		if con.GetLayout() == c.GetLayout() {
			if i+1 < len(layout.Children()) {
				return layout.Children()[i+1]
			}
			// nothing to our right... TODO ask our parent
			break
		}
		//log.Panicln(con, c)
	}
	return nil
}
