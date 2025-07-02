package main

import (
	"bufio"
	"fmt"
	model "idea_bag/model"
	tui "idea_bag/tui"
	"os"
)

type UI interface {
	Run()
}

func main() {
	file, err := os.OpenFile("ideabag.csv", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		fmt.Println("ERROR: failed to open csv file: ", err)
		return
	}
	defer file.Close()

	save := func(es []*model.Entry) {
		file, err := os.OpenFile("ideabag.csv", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
		if err != nil {
			panic("failed to open file for saving")
		}
		defer file.Close()

		var entries []model.Entry
		for _, e := range es {
			entries = append(entries, *e)
		}
		buf := bufio.NewWriter(file)
		model.ToCsv(entries, buf)
		buf.Flush()
	}

	entries := model.FromCsv(bufio.NewReader(file))
	var entryRefs []*model.Entry
	for _, e := range entries {
		entryRefs = append(entryRefs, &e)
	}
	tui := tui.New(entryRefs, save)
	tui.Run()
}
