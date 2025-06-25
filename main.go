package main

import (
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

	entries := model.FromCsv(file)
	tui := tui.New(entries)
	tui.Run()
}
