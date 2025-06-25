package main

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Project = string
type Tag = string
type Tool = string

type Entry struct {
	project Project
	tags    []Tag
	tools   []Tool
}

func (e *Entry) addTag(t string) {
	tag := "#"
	if t[0] == '#' {
		tag = ""
	}
	tag += t
	e.tags = append(e.tags, tag)
}

func (e *Entry) addTool(t string) {
	tool := "&"
	if t[0] == '&' {
		tool = ""
	}
	tool += t
	e.tools = append(e.tools, tool)
}

func (e Entry) String() string {
	return prettyEntry(e)
}

type Model struct {
	AllEntries      []Entry
	FilteredEntries []*Entry
	SelectedEntry   int
	Input           string
	ParsingEntry    Entry
	Msg             string
}

func initModel(entries []Entry) Model {
	filtered := make([]*Entry, len(entries))
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
		p := newEntryParser(m.Input, 2)
		e, err := p.parse()
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

func filterEntry(entries []Entry, text string) []*Entry {
	res := []*Entry{}
	for i := range len(entries) {
		e := &entries[i]
		if strings.Contains(e.String(), strings.TrimSpace(text)) {
			res = append(res, e)
		}
	}
	return res
}

type EntryParser struct {
	i      int
	s      string
	offset int
}

func (p *EntryParser) invalid() bool {
	return slices.Contains([]byte{',', '"'}, p.s[p.i])
}

func (p *EntryParser) parse() (Entry, error) {
	e := Entry{}
	for p.i < len(p.s) {
		switch p.s[p.i] {
		case '#':
			tag, err := p.parseTag()
			if err != nil {
				return Entry{}, errors.New(p.errorMsg(err.Error()))
			}
			e.tags = append(e.tags, tag)
		case '&':
			tool, err := p.parseTool()
			if err != nil {
				return Entry{}, errors.New(p.errorMsg(err.Error()))
			}
			e.tools = append(e.tools, tool)
		case ' ', '\n':
			p.i++
		default:
			if p.i == 0 {
				project, err := p.parseProject()
				if err != nil {
					return Entry{}, errors.New(p.errorMsg(err.Error()))
				}
				e.project = project
			} else {
				return Entry{}, errors.New(p.errorMsg("project name should comes first"))
			}
		}
	}

	return e, nil
}

func (p *EntryParser) errorMsg(msg string) string {
	space := strings.Repeat(" ", p.i+p.offset)
	return fmt.Sprintf("%s^ %s", space, msg)
}

func (p *EntryParser) advance() bool {
	p.i++
	return p.i < len(p.s)
}

func (p *EntryParser) isAtEnd() bool {
	return p.i >= len(p.s)
}

func (p *EntryParser) parseProject() (Project, error) {
	project := []byte{}
	ends := []byte{'\n', '#', '&'}
	for !p.isAtEnd() && !slices.Contains(ends, p.s[p.i]) {
		if p.invalid() {
			return "", errors.New("invalid character")
		}
		project = append(project, p.s[p.i])
		if !p.advance() {
			break
		}
	}
	return string(project), nil
}

func (p *EntryParser) parseTag() (Tag, error) {
	p.i++
	tag := []byte{'#'}
	ends := []byte{' ', '#', '&', '\n'}
	for !p.isAtEnd() && !slices.Contains(ends, p.s[p.i]) {
		if p.invalid() {
			return "", errors.New("invalid character")
		}
		tag = append(tag, p.s[p.i])
		if !p.advance() {
			break
		}
	}
	if len(tag) < 2 {
		return "", errors.New("expect tag name")
	}
	return string(tag), nil
}

func (p *EntryParser) parseTool() (Tool, error) {
	p.i++
	tool := []byte{'&'}
	ends := []byte{' ', '#', '&', '\n'}
	for !p.isAtEnd() && !slices.Contains(ends, p.s[p.i]) {
		if p.invalid() {
			return "", errors.New("invalid character")
		}
		tool = append(tool, p.s[p.i])
		if !p.advance() {
			break
		}
	}
	if len(tool) == 1 {
		return "", errors.New("expect tool name")
	}
	return string(tool), nil
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

func newEntryParser(s string, offset int) EntryParser {
	return EntryParser{0, strings.TrimSpace(s), offset}
}

func main() {
	file, err := os.OpenFile("ideabag.csv", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("ERROR: failed to open csv file: ", err)
		return
	}
	defer file.Close()

	entries := fromCsv(file)

	program := tea.NewProgram(initModel(entries))
	if _, err := program.Run(); err != nil {
		fmt.Println("ERROR: failed to run TEA: ", err)
		os.Exit(1)
	}
}

func prettyEntry(entry Entry) string {
	s := ""
	s += fmt.Sprintf("%s", entry.project)
	for _, t := range entry.tags {
		s += fmt.Sprintf(" %s", t)
	}
	for _, t := range entry.tools {
		s += fmt.Sprintf(" %s", t)
	}
	return s
}

func prettyEntries(entries []Entry) string {
	s := ""
	for _, e := range entries {
		s += prettyEntry(e)
		s += "\n"
	}
	return s
}
