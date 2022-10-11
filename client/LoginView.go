package main

import (
	tui "github.com/marcusolsson/tui-go"
)

type LoginHandler func(string)

type LoginView struct {
	tui.Box
	frame        *tui.Box
	name         *tui.Entry
	loginHandler LoginHandler
}

func NewLoginView() *LoginView {
	// https://github.com/marcusolsson/tui-go/blob/master/example/login/main.go
	// https://github.com/nqbao/go-sandbox/blob/chat/0.0.1/chatserver/tui/loginview.go
	view := &LoginView{}
	view.name = tui.NewEntry()
	view.name.SetFocused(true)
	view.name.SetSizePolicy(tui.Maximum, tui.Maximum)

	label := tui.NewLabel("Enter your name: ")
	view.name.SetSizePolicy(tui.Expanding, tui.Maximum)

	nameBox := tui.NewHBox(
		label,
		view.name,
	)
	nameBox.SetBorder(true)
	nameBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	view.frame = tui.NewVBox(
		tui.NewSpacer(),
		tui.NewPadder(-4, 0, tui.NewPadder(4, 0, nameBox)),
		tui.NewSpacer(),
	)
	view.Append(view.frame)

	return view
}

