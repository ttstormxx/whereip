package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// var sleep_count int = 0

func Resend(url string, data string) (*http.Response, error) {
	// 重新构造body
	dstring := strings.NewReader(data)
	// 使用client发包

	// client := &http.Client{
	// 	Timeout: 30 * time.Second,
	// }
	client := &http.Client{
		// Timeout: 10 * time.Second,
		Timeout: 5 * time.Second,
	}
	// url = url
	req, err := http.NewRequest("POST", url, dstring)

	// 添加请求头
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Connection", "Close")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.116 Safari/537.36")

	if err != nil {
		// fmt.Println(err)
	}
	// resp, err := http.Post(url, contentType, dstring)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	// Resend 函数里不可以添加defer resp.Body.Close(),会报错
	// defer resp.Body.Close()
	return resp, err
}
func PostGetIpLocation(url string, ipsSlice []string) []Ip_location {

	var temp []string
	var resp *http.Response
	// var url string
	for _, value := range ipsSlice {
		value = "\"" + value + "\""
		temp = append(temp, value)
	}

	dataString := strings.Join(temp, ",")
	dataString = "[" + dataString + "]"

	resp, err := Resend(url, dataString)
	for err != nil {
		if strings.Contains(string(err.Error()), "Client.Timeout exceeded") {
			fmt.Println("超时，正在重新发包。。。")
			resp, err = Resend(url, dataString)
		} else {

			fmt.Println(err)
		}
	}
	defer resp.Body.Close()

	// 有时会出现502未访问成功，需重新发送
	// 需更换为for 502则一直发
	// 直接重发POST包，body会丢失，需要重新读取body数据
	for resp.StatusCode == 429 {
		fmt.Println("达到API限制了,程序休眠，将等待60秒之后重新发包。。。")
		// 休眠
		time.Sleep(60 * time.Second)
		// 重新构造body
		resp, err = Resend(url, dataString)
		for err != nil {
			if strings.Contains(string(err.Error()), "Client.Timeout exceeded") {
				fmt.Println("超时，正在重新发包。。。")
				resp, err = Resend(url, dataString)
			} else {

				fmt.Println(err)
			}
		}
		if resp.StatusCode == 200 {
			break
		}
	}
	for resp.StatusCode != 200 {
		fmt.Println("发包失败了，正在重新发送。。。")
		resp, err = Resend(url, dataString)
		for err != nil {
			if strings.Contains(string(err.Error()), "Client.Timeout exceeded") {
				fmt.Println("超时，正在重新发包。。。")
				resp, err = Resend(url, dataString)
			} else {

				fmt.Println(err)
			}
		}
		if resp.StatusCode == 200 {
			break
		}
	}

	// ReadAll只能读取一次，重复读取，内容为空
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}

	jsonbody := make([]Ip_location, 0)

	// 如果unmarshal出错，也应该重新发包
	err = json.Unmarshal([]byte(body), &jsonbody)
	if err != nil {
		// fmt.Println("回包出错，将重新发包。。。")
		if string(err.Error()) == "unexpected end of JSON input" {
			fmt.Println("服务器未完整回包，将重新发包。。。")
		}
		fmt.Println(err)
	}
	return jsonbody
}

func IpSlice(ipsarr []string, n int, url string) []Ip_location {

	var temp []Ip_location
	for i := 1; i <= n; i++ {

		// 似乎api的限制最低是15次每分钟就会被限制
		// if 100*i%4500 == 0
		// API限制时状态码是429
		// if 100*i%1500 == 0 {
		// 	fmt.Println("超过API限制，休眠60秒重新发包。。。")
		// 	sleep_count++
		// 	time.Sleep(65 * time.Second)
		// }
		fmt.Println("第 ", i, " 次发包：")
		slice := ipsarr[(i-1)*100 : 100*i]
		// fmt.Println(slice)
		// 发包操作
		sliceResult := PostGetIpLocation(url, slice)
		for _, value := range sliceResult {
			temp = append(temp, value)
		}

	}
	// 避免最后一次发空包
	if !(n*100 == len(ipsarr)) {
		slice := ipsarr[n*100:]
		// fmt.Println(slice)
		// 发包操作
		sliceResult := PostGetIpLocation(url, slice)
		for _, value := range sliceResult {
			temp = append(temp, value)
		}
	}

	return temp
}

func JsonIpLocation(url string, ips []string) []Ip_location {

	// 切片 最大值100
	var temp []Ip_location
	n := len(ips) / 100

	// 判断是否超过100个URL值，超过则调用IpSlice()
	// 不超过则直接调用 PostGetIpLocation()
	if n == 0 || len(ips) == 100 {

		temp = PostGetIpLocation(url, ips)

	} else {
		temp = IpSlice(ips, n, url)
	}

	return temp
}
