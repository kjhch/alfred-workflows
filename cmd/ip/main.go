package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/kjhch/alfred-workflows/pkg/alfred"
)

func main() {
	wf := alfred.InitWorkflow()

	localIpChan, publicIpChan := make(chan alfred.Item), make(chan alfred.Item)
	go getLocalIp(localIpChan)
	go getPublicIp(publicIpChan)

	localIp, localIpOK := <-localIpChan
	publicIp, publicIpOK := <-publicIpChan
	if localIpOK {
		wf.AddItem(localIp)
	}
	if publicIpOK {
		wf.AddItem(publicIp)
	}

	wf.SendOutput()

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
			// fmt.Println("获取接口地址失败:", err)
			continue
		}

		// 遍历接口的 IP 地址列表
		for _, addr := range addrs {
			// 判断是否为 IPv4 地址
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				result <- alfred.Item{
					Title: ipnet.IP.String(),
				}
				return
			}
		}
	}
	// result <- ""
}

func getPublicIp(result chan<- alfred.Item) {
	defer close(result)
	req, err := http.NewRequest(http.MethodGet, "http://cip.cc/", nil)
	if err != nil {
		println(err.Error())
		return
	}
	req.Header.Set("User-Agent", "curl/7.88.1")
	resp, err := (&http.Client{Timeout: 2 * time.Second}).Do(req)
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
	fmt.Println(bodys)
	lines := strings.Split(bodys, "\n")
	ipInfo := make(map[string]string)
	for _, line := range lines {
		if line == "" {
			continue
		}
		kv := strings.Split(line, ":")
		ipInfo[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}
	result <- alfred.Item{
		Title:    ipInfo["IP"],
		Subtitle: fmt.Sprintf("%v  %v  %v", ipInfo["地址"], ipInfo["数据二"], ipInfo["数据三"]),
	}
}

func getProxyIp(result chan<- alfred.Item) {

}
