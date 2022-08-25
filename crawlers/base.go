package crawlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	util "multilogin_scraping/pkg/utils"
	"net/http"
	"time"

	"github.com/spf13/viper"

	"github.com/tebeka/selenium"
)

type Profile struct {
	Name   string
	Status string `json:"status"`
	Value  string `json:"value"`
	UUID   string `json:"uuid"`
}

type BaseSelenium struct {
	WebDriver selenium.WebDriver
	Profile   *Profile
}

func NewBaseSelenium() *BaseSelenium {
	return &BaseSelenium{}
}

var mla_url string = "/api/v1/profile/start?automation=true&profileId="
var profileURL string = "/api/v2/profile/"

func (ps *Profile) CreateProfile() {
	oses := []string{"win", "mac", "android", "lin"}
	//browsers := []string{"stealthfox", "mimic"}
	browsers := []string{"stealthfox"}
	values := map[string]string{
		"name":    fmt.Sprint(ps.Name, "-Crawler-", util.RandInt()),
		"os":      util.RandSliceStr(oses),
		"browser": util.RandSliceStr(browsers),
	}
	jsonData, err := json.Marshal(values)

	if err != nil {
		log.Fatalln(err)
	}
	time.Sleep(time.Second * 5)
	url := fmt.Sprint(viper.GetString("crawler.multilogin_url"), profileURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))

	defer resp.Body.Close()

	if err != nil {
		log.Fatalln(err)
	}

	// Decode data
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		log.Fatalln(err)
	}
}

// FetchProfile to get URL for remoting
func (ps *Profile) FetchProfile() *Profile {
	time.Sleep(time.Second * 5)
	url := fmt.Sprint(viper.GetString("crawler.multilogin_url"), mla_url, ps.UUID)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("ZillowCrawler Error: ", err)
		return nil
	}
	defer resp.Body.Close()

	// Decode data

	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		fmt.Println("ZillowCrawler Error: ", err)
		return nil
	}
	return ps
}

func (ps *Profile) DeleteProfile() {
	url := fmt.Sprint(viper.GetString("crawler.multilogin_url"), profileURL, ps.UUID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	if _, err := http.DefaultClient.Do(req); err != nil {
		log.Fatalln(err)
	}
}

func (bs *BaseSelenium) StartSelenium(profileName string) *BaseSelenium {
	ps := &Profile{Name: profileName}
	ps.CreateProfile()
	if ps.UUID == "" {
		fmt.Println("ZillowCrawler Error:", ps.Value)
		return nil
	}
	ps.FetchProfile()
	selenium.SetDebug(viper.GetBool("crawler.debug"))
	caps := selenium.Capabilities{}

	// Connect to Selenium
	wd, err := selenium.NewRemote(caps, ps.Value)
	if err != nil {
		fmt.Println("ZillowCrawler Error: ", err)
		return nil
	}

	bs.WebDriver = wd
	bs.Profile = ps
	return bs
}
func (bs *BaseSelenium) StopSelenium() {
	bs.WebDriver.Quit()
	time.Sleep(3 * time.Second)
	bs.Profile.DeleteProfile()
}
