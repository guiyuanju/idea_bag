package tui

import (
	"fmt"
	model "idea_bag/model"
	parser "idea_bag/parser"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
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
		SelectedEntry:   -1,
		Input:           "",
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	instantParse := func() {
		p := parser.New(m.Input, 2)
		e, err := p.Parse()
		if err != nil {
			m.Msg = err.Error()
		} else {
			m.Msg = ""
		}
		m.ParsingEntry = e
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			instantParse()
			m.AllEntries = append(m.AllEntries, m.ParsingEntry)
			m.Input = ""
			m.FilteredEntries = filterEntry(m.AllEntries, m.Input)
			return m, nil
		case "backspace":
			if len(m.Input) > 0 {
				m.Input = m.Input[:len(m.Input)-1]
			}
			m.FilteredEntries = filterEntry(m.AllEntries, m.Input)
			instantParse()
		case "ctrl+u":
			m.Input = ""
			m.FilteredEntries = filterEntry(m.AllEntries, m.Input)
			instantParse()
		default:
			m.Input += msg.String()
			m.FilteredEntries = filterEntry(m.AllEntries, m.Input)
			instantParse()
			return m, nil
		}
	}

	return m, nil
}

func filterEntry(entries []model.Entry, text string) []*model.Entry {
	res := []*model.Entry{}
	for i := range len(entries) {
		e := &entries[i]
		if strings.Contains(e.String(), strings.TrimSpace(text)) {
			res = append(res, e)
		}
	}
	return res
}

func (m Model) View() string {
	s := "> " + m.Input + "▋\n"
	if len(m.Msg) > 0 {
		s += m.Msg + "\n"
	}

	for _, e := range m.FilteredEntries {
		s += fmt.Sprintf("%v\n", e)
	}

	return s
}

type TUI struct {
	entries []model.Entry
}

func New(entries []model.Entry) TUI {
	return TUI{entries}
}

func (t *TUI) Run() {
	program := tea.NewProgram(initModel(t.entries))
	if _, err := program.Run(); err != nil {
		fmt.Println("ERROR: failed to run TEA: ", err)
		os.Exit(1)
	}
}
