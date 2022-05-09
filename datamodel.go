package main


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