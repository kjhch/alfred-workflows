package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/kjhch/alfred-workflows/pkg/alfred"
)

func main() {
	wf := alfred.InitWorkflow()
	result := make(chan alfred.Item)
	go juliangUsage(result)
	go xiequUsage(result)
	i1 := <-result
	wf.AddItem(i1)
	i2 := <-result
	wf.AddItem(i2)

	wf.SendOutput()
}

func juliangUsage(result chan<- alfred.Item) {
	resp, err := http.Get("http://v2.api.juliangip.com/dynamic/balance?trade_no=1822891249461127&sign=3bdf4d2a52bc2e711f7673cef31024b6")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	body := make(map[string]any)
	json.Unmarshal(b, &body)
	if body["data"] != nil {
		balance := body["data"].(map[string]any)["balance"]
		result <- alfred.Item{
			Title: fmt.Sprintf("巨量剩余: %v", balance),
			Arg:   fmt.Sprintf("%v", balance),
		}
	} else {
		result <- alfred.Item{
			Title: fmt.Sprintf("巨量剩余: %v", body["msg"]),
		}
	}
}

func xiequUsage(result chan<- alfred.Item) {
	resp, err := http.Get("http://op.xiequ.cn/ApiUser.aspx?act=suitdt&uid=132239&ukey=EA391AB26FBD72939F91B207D79FD4DD")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	body := make(map[string]any)
	err = json.Unmarshal(b, &body)
	if err != nil {
		result <- alfred.Item{
			Title: fmt.Sprintf("携趣剩余: %v", string(b)),
		}
	}
}
