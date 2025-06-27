package tui

import (
	"fmt"
	model "idea_bag/model"
	parser "idea_bag/parser"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type State int

const (
	Listing = iota
	Adding
)

type Model struct {
	State        State
	EntryList    list.Model
	TextInput    textinput.Model
	ParseErr     string
	ParsingEntry model.Entry
	keys         *listKeyMap
}

type listKeyMap struct {
	openAddInput    key.Binding
	confirmAddInput key.Binding
	closeAddInput   key.Binding
	backspace       key.Binding
}

func newListKeyMap() *listKeyMap {
	return &listKeyMap{
		openAddInput: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add item"),
		),
		confirmAddInput: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm"),
		),
		closeAddInput: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel add"),
		),
		backspace: key.NewBinding(
			key.WithKeys("backspace"),
			key.WithHelp("backspace", "delete"),
		),
	}
}

func initModel(entries []list.Item) Model {
	listKeys := newListKeyMap()

	ti := textinput.New()
	ti.Placeholder = "Add"
	ti.Blur()
	ti.CharLimit = 0
	ti.Width = 20

	entryList := list.New(entries, list.NewDefaultDelegate(), 0, 0)
	// entryList.Title = "Idea Bag"
	entryList.SetShowTitle(false)
	entryList.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			listKeys.openAddInput,
			listKeys.confirmAddInput,
			listKeys.closeAddInput,
		}
	}

	return Model{
		State:     Listing,
		EntryList: entryList,
		TextInput: ti,
		keys:      listKeys,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.TextInput.Validate = func(s string) error {
		p := parser.New(s, 2)
		e, err := p.Parse()
		m.ParsingEntry = e
		return err
	}

	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.EntryList.SetSize(msg.Width-h, msg.Height-v-2)
	case tea.KeyMsg:
		if m.EntryList.FilterState() == list.Filtering {
			break
		}
		switch {
		case key.Matches(msg, m.keys.openAddInput):
			if m.State == Listing {
				m.State = Adding
				m.TextInput.Focus()
				return m, nil
			}

		case key.Matches(msg, m.keys.closeAddInput):
			if m.State == Adding {
				m.State = Listing
				m.TextInput.Reset()
				m.TextInput.Blur()
				return m, nil
			}

		case key.Matches(msg, m.keys.confirmAddInput):
			if m.State == Adding && m.TextInput.Err == nil {
				m.EntryList.InsertItem(0, &m.ParsingEntry)
				m.State = Listing
				m.TextInput.Reset()
				m.TextInput.Blur()
				return m, nil
			}

		case key.Matches(msg, m.keys.backspace):
			if m.State == Listing {
				m.EntryList.RemoveItem(m.EntryList.GlobalIndex())
			}
		}
	}

	var cmd tea.Cmd
	var error string

	switch m.State {
	case Adding:
		m.TextInput, cmd = m.TextInput.Update(msg)
		cmds = append(cmds, cmd)
		if m.TextInput.Err != nil {
			error = m.TextInput.Err.Error()
		} else {
			error = ""
		}
		m.ParseErr = error

	case Listing:
		m.EntryList, cmd = m.EntryList.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
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
	input := ""
	if m.State == Adding {
		input += "  " + m.TextInput.View() + "\n"
		input += "  " + m.ParseErr
	}
	list := docStyle.Render(m.EntryList.View())
	res := lipgloss.JoinVertical(lipgloss.Left, input, list)
	return res
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
