package main

import (
	"fmt"

	"github.com/kjhch/alfred-workflows/pkg/alfred"
)

func main() {
	wf := alfred.InitWorkflow()
	wf.AddItem(alfred.Item{
		Title:    fmt.Sprint(wf.Input),
		Subtitle: "Subtitle",
	})
	wf.SendOutput()
}
