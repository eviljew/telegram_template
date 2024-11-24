package tg_template

import (
	"strconv"
	"strings"
)

type Button struct {
	ID      int
	Row     int
	Text    string
	Type    Type
	Data    string
	Visible Visible
}

func NewButton(raw string) *Button {
	splitRaw := strings.Split(raw, ";")

	btn := &Button{
		ID:      0,
		Row:     0,
		Text:    "",
		Type:    "",
		Data:    "",
		Visible: VisibleOff,
	}

	if len(splitRaw) >= 1 {
		btn.ID, _ = strconv.Atoi(splitRaw[0])
	}

	if len(splitRaw) >= 2 {
		btn.Row, _ = strconv.Atoi(splitRaw[1])
	}

	if len(splitRaw) >= 3 {
		btn.Text = splitRaw[2]
	}

	if len(splitRaw) >= 4 {
		btn.Type = Type(splitRaw[3])
	}

	if len(splitRaw) >= 5 {
		btn.Data = splitRaw[4]
	}

	if len(splitRaw) >= 6 {
		btn.Visible = VisibleFromString(splitRaw[5])
	}

	return btn
}

type Type string

const (
	ButtonType Type = "button"
	LinkType   Type = "link"
)

type Visible int

const (
	VisibleOn  Visible = 1
	VisibleOff Visible = 0
)

func VisibleFromString(s string) Visible {
	intVisible, _ := strconv.Atoi(s)

	return Visible(intVisible)
}
