package main

// IPStruct IP地址查询返回结构体
type IPStruct struct {
	IP       string   `json:"ip"`
	Country  string   `json:"country"`
	Province string   `json:"province"`
	City     string   `json:"city"`
	Location []string `json:"location"`
}

// GDStruct 高德IP查询
type GDStruct struct {
	//Status   string `json:"status"`
	//Info     string `json:"info"`
	//Infocode string `json:"infocode"`
	//Country  string `json:"country"`
	Province string `json:"province"`
	City     string `json:"city"`
	District string `json:"district"`
	Isp      string `json:"isp"`
	//Location string `json:"location"`
	Ip       string `json:"ip"`
}