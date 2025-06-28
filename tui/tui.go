package tui

import (
	"fmt"
	model "idea_bag/model"
	"idea_bag/parser"
	"os"
	"strings"

	"golang.org/x/term"
)

type Model struct {
	AllEntries             []model.Entry
	FilteredEntries        []*model.Entry
	SelectedFilterEntryIdx int
	SelectedEntry          *model.Entry
	Input                  string
	ParsingEntry           model.Entry
	ErrMsg                 string
}

func (m Model) headHeight() int {
	inputHeight := 1
	var errMsgHeight int
	if len(m.ErrMsg) > 0 {
		errMsgHeight = 1
	}
	return inputHeight + errMsgHeight + 1
}

func (m Model) usableHeight() int {
	_, h, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		panic("ERROR: failed to obtain terminal size")
	}
	usableHeight := h - m.headHeight()
	return usableHeight
}

func (m Model) visibleFilteredEntries() []*model.Entry {
	h := m.usableHeight()
	if len(m.FilteredEntries) <= h {
		return m.FilteredEntries
	}

	if m.SelectedFilterEntryIdx >= h {
		return m.FilteredEntries[m.SelectedFilterEntryIdx-h : m.SelectedFilterEntryIdx]
	}

	return m.FilteredEntries
}

func initModel(entries []model.Entry) Model {
	filtered := make([]*model.Entry, len(entries))
	for i := range filtered {
		filtered[i] = &entries[i]
	}
	return Model{
		AllEntries:             entries,
		FilteredEntries:        filtered,
		SelectedFilterEntryIdx: 0,
		Input:                  "",
	}
}

func (m *Model) resetSelectedEntry() {
	if len(m.FilteredEntries) > 0 {
		m.SelectedEntry = m.FilteredEntries[0]
	}
	m.SelectedEntry = nil
}

func view(m Model) string {
	var s string
	// prompt
	s += fmt.Sprintf("> %s█\r\n", m.Input)
	// error
	if len(m.ErrMsg) > 0 {
		s += fmt.Sprintf("%s\r\n", m.ErrMsg)
	}
	// list
	for _, e := range m.visibleFilteredEntries() {
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
		m.resetSelectedEntry()
	} else {
		for i := range len(res) {
			if res[i].Project != m.FilteredEntries[i].Project {
				m.resetSelectedEntry()
				break
			}
		}
	}

	m.FilteredEntries = res
}

func instantParse(m *Model) error {
	p := parser.New(m.Input, 2)
	entry, err := p.Parse()
	if err != nil {
		m.ErrMsg = err.Error()
		return err
	}
	m.ParsingEntry = entry
	m.ErrMsg = ""
	return nil
}

type TUI struct {
	entries []model.Entry
}

func New(entries []model.Entry) TUI {
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
			model.ErrMsg = ""
		case KeyBackspace:
			if len(model.Input) > 0 {
				model.Input = model.Input[:len(model.Input)-1]
			}
			model.ErrMsg = ""
		case KeyCtrlN:
			if len(model.FilteredEntries) > 0 {
				model.SelectedFilterEntryIdx = (model.SelectedFilterEntryIdx + 1) % len(model.FilteredEntries)
			}
		case KeyCtrlP:
			if len(model.FilteredEntries) > 0 {
				// add an extra length to ensure the mod result greater than zero
				model.SelectedFilterEntryIdx = (model.SelectedFilterEntryIdx - 1 + len(model.FilteredEntries)) % len(model.FilteredEntries)
			}
		case KeyEnter:
			if instantParse(&model) == nil {
				model.AllEntries = append(model.AllEntries, model.ParsingEntry)
				model.Input = ""
			}
		default:
			model.Input += string(byte(key))
			model.ErrMsg = ""
		}
	}
}
