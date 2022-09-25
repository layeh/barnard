package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/kennygrant/sanitize"
	"layeh.com/barnard/uiterm"
	"layeh.com/gumble/gumble"
)

const (
	uiViewLogo        = "logo"
	uiViewTop         = "top"
	uiViewStatus      = "status"
	uiViewInput       = "input"
	uiViewInputStatus = "inputstatus"
	uiViewOutput      = "output"
	uiViewTree        = "tree"
)

func esc(str string) string {
	return sanitize.HTML(str)
}

func (b *Barnard) UpdateInputStatus(status string) {
	b.UiInputStatus.Text = status
	b.UiTree.Rebuild()
	b.Ui.Refresh()
}

func (b *Barnard) AddOutputLine(line string) {
	now := time.Now()
	b.UiOutput.AddLine(fmt.Sprintf("[%02d:%02d:%02d] %s", now.Hour(), now.Minute(), now.Second(), line))
}

func (b *Barnard) AddOutputMessage(sender *gumble.User, message string) {
	if sender == nil {
		b.AddOutputLine(message)
	} else {
		b.AddOutputLine(fmt.Sprintf("%s: %s", sender.Name, strings.TrimSpace(esc(message))))
	}
}

func (b *Barnard) OnVoiceToggle(ui *uiterm.Ui, key uiterm.Key) {
	b.ToggleVoice()
	ui.Refresh()
}

func (b *Barnard) ToggleVoice() {
	if b.UiStatus.Text == "  Tx  " {
		b.UiStatus.Text = " Idle "
		b.UiStatus.Fg = uiterm.ColorBlack
		b.UiStatus.Bg = uiterm.ColorWhite
		b.Stream.StopSource()
	} else {
		b.UiStatus.Fg = uiterm.ColorWhite | uiterm.AttrBold
		b.UiStatus.Bg = uiterm.ColorRed
		b.UiStatus.Text = "  Tx  "
		b.Stream.StartSource()
	}
}

func (b *Barnard) OnQuitPress(ui *uiterm.Ui, key uiterm.Key) {
	b.Client.Disconnect()
	b.Ui.Close()
}

func (b *Barnard) OnClearPress(ui *uiterm.Ui, key uiterm.Key) {
	b.UiOutput.Clear()
}

func (b *Barnard) OnScrollOutputUp(ui *uiterm.Ui, key uiterm.Key) {
	b.UiOutput.ScrollUp()
}

func (b *Barnard) OnScrollOutputDown(ui *uiterm.Ui, key uiterm.Key) {
	b.UiOutput.ScrollDown()
}

func (b *Barnard) OnScrollOutputTop(ui *uiterm.Ui, key uiterm.Key) {
	b.UiOutput.ScrollTop()
}

func (b *Barnard) OnScrollOutputBottom(ui *uiterm.Ui, key uiterm.Key) {
	b.UiOutput.ScrollBottom()
}

func (b *Barnard) OnFocusPress(ui *uiterm.Ui, key uiterm.Key) {
	active := b.Ui.Active()
	if active == uiViewInput {
		b.Ui.SetActive(uiViewTree)
	} else if active == uiViewTree {
		b.Ui.SetActive(uiViewInput)
	}
}

func (b *Barnard) OnTextInput(ui *uiterm.Ui, textbox *uiterm.Textbox, text string) {
	if text == "" {
		return
	}
	if b.Client != nil && b.Client.Self != nil {
		b.Client.Self.Channel.Send(text, false)
		b.AddOutputMessage(b.Client.Self, text)
	}
}

func (b *Barnard) OnUiInitialize(ui *uiterm.Ui) {
	ui.Add(uiViewLogo, &uiterm.Label{
		Text: " barnard ",
		Fg:   uiterm.ColorWhite | uiterm.AttrBold,
		Bg:   uiterm.ColorMagenta,
	})

	ui.Add(uiViewTop, &uiterm.Label{
		Fg: uiterm.ColorWhite,
		Bg: uiterm.ColorBlue,
	})

	b.UiStatus = uiterm.Label{
		Text: " Idle ",
		Fg:   uiterm.ColorBlack,
		Bg:   uiterm.ColorWhite,
	}
	ui.Add(uiViewStatus, &b.UiStatus)

	b.UiInput = uiterm.Textbox{
		Fg:    uiterm.ColorWhite,
		Bg:    uiterm.ColorBlack,
		Input: b.OnTextInput,
	}
	ui.Add(uiViewInput, &b.UiInput)

	b.UiInputStatus = uiterm.Label{
		Fg: uiterm.ColorBlack,
		Bg: uiterm.ColorWhite,
	}
	ui.Add(uiViewInputStatus, &b.UiInputStatus)

	b.UiOutput = uiterm.Textview{
		Fg: uiterm.ColorWhite,
		Bg: uiterm.ColorBlack,
	}
	ui.Add(uiViewOutput, &b.UiOutput)

	b.UiTree = uiterm.Tree{
		Generator: b.TreeItem,
		Listener:  b.TreeItemSelect,
		Fg:        uiterm.ColorWhite,
		Bg:        uiterm.ColorBlack,
	}
	ui.Add(uiViewTree, &b.UiTree)

	b.Ui.AddKeyListener(b.OnFocusPress, uiterm.KeyTab)
	b.Ui.AddKeyListener(b.OnVoiceToggle, uiterm.KeyF1)
	b.Ui.AddKeyListener(b.OnQuitPress, uiterm.KeyF10)
	b.Ui.AddKeyListener(b.OnClearPress, uiterm.KeyCtrlL)
	b.Ui.AddKeyListener(b.OnScrollOutputUp, uiterm.KeyPgup)
	b.Ui.AddKeyListener(b.OnScrollOutputDown, uiterm.KeyPgdn)
	b.Ui.AddKeyListener(b.OnScrollOutputTop, uiterm.KeyHome)
	b.Ui.AddKeyListener(b.OnScrollOutputBottom, uiterm.KeyEnd)

	b.AddOutputLine("                Welcome to Barnard                ")
	b.AddOutputLine("--------------------------------------------------")
	b.AddOutputLine("HELP:")
	b.AddOutputLine("F1       : Toggle Voice Transmission")
	b.AddOutputLine("CTRL+L   : Clear chat log")
	b.AddOutputLine("TAB      : Toggle focus between chat and user tree")
	b.AddOutputLine("Page Up  : Scroll chat up")
	b.AddOutputLine("Page Down: Scroll chat down")
	b.AddOutputLine("HOME     : Scroll chat to the top")
	b.AddOutputLine("END      : Scroll chat to the bottom")
	b.AddOutputLine("F10      : Quit")
	b.AddOutputLine("--------------------------------------------------")

	b.start()

	if b.StartTX {
		b.ToggleVoice()
	}

	if b.StartupChannel != "" {
		b.Client.Self.Move(b.Client.Channels.Find(b.StartupChannel))
	}
}

func (b *Barnard) OnUiResize(ui *uiterm.Ui, width, height int) {
	ui.SetBounds(uiViewLogo, 0, 0, 9, 1)
	ui.SetBounds(uiViewTop, 9, 0, width-6, 1)
	ui.SetBounds(uiViewStatus, width-6, 0, width, 1)
	ui.SetBounds(uiViewInput, 0, height-1, width, height)
	ui.SetBounds(uiViewInputStatus, 0, height-2, width, height-1)
	ui.SetBounds(uiViewOutput, 0, 1, width-20, height-2)
	ui.SetBounds(uiViewTree, width-20, 1, width, height-2)
}
