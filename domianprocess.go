package main

import (
	"fmt"
	"net"

	"github.com/dlclark/regexp2"
)

func DomainToIp(misc []string) (ips, faildoamin []string, failint []int, err error) {
	// 正则表达式排除ip，对域名进行处理
	expr := `\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`
	reg, _ := regexp2.Compile(expr, 0)
	// if err != nil {
	// 	return nil, nil, err
	// }
	// var tempIps []string
	for i, value := range misc {

		if isMatch, _ := reg.MatchString(value); !isMatch {
			// Domain to ip
			fmt.Println("正在解析域名：", value)
			ns, err := net.LookupHost(value)
			if err != nil {
				// return nil, err
				// fmt.Println
				// fmt.Println(err)
				faildoamin = append(faildoamin, value)
				failint = append(failint, i)
			} else {
				ips = append(ips, ns[0])
			}
		} else {
			ips = append(ips, value)
		}
	}
	// fmt.Println("返回的err是")
	// fmt.Println(err)
	return ips, faildoamin, failint, err
}
