package main

import (
	"errors"
	"strings"

	"github.com/dlclark/regexp2"
)

func IsIpOrDomainValid(ips []string) (valid bool, err error) {
	// 必须有 .
	// 域名 至少有一个 .
	// . 不可以出现在 开头或结尾
	// 不可以有 ..
	// IP必须有 分隔的 3 个 .

	// 上面过滤完成后，对漏网之鱼进行过滤
	// 域名中只能出现 -， 不可出现其他非法字符 空格以及其！、%、&等非法字符
	// - 连接符不能连续出现  --    可以连续出现   中文域名

	expr_ip := `\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`
	reg_ip, _ := regexp2.Compile(expr_ip, 0)

	expr_false_num := `\d{1,3}\.\d{1,3}`
	reg_false_num, _ := regexp2.Compile(expr_false_num, 0)

	expr_special_chars := `[\\\^\$\|\?\*\+\[\]\{\}\(\)!@#%&/<>~]`
	reg_special, _ := regexp2.Compile(expr_special_chars, 0)

	for _, value := range ips {

		if strings.HasPrefix(value, "-") || strings.HasSuffix(value, "-") {
			err = errors.New(value + " 不是正确的域名/IP/URL，请检查")
			// fmt.Println("pre-suffix- i catched it")
			break
		}
		// if strings.Contains(value, "--") {
		// 	// 竟然有连续-- 的域名，解析后是中文域名
		// 	// 暂时禁用这条规则
		// 	err = errors.New(value + "不是正确的域名/IP/URL，请检查")
		// 	// fmt.Println("-- i catched it")
		// 	break
		// }

		if isMatch, _ := reg_special.MatchString(value); isMatch {
			err = errors.New(value + "不是正确的域名/IP/URL，请检查")
			// fmt.Println("reg_special i catched it")
			break
		}
		if !strings.Contains(value, ".") || strings.Contains(value, "..") {
			// fmt.Println("value:", value)
			err = errors.New(value + "不是正确的域名/IP/URL，请检查")
			// fmt.Println(". .. i catched it")
			break
		}
		if strings.HasPrefix(value, ".") || strings.HasSuffix(value, ".") {
			err = errors.New(value + "不是正确的域名/IP/URL，请检查")
			// fmt.Println("pre-suffix . i catched it")
			break
		}
		if isMatch, _ := reg_ip.MatchString(value); isMatch {
			if strings.Count(value, ".") != 3 {
				err = errors.New(value + "不是正确的域名/IP/URL，请检查")
				// fmt.Println("not ip i catched it")
				break
			}
		} else if isMatch_false_num, _ := reg_false_num.MatchString(value); isMatch_false_num {
			// err = errors.New(value + "不是正确的域名/IP/URL，请检查")
			// // fmt.Println("not ip_little  i catched it")
			// break
			// 这个逻辑有点问题，暂时禁用
		}

	}
	return valid, err
}
