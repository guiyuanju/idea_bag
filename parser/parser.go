package parser

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

type EntryParser struct {
	i      int
	s      string
	offset int
}

func (p *EntryParser) Parse() (string, error) {
	if len(p.s) == 0 {
		return "", errors.New(p.errorMsg("expect project name"))
	}

	e := []string{}
	for p.i < len(p.s) {
		switch p.s[p.i] {
		case '#':
			tag, err := p.parseTag()
			if err != nil {
				return "", errors.New(p.errorMsg(err.Error()))
			}
			e = append(e, tag)

		case ' ', '\n':
			p.i++
		default:
			part, err := p.parseProject()
			if err != nil {
				return "", errors.New(p.errorMsg(err.Error()))
			}
			e = append(e, strings.TrimSpace(part))
		}
	}

	return strings.Join(e, " "), nil
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

func (p *EntryParser) parseProject() (string, error) {
	project := []byte{}
	ends := []byte{'\n', '#'}
	for !p.isAtEnd() && !slices.Contains(ends, p.s[p.i]) {
		project = append(project, p.s[p.i])
		if !p.advance() {
			break
		}
	}
	return string(project), nil
}

func (p *EntryParser) parseTag() (string, error) {
	p.i++
	tag := []byte{'#'}
	ends := []byte{' ', '#', '\n'}
	for !p.isAtEnd() && !slices.Contains(ends, p.s[p.i]) {
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

func New(s string, offset int) EntryParser {
	return EntryParser{0, s, offset}
}
