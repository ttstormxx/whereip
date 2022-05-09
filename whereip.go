package main

import (
	"flag"
	"fmt"
	"strings"
	"time"
	// "time"
)

func main() {
	// 添加 google搜索结果链接爬虫  未完成 #放弃
	// 读URL列表    最后一行读不到，需要一个空行  @已处理
	// 处理URL列表为IP/域名列表   处理不能处理后方无/的URL     @已处理
	// 免费API每分钟最大发包次数45次 免费查询每次发包最大IP值是100 注意切割  切片 @已处理
	// 解析域名，本api无法通过POST解析域名  @已处理
	// 返回CN ip和对应URL struct   @已处理
	// 返回502 重新发包     未完成  @已处理
	// 筛选 fail  success    @已处理
	// 打印失败的查询 到控制台  @已处理
	// 查询失败的地址 寻找其他API   暂时搁置(放弃)
	// 将中国链接写入中国链接.txt  @已处理
	// 增加IP URL检测，是否需要进行split  @已处理
	// 针对上一条，思路不对，如果是IP，主动转换为URL，再进行处理
	// 在UrlToIps函数中连带后方无/的问题，一并处理
	// 当数据不是正确的URL/域名/IP时，NS解析会报错   @已处理
	// CDN？  我看算了吧  #放弃
	// URL/IP/Domain去重？ 不去重  #放弃
	// 参数模式		@已处理
	// 管道符模式    @已处理

	// 域名/IP有效性验证   	@已处理
	// 非有效域名/IP 剔除并展示 其他正常处理    @已处理
	// 域名解析失败，使用新的解析API尝试解析
	// https://ip-api.com/docs/dns
	// 域名解析失败的剔除并记录展示   @已处理
	// 发包失败的重新发包  非 200 的重新发包  @已处理
	// POST包重发 body丢失问题    @已处理
	// 重新发包失败的剔除并记录  目前  一直发 直到成功  @已处理
	// 发包超时重新发包  5s 	@已处理
	// 竟然会收到不完整的response wtf？？？  这里需要重发  应该是超时太长时间了，需要限制超时，超过重发 @已处理
	// 平均超时太长是因为没有挂梯子，需要挂梯子
	// 超过API 每分钟45次 每次100个 限制的，等待重发，似乎在变化，有可能是15次，有可能是20多次 @已处理
	// 添加代理池，超出免费API切换代理请求
	// 并发请求API
	// CDN展示   增加CDN 域名展示，在通用查询的基础上，额外对于CDN解析出的IP逐个查询并提醒
	// 吸收多个域名解析接口，结果对比
	// 默认不保存 增加 o 参数  @已处理
	// 优化大批量查询的稳定性  @已处理
	// 增加proxy选项 -proxy

	var (
		originUrl         []string // 原始URL列表
		originUrl_no_fail []string // 原始URL列表单次去除解析失败的地址
		misc              []string //处理完初始URL后的IP/域名列表
		misc_no_fail      []string //处理完初始URL后的IP/域名列表单次去除解析失败的地址

		ips        []string                                        //域名解析完毕后的IP列表
		faildomain []string                                        //域名解析失败的域名列表
		failint    []int                                           //标记域名解析失败的域名index
		url        string   = "http://ip-api.com/batch?lang=zh-CN" //API地址
		// url := "http://ip-api.com/batch?lang=en"  //en API地址
		// contentType string = "application/json" //请求头
		o string
	)
	start := time.Now()

	url_file_path := flag.String("l", "url.txt", "url文件路径")

	save_country_code := flag.String("c", "CN", "要保存的IP归属国家")

	save_region_code := flag.String("s", "", "要保存的IP归属省份")
	flag.StringVar(&o, "o", "", "保存文件")

	flag.Parse()

	// *url_file_path = strings.ToUpper(*url_file_path)
	*save_country_code = strings.ToUpper(*save_country_code)
	*save_region_code = strings.ToUpper(*save_region_code)

	originUrl, err := ReadLine(*url_file_path)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("输入的url总数为: ", len(originUrl), " 个")

	// 处理URL为IP:PORT列表-->IP列表
	// fmt.Println("处理URL为misc")
	misc, err = UrlToIps(originUrl)
	if err != nil {
		fmt.Println(err)
		return
		// os.Exit(1)
	}

	// 解析域名，本api无法通过POST解析域名
	// DomainToIp(ips)
	fmt.Println("开始域名解析")
	ips, faildomain, failint, _ = DomainToIp(misc)
	fmt.Println()
	if len(faildomain) > 0 {
		fmt.Println("解析失败的域名总共：", len(faildomain), "个", "如下：", faildomain)
	}

	// 根据失败域名下标 去除对应失败域名
	originUrl_no_fail = RemoveFailDomain(originUrl, failint)
	misc_no_fail = RemoveFailDomain(misc, failint)

	fmt.Println()
	fmt.Println("开始获取IP的归属地")
	locations := JsonIpLocation(url, ips)

	myAll := To_Url_Ip_Location(originUrl_no_fail, misc_no_fail, locations)
	// fmt.Println("结构化地址归属成功")
	fmt.Println("全部地址归属如下：")
	for _, value := range myAll {
		fmt.Println(value)
	}
	fmt.Println()
	fmt.Println()

	chinaAll := FilterForeign(myAll, *save_country_code, *save_region_code)

	if *save_country_code == "CN" {
		fmt.Println("指定保存的IP归属 国家 是：", *save_country_code, "(不含港澳台)")
	} else {
		fmt.Println("指定保存的IP归属 国家 是：", *save_country_code)
	}

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

	fmt.Println()
	fmt.Println()
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
	// fmt.Println("\n\n")
	fmt.Println()
	fmt.Println()

	// 写入文件
	if o != "" {
		fmt.Println("文件写入中...")
		fmt.Println()
		// fmt.Println("文件名为：", o)

		// fmt.Println("写入文件。。。")
		ChinaWriteToFile(chinaAll, o)
		// 写入文件
		if strings.HasSuffix(o, ".txt") {

			fmt.Println("写入文件完毕。。。(txt)")
		} else if strings.HasSuffix(o, ".csv") {
			fmt.Println("写入文件完毕。。。(csv)")
		} else {
			fmt.Println("写入文件完毕。。。(txt,csv)")
		}
	}
	elapse := time.Since(start)
	fmt.Println("程序本次运行耗时：", elapse)
}
