# whereip
查询指定IP/域名/URL的归属地  where my  ip  is

# 简介
支持管道符输入
支持参数
优化输入处理

# 用法
`echo baidu.com|whereip.exe`
`echo  183.160.114.27:9999|whereip.exe`
`cat url.txt|whereip.exe`
`whereip.exe  -l url.txt   -c CN   -s  ah`

![image](https://user-images.githubusercontent.com/48342077/166718695-609fd382-856d-4710-8ea9-eeeac8f3cf61.png)

# 参数
-c string
    要保存的IP归属国家 (default "CN")
-l string
    url文件路径 (default "url.txt")
-s string
    要保存的IP归属省份

输入非常随意
url.txt是示例
只要能解析的地址都可以

# 输出
输出到console，txt，csv
每次输出会覆盖之前的结果

代码里能找到初恋般的感觉
