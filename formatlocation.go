package main

func To_Url_Ip_Location(urls, ips []string, locations []Ip_location) []Url_Ip_Location {

	temp := []Url_Ip_Location{}
	for i := 0; i < len(urls); i++ {
		temp = append(temp, Url_Ip_Location{ips[i], urls[i], locations[i]})
	}
	return temp
}
