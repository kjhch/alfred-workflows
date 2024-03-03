package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"testing"
)

func Test(t *testing.T) {
	os.Setenv("HTTP_PROXY", "socks5://localhost:1080")
	// 创建一个 HTTP 客户端
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}

	// 创建一个 GET 请求
	req, err := http.NewRequest("GET", "http://httpbin.org/get", nil)
	req.Header.Set("User-Agent", "curl/7.88.1")
	if err != nil {
		fmt.Println("创建请求失败:", err)
		return
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求失败:", err)
		return
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	// 打印响应状态码
	fmt.Println("响应:", string(b))
}

func TestRegReplace(t *testing.T) {
	s := regexp.MustCompile(`[\s|/]+`).ReplaceAllString("上海市 | 电信", "|")
	fmt.Println(s)
}
