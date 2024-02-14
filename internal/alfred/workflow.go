package alfred

import (
	"encoding/json"
	"fmt"
)

type Workflow struct {
	Items []Item `json:"items"`
}

// Icon specifies the "icon" field of an Item.
type Icon struct {
	Type string `json:"type,omitempty"`
	Path string `json:"path,omitempty"`
}

// Mod specifies the values of an Item.Mods map for the "mods" object.
type Mod struct {
	Variables map[string]string `json:"variables,omitempty"`
	Valid     *bool             `json:"valid,omitempty"`
	Arg       string            `json:"arg,omitempty"`
	Subtitle  string            `json:"subtitle,omitempty"`
	Icon      *Icon             `json:"icon,omitempty"`
}

// Text specifies the "text" field of an Item.
type Text struct {
	Copy      string `json:"copy,omitempty"`
	Largetype string `json:"largetype,omitempty"`
}

// Item specifies the members of the "items" array.
type Item struct {
	Variables    map[string]string `json:"variables,omitempty"`
	UID          string            `json:"uid,omitempty"`
	Title        string            `json:"title"`
	Subtitle     string            `json:"subtitle,omitempty"`
	Arg          string            `json:"arg,omitempty"`
	Icon         *Icon             `json:"icon,omitempty"`
	Autocomplete string            `json:"autocomplete,omitempty"`
	Type         string            `json:"type,omitempty"`
	Valid        *bool             `json:"valid,omitempty"`
	Match        string            `json:"match,omitempty"`
	Mods         map[string]Mod    `json:"mods,omitempty"`
	Text         *Text             `json:"text,omitempty"`
	QuicklookURL string            `json:"quicklookurl,omitempty"`
}

func (wf *Workflow) AddItem(item Item) {
	wf.Items = append(wf.Items, item)
}
func (wf *Workflow) SendFeedBack() {
	wfJSON, _ := json.Marshal(wf)
	fmt.Println(string(wfJSON))
}
