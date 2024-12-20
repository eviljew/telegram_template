package tg_template

import (
	"bufio"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
)

type Template struct {
	Text     string          `json:"text"`
	Buttons  map[int]*Button `json:"buttons"`
	PhotoURL string          `json:"photo_url"`
}

type LangTemplates map[Lang]*Template

func New(name string, lng Lang, data []*Data) (*Template, error) {
	// очень важно, чтобы папка с шаблонами была в корне
	file, err := os.Open("templates/" + name + ".txt")
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Выделяю из шаблона текст по указанному языку, если он конечно указан
	text := string(b)
	if lng != No {
		text = ExtractReText(string(b), fmt.Sprintf(`(?s)<:%v>(.*?)</:%v>`, lng, lng))
	}

	tmpl := &Template{
		Text:    text,
		Buttons: make(map[int]*Button),
	}

	tmpl.setData(data)

	tmpl.setButtons()

	return tmpl, nil
}

func (t *Template) PrepareKeyboardMarkup() *tgbotapi.InlineKeyboardMarkup {
	if len(t.Buttons) == 0 {
		return nil
	}

	// сортирую кнопки по id
	btnKeys := make([]int, 0, len(t.Buttons))
	for id := range t.Buttons {
		btnKeys = append(btnKeys, id)
	}
	sort.Ints(btnKeys)

	keyboard := make(map[int][]tgbotapi.InlineKeyboardButton)
	for idx := range btnKeys {
		btn := t.Buttons[btnKeys[idx]]

		if btn.Visible == VisibleOff {
			continue
		}

		if _, ok := keyboard[btn.Row]; !ok {
			keyboard[btn.Row] = tgbotapi.NewInlineKeyboardRow()
		}

		switch btn.Type {
		case LinkType:
			keyboard[btn.Row] = append(keyboard[btn.Row], tgbotapi.NewInlineKeyboardButtonURL(btn.Text, btn.Data))
			break
		case ButtonType:
			keyboard[btn.Row] = append(keyboard[btn.Row], tgbotapi.NewInlineKeyboardButtonData(btn.Text, btn.Data))
			break
		}
	}

	// сортирую строки кнопок по id
	rowKeys := make([]int, len(keyboard))
	for idx := range keyboard {
		rowKeys = append(rowKeys, idx)
	}
	sort.Ints(rowKeys)

	kbdmkp := make([][]tgbotapi.InlineKeyboardButton, 0)

	for _, idx := range rowKeys {
		if len(keyboard[idx]) != 0 {
			kbdmkp = append(kbdmkp, keyboard[idx])
		}
	}

	if len(keyboard) != 0 {
		mkp := tgbotapi.NewInlineKeyboardMarkup(kbdmkp...)
		return &mkp
	} else {
		return nil
	}
}

func (t *Template) AddData(d ...*Data) {
	t.setData(d)
}

func (t *Template) AddButtons(b ...*Button) {
	t.AddButtons(b...)
}

func (t *Template) SetPhotoURL(photoURL string) {
	t.PhotoURL = photoURL
}

func (t *Template) SetButtonVisible(btnID int, visible Visible) {
	_, ok := t.Buttons[btnID]
	if ok {
		t.Buttons[btnID].Visible = visible
	}
}

func (t *Template) RemoveDataString(pattern string) {
	var result []string
	scanner := bufio.NewScanner(strings.NewReader(t.Text))

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, pattern) {
			result = append(result, line)
		}
	}

	t.Text = strings.Join(result, "\n")
}

func (t *Template) setButtons() {
	rawButtons := ExtractReTextArr(t.Text, `(?s)(<:buttons>.*?</:buttons>)`)
	for _, rawBtn := range rawButtons {
		t.Text = regexp.MustCompile(regexp.QuoteMeta(rawBtn)).ReplaceAllString(t.Text, "")

		btn := NewButton(DelReText(rawBtn, `(<:buttons>)|(</:buttons>)`))

		t.Buttons[btn.ID] = btn
	}
}

func (t *Template) setData(data []*Data) {
	for _, el := range data {
		t.Text = ReplaceReText(t.Text, el.Pattern, el.Replacement)
	}
}

func (t *Template) addButton(buttons []*Button) {
	for _, el := range buttons {
		t.Buttons[el.ID] = el
	}
}
