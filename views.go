package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// 定义GeoLite2数据文件mmdb路径
const dbFile = "mmdblookup --file /usr/local/geoip/GeoLite2-City_20210908/GeoLite2-City.mmdb --ip"

func indexHandler(c *gin.Context) {
	year, month, day := time.Now().Date()
	nowTime := fmt.Sprintf("%d-%s-%d", year, month, day)
	c.JSON(http.StatusOK, gin.H{
		"title":  "Welcome ipQuery",
		"author": "wuhaomiao",
		"data":   nowTime,
	})
}

// GeoLite2数据
func ipGeoLiteHandler(c *gin.Context) {
	// 获取get请求参数
	ip := c.Param("ip")

	// 如果是中国ip，调用高德api接口，返回数据
	countryName := getCountry(ip)[1]
	if countryName != "China" {
		// 如果非中国ip地址，调用GeoLite2数据
		result := selectIP(ip)
		fmt.Println(result)
		c.String(http.StatusOK, result)
		return
	}

	// 中国IP,调用高德api接口
	url := fmt.Sprintf("https://restapi.amap.com/v5/ip?key=<your key>&ip=%s&type=4", ip)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("调用高德api失败：", err)
		return
	}
	defer resp.Body.Close()
	htmlBody, _ := ioutil.ReadAll(resp.Body)
	var gdIP = GDStruct{}
	err = json.Unmarshal(htmlBody, &gdIP)
	if err != nil {
		fmt.Println("解析为json失败：", err)
		return
	}
	c.JSON(http.StatusOK, gdIP)
}

func selectIP(ip string) string {
	// 查询国家
	countryName := getCountry(ip)
	// 查询省市
	provinceName := getProvince(ip)
	// 查询城市
	cityName := getCity(ip)
	// 查询经纬度
	location := getLocation(ip)

	ipCN := IPStruct{
		IP:       ip,
		Country:  countryName[0],
		Province: provinceName[0],
		City:     cityName[0],
		Location: location,
	}
	ipInfo, err := json.Marshal(ipCN)
	if err != nil {
		fmt.Println("转换JSON格式失败", err)
	}
	return string(ipInfo)
}

// 查询国家,返回中文与英文
func getCountry(ip string) []string {
	// 定义返回的中英文名称
	cnName := ""
	enName := ""
	// 获取中文国家
	getCNShell := fmt.Sprintf("%s %s country names zh-CN", dbFile, ip)
	cmd := exec.Command("/bin/sh", "-c", getCNShell)
	//fmt.Println(getCNShell)
	cnNameByte, err := cmd.Output()
	if err == nil {
		cnName = string(cnNameByte) // 如果命令执行成功
	}

	// 获取英文国家
	getENShell := fmt.Sprintf("%s %s country names en", dbFile, ip)
	cmd = exec.Command("/bin/sh", "-c", getENShell)
	//fmt.Println(getENShell)
	enNameByte, err := cmd.Output()
	if err == nil {
		enName = string(enNameByte) // 如果命令执行成功
	}

	// 调整输出的字符串格式
	pat := `".*"` //正则
	//pat := `(?<=").*(?=")`  // 正则，去掉了两边的引号
	ret, _ := regexp.Compile(pat)
	cnName = strings.Trim(ret.FindString(cnName), "\"") // 匹配引号里面的字符串，去掉引号
	enName = strings.Trim(ret.FindString(enName), "\"")

	// 返回名字，第1个为中文，第2个为英文
	countryNames := []string{
		cnName,
		enName,
	}
	return countryNames
}

// 查询省市,返回中文与英文
func getProvince(ip string) []string {
	// 定义返回的中英文名称
	cnName := ""
	enName := ""
	// 获取中文省市，省为列表格式，0取列表第1个字典数据
	getCNShell := fmt.Sprintf("%s %s subdivisions 0 names zh-CN", dbFile, ip)
	cmd := exec.Command("/bin/sh", "-c", getCNShell)
	//fmt.Println(getCNShell)
	cnNameByte, err := cmd.Output()
	if err == nil {
		cnName = string(cnNameByte) // 如果命令执行成功
	}

	// 获取英文省市
	getENShell := fmt.Sprintf("%s %s subdivisions 0 names en", dbFile, ip)
	cmd = exec.Command("/bin/sh", "-c", getENShell)
	//fmt.Println(getENShell)
	enNameByte, err := cmd.Output()
	if err == nil {
		enName = string(enNameByte) // 如果命令执行成功
	}

	// 调整输出的字符串格式
	pat := `".*"` //正则
	//pat := `(?<=").*(?=")`
	ret, _ := regexp.Compile(pat)
	cnName = strings.Trim(ret.FindString(cnName), "\"")
	enName = strings.Trim(ret.FindString(enName), "\"")

	// 如果获取不到中文名，返回英文名
	if cnName == "" {
		cnName = enName
	}

	// 返回名字，第1个为中文，第2个为英文
	countryNames := []string{
		cnName,
		enName,
	}
	return countryNames
}

// 查询城市,返回中文与英文
func getCity(ip string) []string {
	// 定义返回的中英文名称
	cnName := ""
	enName := ""
	// 获取中文城市
	getCNShell := fmt.Sprintf("%s %s city names zh-CN", dbFile, ip)
	cmd := exec.Command("/bin/sh", "-c", getCNShell)
	//fmt.Println(getCNShell)
	cnNameByte, err := cmd.Output()
	if err == nil {
		cnName = string(cnNameByte) // 如果命令执行成功
	}

	// 获取英文城市
	getENShell := fmt.Sprintf("%s %s city names en", dbFile, ip)
	cmd = exec.Command("/bin/sh", "-c", getENShell)
	//fmt.Println(getENShell)
	enNameByte, err := cmd.Output()
	if err == nil {
		enName = string(enNameByte) // 如果命令执行成功
	}

	// 调整输出的字符串格式
	pat := `".*"` //正则
	//pat := `(?<=").*(?=")`
	ret, _ := regexp.Compile(pat)
	cnName = strings.Trim(ret.FindString(cnName), "\"")
	enName = strings.Trim(ret.FindString(enName), "\"")

	// 返回名字，第1个为中文，第2个为英文
	countryNames := []string{
		cnName,
		enName,
	}
	return countryNames
}

// 查询纬度,返回中文与英文
func getLocation(ip string) []string {
	// 定义返回的中英文名称
	cnName := ""
	enName := ""
	// 获取经度
	getCNShell := fmt.Sprintf("%s %s location longitude", dbFile, ip)
	cmd := exec.Command("/bin/sh", "-c", getCNShell)
	//fmt.Println(getCNShell)
	cnNameByte, err := cmd.Output()
	if err == nil {
		cnName = string(cnNameByte) // 如果命令执行成功
	}

	// 获取纬度
	getENShell := fmt.Sprintf("%s %s location latitude", dbFile, ip)
	cmd = exec.Command("/bin/sh", "-c", getENShell)
	//fmt.Println(getENShell)
	enNameByte, err := cmd.Output()
	if err == nil {
		enName = string(enNameByte) // 如果命令执行成功
	}

	// 调整输出的字符串格式
	pat := `[0-9]*\.[0-9]*` //正则
	ret, _ := regexp.Compile(pat)
	cnName = strings.Trim(ret.FindString(cnName), "\"")
	enName = strings.Trim(ret.FindString(enName), "\"")

	// 返回名字，第1个为中文，第2个为英文
	countryNames := []string{
		cnName,
		enName,
	}
	return countryNames
}
