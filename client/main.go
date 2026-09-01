package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	connection, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error connecting to the server please try again...")
		os.Exit(1)
	}
	defer connection.Close()

	app := tview.NewApplication()

	// popup asking for the username
	nameForm := tview.NewForm().
		AddInputField("Username:", "", 20, nil, nil).
		AddButton("Join", nil)
	nameForm.SetBorder(true).SetTitle(" go-chat ")
	nameForm.SetFocus(0)

	nameForm.GetButton(0).SetSelectedFunc(func() {
		username := strings.TrimSpace(nameForm.GetFormItemByLabel("Username:").(*tview.InputField).GetText())
		if username == "" {
			return
		}
		if _, err := connection.Write([]byte(username + "\n")); err != nil {
			fmt.Fprintln(os.Stderr, "failed to join:", err)
			os.Exit(1)
		}
		runChat(app, connection, username)
	})

	if err := app.SetRoot(nameForm, true).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func runChat(app *tview.Application, connection net.Conn, username string) {
	view := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetChangedFunc(func() { app.Draw() })
	input := tview.NewInputField().
		SetPlaceholder("Type a message (exit to leave)").
		SetFieldWidth(0)

	var mu sync.Mutex
	print := func(s string) {
		mu.Lock()
		defer mu.Unlock()
		view.Write([]byte(s + "\n"))
	}
	// color join/leave notices yellow
	notice := func(s string) {
		print("[yellow]" + s + "[-]")
	}

	chat := layout(view, input)

	// receive loop: server messages are already "<name>: <text>"
	go func() {
		sc := bufio.NewScanner(connection)
		for sc.Scan() {
			line := sc.Text()
			if strings.HasSuffix(line, " joined the chat...") || strings.HasSuffix(line, " left the chat...") {
				notice(line)
			} else {
				print(line)
			}
		}
		// server closed the connection
		print("[red]disconnected from server[-]")
		app.Stop()
	}()

	input.SetDoneFunc(func(key tcell.Key) {
		if key != tcell.KeyEnter {
			return
		}
		msg := strings.TrimSpace(input.GetText())
		input.SetText("")
		if msg == "" {
			return
		}
		if msg == "exit" {
			connection.Close()
			app.Stop()
			return
		}
		if _, err := connection.Write([]byte(msg + "\n")); err != nil {
			print("[red]failed to send: " + err.Error() + "[-]")
			return
		}
		// the server doesn't echo your own message back, show it locally
		print(username + ": " + msg)
	})

	// brief "Welcome X!" popup, then back to the chat
	modal := tview.NewModal().
		SetText("[yellow]Welcome to go-chat " + username + "![-]").
		AddButtons([]string{"Chat"})
	app.SetRoot(modal, false)
	time.AfterFunc(1500*time.Millisecond, func() {
		app.SetRoot(chat, true).SetFocus(input)
	})
}

func layout(view *tview.TextView, input *tview.InputField) *tview.Flex {
	view.SetBorder(true).SetTitle(" go-chat ")
	input.SetBorder(true)
	compose := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(view, 0, 1, false).
		AddItem(input, 3, 0, true)
	return compose
}
