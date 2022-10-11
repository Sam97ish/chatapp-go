package main

import (
	"fmt"
	"github.com/Sam97ish/chatapp-go/proto/service"
	"github.com/marcusolsson/tui-go"
	"time"
)

type ChatView struct {
	tui.Box
	chat    *tui.Box
	chatbox *tui.Box
	input   *tui.Entry
}

func NewChatView() *ChatView {
	// https://github.com/marcusolsson/tui-go/blob/master/example/login/main.go
	// https://github.com/nqbao/go-sandbox/blob/chat/0.0.1/chatserver/tui/loginview.go
	cview := &ChatView{}
	cview.chatbox = tui.NewVBox()
	chatboxScroll := tui.NewScrollArea(cview.chatbox)
	chatboxScroll.SetAutoscrollToBottom(true)

	chatboxView := tui.NewVBox(chatboxScroll)
	chatboxView.SetBorder(true)

	cview.input = tui.NewEntry()
	cview.input.SetFocused(true)
	cview.input.SetSizePolicy(tui.Expanding, tui.Maximum)

	inputBox := tui.NewHBox(cview.input)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	cview.chat = tui.NewVBox(chatboxView, inputBox)
	cview.chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	return cview
}

func (c *ChatView) AddMessage(msg *service.Message) {

	msgTime, errTime := time.Parse(time.RFC1123Z, msg.Timestamp)
	if errTime != nil {
		fmt.Printf("error parsing time: %v", errTime)
	}

	c.chatbox.Append(tui.NewHBox(
		tui.NewLabel(msgTime.Format(time.Stamp)),
		tui.NewPadder(1, 0, tui.NewLabel(fmt.Sprintf("<%s>", msg.User.Name))),
		tui.NewLabel(msg.Content),
		tui.NewSpacer(),
	))
}
