package main

import (
	"bufio"
	"fmt"
	"net"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	input    *tview.InputField
	pages    *tview.Pages
	form     *tview.Form
	chatArea *tview.TextView
	name     string
)

func handleClient(c net.Conn) {
	reader := bufio.NewReader(c)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Lost connection")
			panic(err)
		}
		fmt.Fprintf(chatArea, "%s", msg)
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:1922")
	if err != nil {
		panic(err)
	}
	go handleClient(conn)

	app := tview.NewApplication()
	chatArea = tview.NewTextView().SetDynamicColors(true).SetChangedFunc(func() {
		app.Draw()
	})
	bar := tview.NewTextView().SetText("Press Ctrl-C to exit")
	bar.SetBackgroundColor(tcell.ColorRebeccaPurple)
	input = tview.NewInputField().SetLabel(":>").
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetDoneFunc(func(key tcell.Key) {
			fmt.Fprintf(conn, "[#5aff00]%s[-]: %s\n", name, input.GetText())
			input.SetText("")
		})

	flex := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(chatArea, 0, 1, false).AddItem(bar, 1, 1, false).AddItem(input, 1, 1, true)

	form = tview.NewForm().
		AddInputField("What's your name", "", 20, nil, nil).
		AddButton("OK", func() {
			pages.SwitchToPage("main")
			name = form.GetFormItem(0).(*tview.InputField).GetText()
			input.SetLabel(name + " :> ")
			fmt.Fprintf(conn, "[red]%s joined[-]\n", name)
		})
	form.SetBorder(true)
	form.SetRect(20, 10, 40, 7)

	pages = tview.NewPages()
	pages.AddPage("inputname", form, false, true)
	pages.AddPage("main", flex, true, false)
	if err := app.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
