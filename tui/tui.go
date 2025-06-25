package tui

import (
	"fmt"
	model "idea_bag/model"
	parser "idea_bag/parser"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Model struct {
	EntryList list.Model
	// FilteredEntries []*model.Entry
	// SelectedEntry int
	// Input         string
	TextInput    textinput.Model
	ParsingEntry model.Entry
	Msg          string
}

func initModel(entries []list.Item) Model {
	// filtered := make([]*model.Entry, len(entries))
	// for i := range filtered {
	// 	filtered[i] = &entries[len(entries)-i-1]
	// }
	ti := textinput.New()
	ti.Placeholder = "Search or add"
	ti.Focus()
	ti.CharLimit = 0
	ti.Width = 20

	entryList := list.New(entries, list.NewDefaultDelegate(), 0, 0)
	entryList.Title = "Idea Bag"

	// fmt.Println(entries[0].FilterValue())
	// os.Exit(0)

	return Model{
		EntryList: entryList,
		// FilteredEntries: filtered,
		// SelectedEntry: -1,
		// Input:     "",
		TextInput: ti,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.TextInput.Validate = func(s string) error {
		p := parser.New(s, 2)
		e, err := p.Parse()
		m.ParsingEntry = e
		return err
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			if m.TextInput.Err != nil {
				return m, nil
			}
			// m.AllEntries = append(m.AllEntries, m.ParsingEntry)
			// m.FilteredEntries = filterEntry(m.AllEntries, m.Input)
			m.TextInput.Reset()
			return m, nil
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.EntryList.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	// m.TextInput, cmd = m.TextInput.Update(msg)
	// if m.TextInput.Err != nil {
	// 	m.Msg = m.TextInput.Err.Error()
	// } else {
	// 	m.Msg = ""
	// }

	m.EntryList, cmd = m.EntryList.Update(msg)

	// m.FilteredEntries = filterEntry(m.AllEntries, m.TextInput.Value())
	// instantParse()
	return m, cmd
}

func filterEntry(entries []model.Entry, text string) []*model.Entry {
	tokens := strings.FieldsFunc(text, func(r rune) bool { return strings.Contains(" \t\n#&", string(r)) })
	res := []*model.Entry{}
outer:
	for i := range len(entries) {
		e := &entries[len(entries)-i-1]
		for _, token := range tokens {
			if !strings.Contains(e.String(), token) {
				continue outer
			}
		}
		res = append(res, e)
	}
	return res
}

func (m Model) View() string {
	// s := ""
	// s += m.TextInput.View() + "\n"
	// if len(m.Msg) > 0 {
	// 	s += m.Msg + "\n"
	// }

	s := docStyle.Render(m.EntryList.View())
	// for _, e := range m.FilteredEntries {
	// 	s += fmt.Sprintf("%v\n", e)
	// }

	return s
}

type TUI struct {
	entries []model.Entry
}

func New(entries []model.Entry) TUI {
	return TUI{entries}
}

func (t *TUI) Run() {
	var items []list.Item
	for i := range len(t.entries) {
		items = append(items, &t.entries[i])
	}

	program := tea.NewProgram(initModel(items))
	if _, err := program.Run(); err != nil {
		fmt.Println("ERROR: failed to run TEA: ", err)
		os.Exit(1)
	}
}
