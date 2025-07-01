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
	Entries       []*model.Entry
	Filtered      []*model.Entry
	Visible       []*model.Entry
	SelectedEntry *model.Entry
	Input         string
	ParsingEntry  model.Entry
	Msg           string
}

func initModel(entries []*model.Entry) Model {
	return Model{
		Entries:       entries,
		SelectedEntry: nil,
		Input:         "",
	}
}

func (m *Model) Update() {
	m.updateFilter()
	m.updateVisible()
}

func (m *Model) AddEntry(e model.Entry) {
	m.Entries = append(m.Entries, &e)
}

func (m *Model) IndexOfEntryInFiltered(entry *model.Entry) (int, bool) {
	for i, e := range m.Filtered {
		if e == entry {
			return i, true
		}
	}
	return 0, false
}

func (m *Model) SelectNextEntry() {
	idx, ok := m.IndexOfEntryInFiltered(m.SelectedEntry)
	if !ok {
		if len(m.Filtered) > 0 {
			m.SelectedEntry = m.Filtered[0]
		}
		return
	}
	if idx+1 >= len(m.Filtered) {
		m.SelectedEntry = m.Filtered[0]
		return
	}
	m.SelectedEntry = m.Filtered[idx+1]
}

func (m *Model) SelectPrevEntry() {
	idx, ok := m.IndexOfEntryInFiltered(m.SelectedEntry)
	if !ok {
		if len(m.Filtered) > 0 {
			m.SelectedEntry = m.Filtered[0]
		}
		return
	}
	if idx-1 < 0 {
		m.SelectedEntry = m.Filtered[len(m.Filtered)-1]
		return
	}
	m.SelectedEntry = m.Filtered[idx-1]
}

func (m *Model) updateFilter() {
	text := m.Input
	tokens := strings.Fields(text)
	res := []*model.Entry{}
outer:
	for i := range len(m.Entries) {
		// reverse, make the newest added item shows at the top
		e := m.Entries[len(m.Entries)-i-1]
		for _, token := range tokens {
			if !strings.Contains(e.String(), token) {
				continue outer
			}
		}
		res = append(res, e)
	}

	m.Filtered = res
}

func (m *Model) updateVisible() {
	m.Visible = m.Filtered
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

func view(m Model) string {
	var s string
	// prompt
	s += fmt.Sprintf("> %s█\r\n", m.Input)
	// error
	if len(m.Msg) > 0 {
		s += fmt.Sprintf("%s\r\n", m.Msg)
	}
	// list
	for _, e := range m.Visible {
		cur := BgBlueMatched(e.String(), strings.Fields(m.Input))
		if m.SelectedEntry == e {
			cur = fmt.Sprintf("[ %s ]", cur)
			cur = Text(cur).Bold().String()
		} else {
			cur = fmt.Sprintf("  %s", cur)
		}
		s += cur + "\r\n"
	}
	return s
}

type TUI struct {
	entries []*model.Entry
}

func New(entries []*model.Entry) TUI {
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
		model.Update()
		// clear the screen, not cross-platform!
		fmt.Print("\033[H\033[2J")
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
			model.SelectNextEntry()
		case KeyCtrlP:
			model.SelectPrevEntry()
		case KeyEnter:
			if instantParse(&model) == nil {
				model.AddEntry(model.ParsingEntry)
				model.Input = ""
			}
		default:
			model.Input += string(byte(key))
			model.Msg = ""
		}
	}
}
