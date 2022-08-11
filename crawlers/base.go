package crawlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"multilogin_scraping/utils"
	"net/http"
	"time"

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

const mla_url string = "http://127.0.0.1:35000/api/v1/profile/start?automation=true&profileId="
const createProfileURL string = "http://127.0.0.1:35000/api/v2/profile"
const deleteProfileURL string = "http://127.0.0.1:35000/api/v2/profile"

func (ps *Profile) CreateProfile() {
	oses := []string{"win", "mac", "android", "lin"}
	browsers := []string{"stealthfox", "mimic"}
	values := map[string]string{
		"name":    fmt.Sprint(ps.Name, "-Crawler-", utils.RandInt()),
		"os":      utils.RandSliceStr(oses),
		"browser": utils.RandSliceStr(browsers),
	}
	jsonData, err := json.Marshal(values)

	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(createProfileURL, "application/json", bytes.NewBuffer(jsonData))

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
func (ps *Profile) FetchProfile() {
	url := fmt.Sprint(mla_url, ps.UUID)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// Decode data
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		log.Fatalln(err)
	}
}

func (ps *Profile) DeleteProfile() {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprint(deleteProfileURL, "/", ps.UUID), nil)
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
	ps.FetchProfile()
	selenium.SetDebug(true)
	caps := selenium.Capabilities{}

	// Connect to Selenium
	wd, err := selenium.NewRemote(caps, ps.Value)
	if err != nil {
		log.Fatalln(err)
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
