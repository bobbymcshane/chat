package main

import (
	"github.com/nsf/termbox-go"
)

type orientation int

const (
	unspecified = iota
	horizontal  = iota
	vertical    = iota
	text        = iota
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

type Manager struct {
	topLayout   ContainerNavigator
	orientation orientation // unspecified, horizontal, vertical
	focused     ContainerNavigator
}

func NewManager() *Manager {
	man := &Manager{}
	man.topLayout = &VerticalLayout{&Layout{}}
	man.SetFocus(man.topLayout)
	return man
}

func (manager *Manager) Focused() ContainerNavigator {
	return manager.focused
}

func (manager *Manager) SetFocus(toFocus Container) {
	if manager.focused == nil {
		manager.focused = toFocus.(ContainerNavigator)
	} else {
		if !manager.focused.IsFocused() {
			panic("Our focused container is not focused")
		}
		manager.focused.UnFocus()
		manager.focused = toFocus.(ContainerNavigator)
	}
	manager.focused.Focus()
}

func (manager *Manager) Focus(d direction) {
	if manager.focused == nil {
		// if nothing is currently focused, focus the first thing
		manager.focused = manager.topLayout
		manager.topLayout.Focus()
		return
	}

	if !manager.focused.IsFocused() {
		panic("Our focused container is not focused")
	}

	switch d {
	case in:
		if i := manager.focused.In(); i != nil {
			manager.SetFocus(i)
		}
	case out:
		if o := manager.focused.Out(); o != nil {
			manager.SetFocus(o)
		}
	case down:
		if b := manager.focused.Below(); b != nil {
			manager.SetFocus(b)
		}
	case up:
		if a := manager.focused.Above(); a != nil {
			manager.SetFocus(a)
		}
	case left:
		if l := manager.focused.Left(); l != nil {
			manager.SetFocus(l)
		}
	case right:
		if r := manager.focused.Right(); r != nil {
			manager.SetFocus(r)
		}
	}
}

// create a new container inside of parent
func (manager *Manager) newContainer(o orientation) {
	parent := manager.focused
	numContainers := 1
	if len(parent.Children()) == 0 {
		// make two containers if we are making the first container in a container
		numContainers++
	}
	for i := 0; i < numContainers; i++ {
		var newContainer ContainerNavigator
		if o == text {
			numContainers = 1
			newContainer = NewTextLayout()
		} else if o == horizontal {
			newContainer = &HorizontalLayout{&Layout{}}
		} else {
			newContainer = &VerticalLayout{&Layout{}}
		}
		newContainer.SetParent(parent)
		parent.Append(newContainer)
	}
}

func (manager Manager) Draw(cells [][]termbox.Cell) {
	manager.topLayout.Draw(cells)
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

	var manager *Manager = NewManager()
	for i := 0; i < 3; i++ {
		manager.newContainer(vertical)
	}

	// focus first container at lowest level
	var toFocus Container = manager.topLayout.Children()[0]
	for ; len(toFocus.Children()) > 0; toFocus = toFocus.Children()[0] {
	}
	manager.SetFocus(toFocus)

	// TODO: draw windows here
	manager.Draw(getBuffer())

	termbox.Flush()
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Ch {
			case 'q':
				return
			case 'd':
				for {
					if manager.Focused().Parent() == nil {
						// deleted the last container. end program
						return
					}

					if len(manager.Focused().Parent().Children()) > 1 {
						break
					} else {
						// TODO: move this logic into the container Delete
						// delete parent container if we are deleting the only container in the parent
						manager.Focus(out)
					}
				}

				toDelete := manager.Focused()
				for i, c := range toDelete.Parent().Children() {
					if c == toDelete {
						// focus sibling container. NOTE: we guaranteed above we would have one
						var toFocus Container
						if i > 0 {
							toFocus = toDelete.Parent().Children()[i-1]
						} else if (i + 1) < len(toDelete.Parent().Children()) {
							toFocus = toDelete.Parent().Children()[i+1]
						} else {
							panic("If we are the last container in our parent we should be deleting our parent")
						}
						manager.SetFocus(toFocus)
						toDelete.Delete()
					}
				}
			case 'c':
				manager.newContainer(vertical)
			case 't':
				manager.newContainer(text)
			case 'i':
				manager.Focus(in)
			case 'o':
				manager.Focus(out)
			case 'h':
				manager.Focus(left)
			case 'j':
				manager.Focus(down)
			case 'k':
				manager.Focus(up)
			case 'l':
				manager.Focus(right)
			case 0:
				switch ev.Key {
				case termbox.KeyCtrlV:
					// TODO:
					_, ok := manager.Focused().(*VerticalLayout)
					if !ok {
						// convert to vertical layout
						newLayout := NewVerticalLayout(manager.Focused())
						manager.Focused().Overwrite(newLayout)
						if manager.Focused().Parent() == nil {
							// this is our top-level layout
							manager.topLayout = newLayout
						}
						manager.SetFocus(newLayout)
					}
				case termbox.KeyCtrlC:
					// TODO:
					_, ok := manager.Focused().(*HorizontalLayout)
					if !ok {
						// convert to horizontal layout
						newLayout := NewHorizontalLayout(manager.Focused())
						manager.Focused().Overwrite(newLayout)
						if manager.Focused().Parent() == nil {
							// this is our top-level layout
							manager.topLayout = newLayout
						}
						manager.SetFocus(newLayout)
					}
				}
			}
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			manager.Draw(getBuffer())
			termbox.Flush()

		case termbox.EventResize:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			manager.Draw(getBuffer())
			termbox.Flush()
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
