package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/kjhch/alfred-workflows/pkg/alfred"
)

func main() {
	wf := alfred.InitWorkflow()

	if len(wf.Input) > 0 {
		queryArgIp(wf)
	} else {
		queryIp(wf)
	}

	wf.SendOutput()
}

func queryIp(wf *alfred.Workflow) {
	localIpChan, publicIpChan, proxyIpChan := make(chan alfred.Item), make(chan alfred.Item), make(chan alfred.Item)
	go getLocalIp(localIpChan)
	go getPublicIp(publicIpChan)
	go getProxyIp(proxyIpChan)

	localIp, localIpOK := <-localIpChan
	publicIp, publicIpOK := <-publicIpChan
	proxyIp, proxyIpOK := <-proxyIpChan
	if localIpOK {
		wf.AddItem(localIp)
	}
	if publicIpOK {
		wf.AddItem(publicIp)
	}
	if proxyIpOK {
		wf.AddItem(proxyIp)
	}
}

func queryArgIp(wf *alfred.Workflow) {
	args := wf.Input
	itemChan := make(chan alfred.Item)
	for _, arg := range args {
		ip := arg
		go func() {
			ipInfo := getPublicIpInfo(ip)
			if ipInfo != nil {
				itemChan <- alfred.Item{
					Title:    ipInfo["IP"],
					Subtitle: fmt.Sprintf("%v  %v  %v", ipInfo["地址"], ipInfo["数据二"], ipInfo["数据三"]),
					Arg:      ipInfo["IP"],
					Icon:     &alfred.Icon{Path: "icon.png"},
				}
			} else {
				itemChan <- alfred.Item{
					Title: ip,
					Arg:   ip,
					Icon:  &alfred.Icon{Path: "icon.png"},
				}
			}

		}()
	}

	for i := 0; i < len(args); i++ {
		item := <-itemChan
		wf.AddItem(item)
	}
}

func getLocalIp(result chan<- alfred.Item) {
	defer close(result)
	interfaces, err := net.Interfaces()
	if err != nil {
		println(err.Error())
		return
	}

	// 遍历每个网络接口
	for _, iface := range interfaces {
		// 排除 loopback 接口和未启用的接口
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		// 获取接口的 IP 地址列表
		addrs, err := iface.Addrs()
		if err != nil {
			println(err.Error())
			continue
		}

		// 遍历接口的 IP 地址列表
		for _, addr := range addrs {
			// 判断是否为 IPv4 地址
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				result <- alfred.Item{
					Title:    "内网IP: " + ipnet.IP.String(),
					Subtitle: iface.Name + "  " + iface.HardwareAddr.String(),
					Arg:      ipnet.IP.String(),
					Icon:     &alfred.Icon{Path: "icon.png"},
				}
				return
			}
		}
	}
}

func getPublicIp(result chan<- alfred.Item) {
	defer close(result)
	ipInfo := getPublicIpInfo("")
	if ipInfo == nil {
		return
	}
	result <- alfred.Item{
		Title:    ipInfo["IP"],
		Subtitle: fmt.Sprintf("%v  %v  %v", ipInfo["地址"], ipInfo["数据二"], ipInfo["数据三"]),
		Arg:      ipInfo["IP"],
		Icon:     &alfred.Icon{Path: "icon.png"},
	}
}

func getProxyIp(result chan<- alfred.Item) {
	defer close(result)
	req, err := http.NewRequest(http.MethodGet, "https://www.google.com/search?q=ip", nil)
	if err != nil {
		println(err.Error())
		return
	}
	req.Header.Set("User-Agent", "curl/7.88.1")
	resp, err := (&http.Client{
		Timeout: 3 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		}}).Do(req)
	if err != nil {
		println(err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		println(err.Error())
		return
	}
	bodys := string(body)
	// println(bodys)
	re, err := regexp.Compile(`\(Client IP address: (.*)\)`)
	if err != nil {
		println(err.Error())
		return
	}
	matches := re.FindAllStringSubmatch(bodys, -1)
	for _, match := range matches {
		ip := match[1]
		ipInfo := getPublicIpInfo(ip)
		result <- alfred.Item{
			Title:    "代理IP: " + ip,
			Subtitle: fmt.Sprintf("%v  %v  %v", ipInfo["地址"], ipInfo["数据二"], ipInfo["数据三"]),
			Arg:      ip,
			Icon:     &alfred.Icon{Path: "icon.png"},
		}
	}
}

func getPublicIpInfo(ip string) map[string]string {
	req, err := http.NewRequest(http.MethodGet, "http://cip.cc/"+ip, nil)
	if err != nil {
		println(err.Error())
		return nil
	}
	req.Header.Set("User-Agent", "curl/7.88.1")
	resp, err := (&http.Client{Timeout: 2 * time.Second}).Do(req)
	if err != nil {
		println(err.Error())
		return nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		println(err.Error())
		return nil
	}
	bodys := string(body)
	lines := strings.Split(bodys, "\n")
	ipInfo := make(map[string]string)
	for _, line := range lines {
		if line == "" {
			continue
		}
		kv := strings.Split(line, ":")
		ipInfo[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}
	return ipInfo
}
