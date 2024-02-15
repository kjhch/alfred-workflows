package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kjhch/alfred-workflows/pkg/alfred"
)

func main() {
	wf := alfred.InitWorkflow()

	if len(wf.Input) > 0 {
		timeArg(wf)
	} else {
		timeNow(wf)
	}

	wf.SendOutput()
}

func timeNow(wf *alfred.Workflow) {
	now := time.Now()
	for _, item := range timeItems(now) {
		wf.AddItem(item)
	}
}

func timeArg(wf *alfred.Workflow) {
	arg := strings.TrimSpace(wf.Input[0])

	defaultFlag := false
	var items []alfred.Item
	if matched, _ := regexp.MatchString(`^\d+-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`, arg); matched {
		println(arg)
		t, _ := time.ParseInLocation(time.DateTime, arg, time.Local)
		items = timeItems(t)
	} else if matched, _ := regexp.MatchString(`^\d+-\d{2}-\d{2}$`, arg); matched {
		t, _ := time.ParseInLocation(time.DateOnly, arg, time.Local)
		items = timeItems(t)
	} else if matched, _ := regexp.MatchString(`^\d{10}$`, arg); matched {
		num, _ := strconv.ParseInt(arg, 10, 64)
		t := time.Unix(num, 0)
		items = timeItems(t)
	} else if matched, _ := regexp.MatchString(`^\d{13}$`, arg); matched {
		num, _ := strconv.ParseInt(arg, 10, 64)
		t := time.UnixMilli(num)
		items = timeItems(t)
	} else {
		items = timeItems(time.Now())
		defaultFlag = true

	}
	for _, item := range items {
		if defaultFlag {
			item.Subtitle = "未匹配到任意时间格式，默认显示当前时间"
		}
		wf.AddItem(item)
	}
}

func timeItems(t time.Time) []alfred.Item {
	t = t.In(time.Local)
	return []alfred.Item{
		{
			Title: fmt.Sprintf("秒: %v", t.Unix()),
			Arg:   fmt.Sprintf("%v", t.Unix()),
		},
		{
			Title: fmt.Sprintf("毫秒: %v", t.UnixMilli()),
			Arg:   fmt.Sprintf("%v", t.UnixMilli()),
		},
		{
			Title: fmt.Sprintf("日期: %v", t.Format(time.DateOnly)),
			Arg:   fmt.Sprintf("%v", t.Format(time.DateOnly)),
		},
		{
			Title: fmt.Sprintf("时间: %v", t.Format(time.DateTime)),
			Arg:   fmt.Sprintf("%v", t.Format(time.DateTime)),
		},
	}

}
