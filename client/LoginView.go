package main

import (
	tui "github.com/marcusolsson/tui-go"
)

type LoginHandler func(string)

type LoginView struct {
	tui.Box
	root  *tui.Box
	input *tui.Entry
}

var logo = ` _______                               ______ __           __   
|   _   |.----.----.---.-.-----.-----.|      |  |--.---.-.|  |_ 
|       ||   _|  __|  _  |     |  -__||   ---|     |  _  ||   _|
|___|___||__| |____|___._|__|__|_____||______|__|__|___._||____|
                                                                `

func NewLoginView() *LoginView {
	// https://github.com/marcusolsson/tui-go/blob/master/example/login/main.go
	// https://github.com/nqbao/go-sandbox/blob/chat/0.0.1/chatserver/tui/loginview.go
	Lview := &LoginView{}
	Lview.input = tui.NewEntry()
	Lview.input.SetFocused(true)
	inputBox := tui.NewHBox(
		Lview.input,
	)
	status := tui.NewStatusBar("Ready.")
	window := tui.NewVBox(
		tui.NewPadder(10, 1, tui.NewLabel(logo)),
		tui.NewPadder(12, 0, tui.NewLabel("Welcome to ArcaneChat! Press Enter to continue.")),
		tui.NewPadder(12, 0, tui.NewLabel("ESC to Quit..")),
		tui.NewPadder(-4, 0, tui.NewPadder(4, 0, inputBox)),
	)
	window.SetBorder(true)

	wrapper := tui.NewVBox(
		tui.NewSpacer(),
		window,
		tui.NewSpacer(),
	)
	content := tui.NewHBox(tui.NewSpacer(), wrapper, tui.NewSpacer())

	Lview.root = tui.NewVBox(
		content,
		status,
	)

	return Lview

}
