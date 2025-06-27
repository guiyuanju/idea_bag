package tui

import (
	"fmt"
	model "idea_bag/model"
	"idea_bag/parser"
	"os"
	"strings"

	"github.com/charmbracelet/x/term"
)

type Model struct {
	AllEntries      []model.Entry
	FilteredEntries []*model.Entry
	SelectedEntry   int
	Input           string
	ParsingEntry    model.Entry
	Msg             string
}

func initModel(entries []model.Entry) Model {
	filtered := make([]*model.Entry, len(entries))
	for i := range filtered {
		filtered[i] = &entries[i]
	}
	return Model{
		AllEntries:      entries,
		FilteredEntries: filtered,
		SelectedEntry:   0,
		Input:           "",
	}
}

func view(m Model) string {
	var s string
	// prompt
	s += fmt.Sprintf("> %s█\r\n", m.Input)
	// error
	if len(m.Msg) > 0 {
		s += fmt.Sprintf("%s\r\n", m.Msg)
	}
	// list
	for i, e := range m.FilteredEntries {
		cur := BgBlueMatched(e.String(), strings.Fields(m.Input))
		if m.SelectedEntry == i {
			cur = fmt.Sprintf("[ %s ]", cur)
			cur = Text(cur).Bold().String()
		} else {
			cur = fmt.Sprintf("  %s", cur)
		}
		s += cur + "\r\n"
	}
	return s
}

func filterEntry(m *Model) {
	text := m.Input
	tokens := strings.Fields(text)
	res := []*model.Entry{}
outer:
	for i := range len(m.AllEntries) {
		// reverse, make the newest added item shows at the top
		e := &m.AllEntries[len(m.AllEntries)-i-1]
		for _, token := range tokens {
			if !strings.Contains(e.String(), token) {
				continue outer
			}
		}
		res = append(res, e)
	}
	// any change to filter view reset the selected item
	if len(m.FilteredEntries) != len(res) {
		m.SelectedEntry = 0
	} else {
		for i := range len(res) {
			if res[i].Project != m.FilteredEntries[i].Project {
				m.SelectedEntry = 0
			}
		}
	}

	m.FilteredEntries = res
}

func instantParse(m *Model) error {
	p := parser.New(m.Input, 2)
	entry, err := p.Parse()
	if err != nil {
		m.Msg = err.Error()
		return err
	}
	m.ParsingEntry = entry
	m.Msg = ""
	return nil
}

type TUI struct {
	entries []model.Entry
}

func New(entries []model.Entry) TUI {
	return TUI{entries}
}

func (t *TUI) Run() {
	oldState, err := term.MakeRaw(os.Stdin.Fd())
	if err != nil {
		panic(err)
	}
	defer term.Restore(os.Stdin.Fd(), oldState)

	HideCursor()
	defer ShowCursor()

	model := initModel(t.entries)

	buf := make([]byte, 3)
	for {
		// clear the screen, not cross-platform!
		fmt.Print("\033[H\033[2J")
		filterEntry(&model)
		fmt.Print(view(model))

		n, err := os.Stdin.Read(buf)
		if err != nil {
			panic(err)
		}
		if n != 1 {
			continue
		}

		key := Key(buf[0])
		switch key {
		case KeyCtrlC:
			return
		case KeyCtrlU:
			model.Input = ""
			model.Msg = ""
		case KeyBackspace:
			if len(model.Input) > 0 {
				model.Input = model.Input[:len(model.Input)-1]
			}
			model.Msg = ""
		case KeyCtrlN:
			if len(model.FilteredEntries) > 0 {
				model.SelectedEntry = (model.SelectedEntry + 1) % len(model.FilteredEntries)
			}
		case KeyCtrlP:
			if len(model.FilteredEntries) > 0 {
				// add an extra length to ensure the mod result greater than zero
				model.SelectedEntry = (model.SelectedEntry - 1 + len(model.FilteredEntries)) % len(model.FilteredEntries)
			}
		case KeyEnter:
			if instantParse(&model) == nil {
				model.AllEntries = append(model.AllEntries, model.ParsingEntry)
				model.Input = ""
			}
		default:
			model.Input += string(byte(key))
			model.Msg = ""
		}
	}
}
