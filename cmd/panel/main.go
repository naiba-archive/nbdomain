package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	whois "github.com/likexian/whois-go"
	parser "github.com/likexian/whois-parser-go"
	"github.com/naiba/nbdomain"
	"github.com/naiba/nbdomain/controller"
	"github.com/naiba/nbdomain/model"
)

var licenseIP string
var licenseDomain string

var licenseUntil time.Time
var loc *time.Location

func init() {
	var err error
	licenseUntil, err = time.Parse("2006-01-02", "2099-12-31")
	if err != nil {
		panic(err)
	}
	loc, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Println("请安装 TimeZone 服务")
		os.Exit(0)
	}
	nbdomain.DB.AutoMigrate(
		model.User{},
		model.Panel{},
		model.Cat{},
		model.Domain{},
		model.Offer{},
	)
}

func main() {
	if err := checkLicense(); err == nil {
		controller.Web()
		go updateWhois()
		go license()
		select {}
	} else {
		log.Println("检查授权失败，请联系奶爸：", err)
	}
}

type worldTime struct {
	ClientIP string    `json:"client_ip,omitempty"`
	Datetime time.Time `json:"datetime,omitempty"`
}

func license() {
	var errTime int
	for {
		if errTime > 0 {
			log.Println("授权验证失败，正在重试：", errTime)
			time.Sleep(time.Hour * 2)
		}
		if errTime > rand.Intn(10)+10 {
			log.Println("授权验证失败，请联系奶爸")
			os.Exit(0)
		}
		if checkLicense() != nil {
			errTime++
		}
		time.Sleep(time.Hour)
	}
}

func checkLicense() error {
	resp, err := http.Get("https://worldtimeapi.org/api/timezone/Asia/Shanghai")
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var tm worldTime
	err = json.Unmarshal(body, &tm)
	if err != nil {
		return err
	}
	if tm.ClientIP != licenseIP || nbdomain.CF.Web.Domain != licenseDomain ||
		time.Now().In(loc).After(licenseUntil) {
		log.Println("本产品未经授权，或授权已失效，请联系奶爸", nbdomain.CF.Web.Domain)
		os.Exit(0)
	}
	return nil
}

func updateWhois() {
	var domains []model.Domain
	for {
		nbdomain.DB.Where("whois_update is NULL OR whois_update = 0 OR DATEDIFF(now(),domains.whois_update) > 7").Find(&domains)
		for _, domain := range domains {
			result, err := whois.Whois(domain.Domain)
			var create, expire time.Time
			var register string
			if err == nil {
				var parsed parser.WhoisInfo
				parsed, err = parser.Parse(result)
				if err == nil {
					create, _ = parseTime(parsed.Domain.CreatedDate)
					expire, _ = parseTime(parsed.Domain.ExpirationDate)
					register = parsed.Registrar.Name
				}
			}
			now := time.Now()
			nbdomain.DB.Model(&domain).UpdateColumns(model.Domain{
				Registrar:   register,
				Create:      &create,
				Expire:      &expire,
				WhoisUpdate: &now,
			})
			time.Sleep(time.Minute)
		}
		time.Sleep(time.Hour)
	}
}

var timeLayouts = []string{
	"2006-01-02T15:04:05-0700",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05Z",
	"2006-01-02",
	"02.01.2006 15:04:05",
	time.RFC1123,     //= "Mon, 02 Jan 2006 15:04:05 MST"
	time.RFC1123Z,    //= "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	time.RFC850,      //= "Monday, 02-Jan-06 15:04:05 MST"
	time.RFC3339,     //= "2006-01-02T15:04:05Z07:00"
	time.RFC3339Nano, //= "2006-01-02T15:04:05.999999999Z07:00"
}

func parseTime(t string) (tt time.Time, e error) {
	for _, layout := range timeLayouts {
		tt, e = time.Parse(layout, t)
		if e == nil {
			return
		}
	}
	e = errors.New("解析失败")
	return
}
