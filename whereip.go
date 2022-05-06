package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dlclark/regexp2"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	// "time"
)

type Ip_location struct {
	Status      string `json:"status"`
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	Region      string `json:"region"`
	RegionName  string `json:"regionName"`
	City        string `json:"city"`
	Query       string `json:"query"`
}

type Url_Ip_Location struct {
	Ip       string
	Url      string
	Location Ip_location
}

func ReadLine(filename string) ([]string, error) {

	var result []string
	// pip begins

	stat, _ := os.Stdin.Stat()

	if (stat.Mode() & os.ModeCharDevice) == 0 {

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			result = append(result, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		// pipe ends
	} else {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer func() {
			f.Close()
			// fmt.Println("文件关闭成功")
		}()
		reader := bufio.NewReader(f)

		for {
			// 这里文本最后一行读不到，需要处理,已处理
			// 或者最后一行置空
			// line,err :=buf.ReadString('\n')
			line, _, err := reader.ReadLine()

			if err != nil {
				if err == io.EOF {
					// 读取结束，报EOF

					// fmt.Println("读取结束")
					break
				}
				return nil, err
			}
			linestr := string(line)
			result = append(result, linestr)
		}
	}
	var temp []string
	for _, value := range result {

		// 处理两头空白字符
		value = strings.TrimSpace(value)
		// 抛弃空行
		expr := `^$`
		reg, _ := regexp2.Compile(expr, 0)
		if isMatch, _ := reg.MatchString(value); !isMatch {
			temp = append(temp, value)

		}
	}
	result = temp
	return result, nil
}

func UrlToIps(urls []string) []string {
	// 处理URL为IP:PORT列表
	// 正则表达式稍有问题
	// URL后面如果没有/则无法匹配
	// http://123.123.123.123 无法匹配
	expr := `(?<=://).+?(?=/)`
	reg, _ := regexp2.Compile(expr, 0)

	// 对于URL后方无/的，主动添加/
	// http://123.123.123.123
	// 变为
	// http://123.123.123.123/
	// 同时针对单行为IP/域名的情况，主动添加http://xxxx/
	//
	expr_http_finder := `^(http://|https://)`
	reg2, _ := regexp2.Compile(expr_http_finder, 0)

	var temp []string
	for _, value := range urls {
		// 任何value，后方均➕/
		value = value + "/"
		// 查找开头是否为http://或https://，没有则加上

		if isMatch, _ := reg2.MatchString(value); !isMatch {
			value = "http://" + value
		}

		match, _ := reg.FindStringMatch(value)
		ipPort := match.String()
		ipPort = strings.Split(ipPort, ":")[0]
		// 处理 IP:PORT列表为IP列表
		temp = append(temp, ipPort)
	}
	return temp

}

func PostGetIpLocation(url, contentType string, ipsSlice []string) []Ip_location {

	var temp []string
	for _, value := range ipsSlice {
		value = "\"" + value + "\""
		temp = append(temp, value)
	}

	dataString := strings.Join(temp, ",")
	dataString = "[" + dataString + "]"

	dstring := strings.NewReader(dataString)
	resp, err := http.Post(url, contentType, dstring)

	// 有时会出现502未访问成功，需重新发送
	// 需更换为for 502则一直发
	// 发现重新发也没用，稍等1分钟即可
	if resp.StatusCode == 502 {
		// fmt.Println("现在超过api的限制了")
		// fmt.Println("程序将休眠1分钟，然后重新发送。。。")
		fmt.Println("502了，准备重新发包。。。")
		fmt.Println("目前重发过程有问题，重新执行程序即可")
		// time.Sleep(60 * time.Second)
		// 重新发包，body丢失，不知啥原因，暂时搁置
		// resp, err = http.Post(url, contentType, dstring)
	}

	if err != nil {
		fmt.Println("post failed, err:", err)
		// return err
	}
	defer resp.Body.Close()
	if err != nil {
	}

	body, err := ioutil.ReadAll(resp.Body)

	jsonbody := make([]Ip_location, 0)

	err = json.Unmarshal([]byte(body), &jsonbody)
	if err != nil {
		fmt.Println(err)
	}
	return jsonbody
}

func IpSlice(ipsarr []string, n int, url, contentType string) []Ip_location {

	var temp []Ip_location
	for i := 0; i <= n; i++ {
		if i == 0 {
			i += 1
			slice := ipsarr[(i-1)*100 : 100*i]
			// fmt.Println(slice)
			// 发包操作
			sliceResult := PostGetIpLocation(url, contentType, slice)
			for _, value := range sliceResult {
				temp = append(temp, value)
			}

		} else {
			slice := ipsarr[(i-1)*100 : 100*i]
			// fmt.Println(slice)
			// 发包操作
			sliceResult := PostGetIpLocation(url, contentType, slice)
			for _, value := range sliceResult {
				temp = append(temp, value)
			}
		}

	}
	// 避免最后一次发空包
	if !(n*100 == len(ipsarr)) {
		slice := ipsarr[n*100 : len(ipsarr)]
		// fmt.Println(slice)
		// 发包操作
		sliceResult := PostGetIpLocation(url, contentType, slice)
		for _, value := range sliceResult {
			temp = append(temp, value)
		}
	}
	return temp
}

func JsonIpLocation(url, contentType string, ips []string) []Ip_location {

	// 切片 最大值100
	var temp []Ip_location
	n := len(ips) / 100

	// 判断是否超过100个URL值，超过则调用IpSlice()
	// 不超过则直接调用 PostGetIpLocation()
	if n == 0 || len(ips) == 100 {

		temp = PostGetIpLocation(url, contentType, ips)

	} else {
		temp = IpSlice(ips, n, url, contentType)
	}

	return temp
}

func To_Url_Ip_Location(urls, ips []string, locations []Ip_location) []Url_Ip_Location {

	temp := []Url_Ip_Location{}
	for i := 0; i < len(urls); i++ {
		temp = append(temp, Url_Ip_Location{ips[i], urls[i], locations[i]})
	}
	return temp
}

func FilterForeign(urlIpLocation []Url_Ip_Location, save_country_code, save_region_code string) []Url_Ip_Location {
	var temp []Url_Ip_Location
	if save_region_code == "" {
		for _, value := range urlIpLocation {
			if value.Location.CountryCode == save_country_code {
				temp = append(temp, value)
			}
		}
	} else {
		for _, value := range urlIpLocation {
			if value.Location.CountryCode == save_country_code && value.Location.Region == save_region_code {
				temp = append(temp, value)
			}
		}
	}

	return temp
}

func DomainToIp(ips []string) []string {
	// 正则表达式排除ip，对域名进行处理
	expr := `\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`
	reg, err := regexp2.Compile(expr, 0)
	if err != nil {
		fmt.Println(err)
	}
	var tempIps []string
	for _, value := range ips {

		if isMatch, _ := reg.MatchString(value); !isMatch {
			// Domain to ip
			ns, err := net.LookupHost(value)
			if err != nil {
				fmt.Println(err)
			}
			tempIps = append(tempIps, ns[0])
		} else {
			tempIps = append(tempIps, value)
		}
	}
	return tempIps
}

func ChinaWriteToFile(chinaAll []Url_Ip_Location) bool {
	var tempString string
	var tempStringCsv string
	for _, value := range chinaAll {
		// fmt.Println(value)
		// 写入文件

		tempString += value.Url + "   " + value.Ip + "   " + value.Location.CountryCode + "   " + value.Location.RegionName + "\n"
		tempStringCsv += value.Url + "," + value.Ip + "," + value.Location.CountryCode + "," + value.Location.RegionName + "\n"
		// 写入txt
		err := ioutil.WriteFile("results.txt", []byte(tempString), 0666)
		if err != nil {
			fmt.Println(err)
		}
		// 写入CSV文件
		err1 := ioutil.WriteFile("results.csv", []byte(tempStringCsv), 0666)
		if err1 != nil {
			fmt.Println(err1)
		}
	}
	return true
}

func main() {
	// 添加 google搜索结果链接爬虫  未完成 #放弃
	// 读URL列表    最后一行读不到，需要一个空行  @已处理
	// 处理URL列表为IP/域名列表   处理不能处理后方无/的URL     @已处理
	// 免费API每分钟最大发包次数45次 免费查询每次发包最大IP值是100 注意切割  切片 @已处理
	// 解析域名，本api无法通过POST解析域名  @已处理
	// 返回CN ip和对应URL struct   @已处理
	// 返回502 重新发包     未完成 #放弃
	// 筛选 fail  success    @已处理
	// 打印失败的查询 到控制台  @已处理
	// 查询失败的地址 寻找其他API   暂时搁置(放弃)
	// 将中国链接写入中国链接.txt  @已处理
	// 增加IP URL检测，是否需要进行split  @已处理
	// 针对上一条，思路不对，如果是IP，主动转换为URL，再进行处理
	// 在UrlToIps函数中连带后方无/的问题，一并处理
	// 当数据不是正确的URL/域名/IP时，NS解析会报错   暂时搁置(放弃)
	// CDN？  我看算了吧  #放弃
	// URL/IP/Domain去重？ 不去重  #放弃
	// 参数模式		@已处理
	// 管道符模式    @已处理

	url_file_path := flag.String("l", "url.txt", "url文件路径")

	save_country_code := flag.String("c", "CN", "要保存的IP归属国家")

	save_region_code := flag.String("s", "", "要保存的IP归属省份")

	flag.Parse()

	// *url_file_path = strings.ToUpper(*url_file_path)
	*save_country_code = strings.ToUpper(*save_country_code)
	*save_region_code = strings.ToUpper(*save_region_code)

	urls, err := ReadLine(*url_file_path)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("url总数为: ", len(urls), " 个")

	// 处理URL为IP:PORT列表-->IP列表
	misc := UrlToIps(urls)

	// 解析域名，本api无法通过POST解析域名
	// DomainToIp(ips)
	ips := DomainToIp(misc)

	// API
	// en Api
	// url := "http://ip-api.com/batch?lang=en"
	// CN Api
	url := "http://ip-api.com/batch?lang=zh-CN"
	contentType := "application/json"

	locations := JsonIpLocation(url, contentType, ips)

	myAll := To_Url_Ip_Location(urls, misc, locations)

	fmt.Println("全部地址归属如下：")
	for _, value := range myAll {
		fmt.Println(value)
	}
	fmt.Println("\n\n")
	// fmt.Println("其中中国IP的地址总数为： ",len(chinaAll)," 个")
	// fmt.Println("其中中国IP的地址为：")
	chinaAll := FilterForeign(myAll, *save_country_code, *save_region_code)

	fmt.Println("指定保存的IP归属 国家 是：", *save_country_code)
	if *save_region_code != "" {
		fmt.Println("指定保存的IP归属 省份代码 是：", *save_region_code)
	} else {
		fmt.Println("指定保存的IP归属 省份代码 是： 未指定省份")
	}
	per := float32(len(chinaAll)) / float32(len(myAll))
	per = 100 * per
	fmt.Println("指定IP归属地的IP地址总数为： ", len(chinaAll), " 个")
	fmt.Println("占比:%", per)
	fmt.Println("指定的IP地址为：")
	for _, value := range chinaAll {
		fmt.Println(value)
	}

	fmt.Println("\n\n")
	// 打印查询失败的链接组

	// failCount := 0
	// var failCount int = 0
	var failArr []Url_Ip_Location
	for _, value := range myAll {
		if value.Location.Status == "fail" {
			// fmt.Println(value)
			failArr = append(failArr, value)

		}
	}

	fmt.Println("查询失败总数为： ", len(failArr), " 个")
	fmt.Println("查询失败的地址如下：")
	for _, value := range failArr {
		fmt.Println(value)
	}
	fmt.Println("\n\n")

	// 写入文件
	// fmt.Println("写入文件。。。")
	ChinaWriteToFile(chinaAll)
	// 写入文件
	fmt.Println("写入文件完毕。。。(txt,csv)")

}
