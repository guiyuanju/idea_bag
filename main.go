package main

import (
	tui "idea_bag/tui"
	"os"
	"strings"
)

type UI interface {
	Run()
}

const FILE string = "ideabag.txt"

func main() {
	save := func(es []*string) {
		var entries []string
		for _, e := range es {
			entries = append(entries, *e)
		}
		os.WriteFile(FILE, []byte(strings.Join(entries, "\n")), 0644)
	}

	bs, err := os.ReadFile(FILE)
	if err != nil {
		panic(err)
	}
	entries := strings.Split(string(bs), "\n")
	var entryRefs []*string
	for _, e := range entries {
		if e == "" {
			continue
		}
		entryRefs = append(entryRefs, &e)
	}
	tui := tui.New(entryRefs, save)
	tui.Run()
}
