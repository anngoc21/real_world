package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
	"log"
	"net/http"
	"os"
	"sort"
)

type IPInfo struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
	Readme   string `json:"readme"`
	Count    int64
	Time     int64
}

func (receiver IPInfo) string() string {
	return fmt.Sprintf("%s:%s\n-%s-%s-%d", "ip", receiver.IP, receiver.Loc, receiver.Timezone, receiver.Count)
}

var IpInfos = []*IPInfo{}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		if cookie, _ := c.Request.Cookie("ipinfo"); cookie != nil {
			sDec, _ := b64.StdEncoding.DecodeString(cookie.Value)
			fmt.Println(string(sDec))
			ip := IPInfo{}
			err := json.Unmarshal([]byte(sDec), &ip)
			if err == nil {
				addIPInfo(ip)
			}
		}
		c.HTML(http.StatusOK, "index.tmpl.html", gin.H{
			"top":  sortTop100(),
			"last": getLast(),
		})
	})

	router.Run(":" + port)
}
func getLast() []*IPInfo {
	sort.Slice(IpInfos, func(i, j int) bool {
		if IpInfos[i].Time > IpInfos[j].Time {
			return true
		}
		return false
	})
	if len(IpInfos) > 100 {
		return IpInfos[:100]
	}
	return IpInfos
}
func sortTop100() []*IPInfo {
	sort.Slice(IpInfos, func(i, j int) bool {
		if IpInfos[i].Count > IpInfos[j].Count {
			return true
		}
		return false
	})
	if len(IpInfos) > 100 {
		return IpInfos[:100]
	}
	return IpInfos
}
func addIPInfo(ip IPInfo) {
	flag := false
	for _, p := range IpInfos {
		if p.IP == ip.IP {
			p.Count += 1
			flag = true
			break
		}
	}
	if !flag {
		ip.Count = 2
		IpInfos = append(IpInfos, &ip)
	}
}
