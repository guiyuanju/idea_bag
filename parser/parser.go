package parser

import (
	"errors"
	"fmt"
	model "idea_bag/model"
	"slices"
	"strings"
)

type EntryParser struct {
	i      int
	s      string
	offset int
}

func (p *EntryParser) invalid() bool {
	return slices.Contains([]byte{',', '"'}, p.s[p.i])
}

func (p *EntryParser) Parse() (model.Entry, error) {
	if len(p.s) == 0 {
		return model.Entry{}, errors.New(p.errorMsg("expect project name"))
	}

	e := model.Entry{}
	for p.i < len(p.s) {
		switch p.s[p.i] {
		case '#':
			if len(e.Project) == 0 {
				return model.Entry{}, errors.New(p.errorMsg("should provide project name first"))
			}
			tag, err := p.parseTag()
			if err != nil {
				return model.Entry{}, errors.New(p.errorMsg(err.Error()))
			}
			e.AddTag(tag)

		case '&':
			if len(e.Project) == 0 {
				return model.Entry{}, errors.New(p.errorMsg("should provide project name first"))
			}
			tool, err := p.parseTool()
			if err != nil {
				return model.Entry{}, errors.New(p.errorMsg(err.Error()))
			}
			e.AddTool(tool)
		case ' ', '\n':
			p.i++
		default:
			if p.i == 0 {
				project, err := p.parseProject()
				if err != nil {
					return model.Entry{}, errors.New(p.errorMsg(err.Error()))
				}
				e.SetProject(project)
			} else {
				return model.Entry{}, errors.New(p.errorMsg("project name should comes first"))
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

func (p *EntryParser) parseProject() (model.Project, error) {
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

func (p *EntryParser) parseTag() (model.Tag, error) {
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

func (p *EntryParser) parseTool() (model.Tool, error) {
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

func New(s string, offset int) EntryParser {
	return EntryParser{0, strings.TrimSpace(s), offset}
}
