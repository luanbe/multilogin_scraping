package crawlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tebeka/selenium/firefox"
	"go.uber.org/zap"
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
	logger    *zap.Logger
}

func NewBaseSelenium(logger *zap.Logger) *BaseSelenium {
	return &BaseSelenium{logger: logger}
}

var mla_url string = "/api/v1/profile/start?automation=true&profileId="
var profileURL string = "/api/v2/profile/"

func (ps *Profile) CreateProfile() error {
	oses := []string{"win", "mac", "android", "lin"}
	//browsers := []string{"stealthfox", "mimic"}
	browsers := []string{"stealthfox"}
	values := &map[string]string{
		"name":    fmt.Sprint(ps.Name, "-Crawler-", util.RandInt()),
		"os":      util.RandSliceStr(oses),
		"browser": util.RandSliceStr(browsers),
	}
	jsonData, err := json.Marshal(values)

	if err != nil {
		return err
	}
	url := fmt.Sprint(viper.GetString("crawler.multilogin_url"), profileURL)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))

	defer resp.Body.Close()

	if err != nil {
		return err
	}

	// Decode data
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return err
	}
	return nil
}

// FetchProfile to get URL for remoting
func (ps *Profile) FetchProfile() error {
	url := fmt.Sprint(viper.GetString("crawler.multilogin_url"), mla_url, ps.UUID)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode data
	if err := json.NewDecoder(resp.Body).Decode(ps); err != nil {
		return err
	}
	return nil
}

func (ps *Profile) DeleteProfile() error {
	url := fmt.Sprint(viper.GetString("crawler.multilogin_url"), profileURL, ps.UUID)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	if _, err := http.DefaultClient.Do(req); err != nil {
		return err
	}
	return nil
}

func (bs *BaseSelenium) StartSelenium(profileName string) error {
	ps := &Profile{Name: profileName}
	if err := ps.CreateProfile(); err != nil {
		return err
	}
	if ps.UUID == "" {
		return fmt.Errorf(ps.Value)
	}
	if err := ps.FetchProfile(); err != nil {
		return err
	}
	time.Sleep(3 * time.Second)
	selenium.SetDebug(viper.GetBool("crawler.debug"))
	caps := selenium.Capabilities{}
	caps.AddFirefox(firefox.Capabilities{Args: []string{"--headless"}})
	// Connect to Selenium
	wd, err := selenium.NewRemote(caps, ps.Value)
	if err != nil {
		return err
	}

	bs.WebDriver = wd
	bs.Profile = ps
	return nil
}
func (bs *BaseSelenium) StopSessionBrowser(browserQuit bool) error {
	bs.logger.Info(fmt.Sprint("Stop browser & delete profile on Multilogin App: ", bs.Profile.UUID))
	if browserQuit == true {
		if err := bs.WebDriver.Quit(); err != nil {
			return err
		}
	}
	time.Sleep(5 * time.Second)
	if err := bs.Profile.DeleteProfile(); err != nil {
		return err
	}
	return nil
}
