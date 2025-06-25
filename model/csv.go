package model

import (
	"bufio"
	"io"
	"strings"
)

func ToCsv(entries []Entry, writer io.Writer) error {
	header := []byte("project,tags,tools\n")
	_, err := writer.Write(header)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		line := []byte(entry.project)
		line = append(line, byte(','))

		for _, tag := range entry.tags {
			line = append(line, []byte(tag)...)
		}
		line = append(line, byte(','))

		for _, tool := range entry.tools {
			line = append(line, []byte(tool)...)
		}
		line = append(line, byte('\n'))

		_, err := writer.Write(line)
		if err != nil {
			return err
		}
	}

	return nil
}

func FromCsv(reader io.Reader) []Entry {
	var res []Entry
	scanner := bufio.NewScanner(reader)
	scanner.Scan() // jump the head
	for scanner.Scan() {
		entry := Entry{}
		line := scanner.Text()
		fields := strings.Split(line, ",")
		entry.SetProject(fields[0])

		var tags []string
		curTag := ""
		for _, c := range fields[1] {
			if c == '#' {
				if len(curTag) > 1 {
					tags = append(tags, curTag)
				}
				curTag = string(c)
			} else {
				curTag += string(c)
			}
		}
		if len(curTag) > 1 {
			tags = append(tags, curTag)
		}
		entry.tags = tags

		var tools []string
		curTool := ""
		for _, c := range fields[2] {
			if c == '&' {
				if len(curTool) > 1 {
					tools = append(tools, curTool)
				}
				curTool = string(c)
			} else {
				curTool += string(c)
			}
		}
		if len(curTool) > 1 {
			tools = append(tools, curTool)
		}
		entry.tools = tools

		res = append(res, entry)
	}

	return res
}
