package model

import (
	"fmt"
	"slices"
	"strings"
)

type Project = string
type Tag = string
type Tool = string

type Entry struct {
	Project Project
	tags    []Tag
	tools   []Tool
}

func (e *Entry) SetProject(s string) {
	e.Project = strings.TrimSpace(s)
}

func (e *Entry) AddTag(t string) {
	tag := "#"
	if t[0] == '#' {
		tag = ""
	}
	tag += t
	e.tags = append(e.tags, tag)
}

func (e *Entry) AddTool(t string) {
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

func prettyEntry(entry Entry) string {
	s := ""
	s += fmt.Sprintf("%s", entry.Project)
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

func (e Entry) Title() string {
	return e.Project
}

func (e Entry) Description() string {
	all := slices.Concat(e.tags, e.tools)
	if len(all) == 0 {
		return "no tags or tools"
	}
	return strings.Join(all, " ")
}

func (e *Entry) FilterValue() string {
	return e.Title() + " " + e.Description()
}
