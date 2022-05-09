package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"

	"github.com/dlclark/regexp2"
)

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

		// 有些程序输出的貌似不是正常的空行  与unix有关推测
		// 此时应使用 cat读取  使用 l参数会报错，无法排除空行
		if isMatch, _ := reg.MatchString(value); !isMatch {
			temp = append(temp, value)

		}
	}
	result = temp
	return result, nil
}

func UrlToIps(urls []string) ([]string, error) {
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
	_, err := IsIpOrDomainValid(temp)
	if err != nil {
		// fmt.Println(err)
		return nil, err
	}
	return temp, err

}
