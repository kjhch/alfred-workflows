package main

import "github.com/kjhch/alfred-workflows/internal/alfred"

func main() {
	wf := new(alfred.Workflow)
	wf.AddItem(alfred.Item{
		Title:    "Title",
		Subtitle: "Subtitle",
	})
	wf.SendFeedBack()
}
