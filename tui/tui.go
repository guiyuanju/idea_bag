package tui

import (
	"fmt"
	model "idea_bag/model"
	"idea_bag/parser"
	"os"
	"slices"
	"strings"

	"golang.org/x/term"
)

type Model struct {
	Entries         []*model.Entry
	Filtered        []*model.Entry
	Visible         []*model.Entry
	VisibleStartIdx int
	SelectedEntry   *model.Entry
	Input           string
	ParsingEntry    model.Entry
	Msg             string
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
	m.SelectedEntry = m.Entries[len(m.Entries)-1]
}

func (m *Model) DelEntry(e *model.Entry) {
	if e == nil {
		return
	}
	m.Entries = slices.DeleteFunc(m.Entries, func(x *model.Entry) bool { return e == x })
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

	// reset selected entry to the first entry if it is filtered out
	if idx := slices.Index(res, m.SelectedEntry); idx < 0 && len(res) > 0 {
		m.SelectedEntry = res[0]
	}

	m.Filtered = res
}

func (m *Model) AvailableLines() int {
	_, h, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		panic("failed to get term size")
	}
	return h - 4
}

func (m *Model) updateVisible() {
	h := m.AvailableLines()
	start := m.VisibleStartIdx
	end := start + h
	if end >= len(m.Filtered) {
		end = len(m.Filtered) - 1
	}

	idx, ok := m.IndexOfEntryInFiltered(m.SelectedEntry)
	if !ok {
		m.Visible = m.Filtered
		return
	}

	if idx > end {
		end = idx
		start = end - h
		if start < 0 {
			start = 0
		}
	}

	if idx < start {
		start = idx
		end = start + h
		if end >= len(m.Filtered) {
			end = len(m.Filtered) - 1
		}
	}

	m.VisibleStartIdx = start

	m.Visible = m.Filtered[start : end+1]
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
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

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
		case KeyCtrlD:
			model.DelEntry(model.SelectedEntry)
		default:
			model.Input += string(byte(key))
			model.Msg = ""
		}
	}
}
