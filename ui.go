package main

import ui "github.com/gizak/termui"

type OnInputFn func(message string)
type OnCloseFn func()

type ChatWindow struct {
	messages      [][]string
	messageWindow *ui.Table
	scrolledRows  int
}

func (chat *ChatWindow) Height() int {
	numRows := chat.messageWindow.Height
	if chat.messageWindow.Border {
		// subtract two for the border
		numRows -= 2
	}

	if chat.messageWindow.Separator {
		// divide by two for our separators
		numRows /= 2
	}
	return numRows
}

func (chat *ChatWindow) renderMessages() {
	numRows := chat.Height()

	chat.messageWindow.FgColors = make([]ui.Attribute, numRows)
	chat.messageWindow.BgColors = make([]ui.Attribute, numRows)
	for i := 0; i < numRows; i++ {
		chat.messageWindow.FgColors[i] = chat.messageWindow.FgColor
		chat.messageWindow.BgColors[i] = chat.messageWindow.BgColor
	}
	numMessages := len(chat.messages)
	firstCut := numMessages - numRows - chat.scrolledRows
	if firstCut < 0 {
		firstCut = 0
	}

	chat.messageWindow.Rows = chat.messages[firstCut : numMessages-chat.scrolledRows]
	ui.Render(chat.messageWindow)
}

func (chat *ChatWindow) AddMessage(user string, message string) {
	row := []string{user, message}
	chat.messages = append(chat.messages, row)
	chat.renderMessages()
}

func (chat *ChatWindow) Start(onInput OnInputFn, onClose OnCloseFn) {
	if err := ui.Init(); err != nil {
		panic(err)
	}

	ui.Handle("/sys/kbd/<escape>", func(ui.Event) {
		// quit
		ui.StopLoop()
		onClose()
	})

	// try to make a table with one row per message
	chat.messages = [][]string{}
	chat.messageWindow = ui.NewTable()
	chat.messageWindow.Rows = chat.messages // type [][]string
	chat.messageWindow.FgColor = ui.ColorWhite
	chat.messageWindow.BgColor = ui.ColorDefault
	chat.messageWindow.Height = ui.TermHeight() - 3
	chat.messageWindow.Width = ui.TermWidth()
	chat.messageWindow.Y = 0
	chat.messageWindow.X = 0
	chat.messageWindow.Border = true
	chat.messageWindow.Separator = true

	ui.Render(chat.messageWindow)

	inputBox := ui.NewPar("")
	inputBox.Height = 3
	inputBox.Width = ui.TermWidth()
	inputBox.TextFgColor = ui.ColorWhite
	inputBox.BorderLabel = "Input"
	inputBox.BorderFg = ui.ColorCyan
	inputBox.Y = chat.messageWindow.Y + chat.messageWindow.Height
	ui.Render(inputBox)

	ui.Handle("/sys/kbd/<enter>", func(ui.Event) {
		onInput(inputBox.Text)

		inputBox.Text = ""
		ui.Render(inputBox)
	})

	ui.Handle("/sys/kbd/<up>", func(ui.Event) {
		// scroll up
		if chat.scrolledRows < (len(chat.messages) - chat.Height()) {
			chat.scrolledRows += 1
			chat.renderMessages()
		}
	})

	ui.Handle("/sys/kbd/<down>", func(ui.Event) {
		// scroll down
		chat.scrolledRows -= 1
		if chat.scrolledRows < 0 {
			chat.scrolledRows = 0
		}
		chat.renderMessages()
	})

	ui.Handle("/sys/kbd/C-8", func(ui.Event) {
		// backspace
		if length := len(inputBox.Text); length > 0 {
			inputBox.Text = inputBox.Text[:length-1]
			ui.Render(inputBox)
		}
	})

	ui.Handle("/sys/kbd/<space>", func(ui.Event) {
		inputBox.Text += " "
		ui.Render(inputBox)
	})

	ui.Handle("/sys/kbd/", func(e ui.Event) {
		event := e.Data.(ui.EvtKbd)
		if len(event.KeyStr) == 1 {
			inputBox.Text += event.KeyStr
			ui.Render(inputBox)
		}
	})

	go func() {
		ui.Loop()
	}()
}

func (chat *ChatWindow) Close() {
	ui.Close()
}
