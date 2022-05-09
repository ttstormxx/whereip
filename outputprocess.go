package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

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

func ChinaWriteToFile(chinaAll []Url_Ip_Location, path string) bool {
	var tempString string
	var tempStringCsv string
	var writetotxt, writetocsv bool
	var tempPath string = path

	if strings.HasSuffix(path, ".txt") {
		writetotxt = true
	}
	if strings.HasSuffix(path, ".csv") {
		writetocsv = true
	}
	if !writetotxt && !writetocsv {
		// tempPath = path

		writetotxt = true
		writetocsv = true
	}
	for _, value := range chinaAll {
		// fmt.Println(value)
		// 写入文件

		tempString += value.Url + "   " + value.Ip + "   " + value.Location.CountryCode + "   " + value.Location.RegionName + "\n"
		tempStringCsv += value.Url + "," + value.Ip + "," + value.Location.CountryCode + "," + value.Location.RegionName + "\n"
	}
	if writetotxt {
		if !strings.HasSuffix(path, ".txt") {
			path = tempPath + ".txt"
		}
		// 写入txt
		err := ioutil.WriteFile(path, []byte(tempString), 0666)
		if err != nil {
			fmt.Println(err)
		}
	}
	if writetocsv {
		if !strings.HasSuffix(path, ".csv") {
			path = tempPath + ".csv"
		}
		// 写入CSV文件

		err1 := ioutil.WriteFile(path, []byte(tempStringCsv), 0666)
		if err1 != nil {
			fmt.Println(err1)
		}
	}
	return true
}
