package crawlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	util2 "multilogin_scraping/pkg/utils"
	"net/http"
	"time"

	"github.com/spf13/viper"

	"github.com/tebeka/selenium"
)

type Profile struct {
	Name        string
	Status      string      `json:"status"`
	Value       string      `json:"value"`
	UUID        string      `json:"uuid"`
	BrowserName string      `json:"browser_name"`
	Proxy       util2.Proxy `json:"proxy"`
	ProxyStatus bool        `json:"proxy_status"`
}

type BaseSelenium struct {
	WebDriver selenium.WebDriver
	Profile   *Profile
	Logger    *zap.Logger
}

type BaseProfileInit struct {
	Name    string `json:"name"`
	Browser string `json:"browser"`
	OS      string `json:"os"`
	Network struct {
		Proxy util2.Proxy `json:"proxy"`
	} `json:"network"`
}

func NewBaseSelenium(logger *zap.Logger) *BaseSelenium {
	return &BaseSelenium{Logger: logger}
}

var mlaUrl string = "/api/v1/profile/start?automation=true&profileId="
var profileURL string = "/api/v2/profile/"

func (ps *Profile) CreateProfile() error {
	// Create random profile multilogin app
	oses := []string{"win", "mac", "lin"}
	//browsers := []string{"stealthfox", "mimic"}
	browsers := []string{"stealthfox"}
	//values := &map[string]string{
	//	"name":    fmt.Sprint(ps.Name, "-Crawler-", util2.RandInt()),
	//	"os":      util2.RandSliceStr(oses),
	//	"browser": util2.RandSliceStr(browsers),
	//}
	ps.BrowserName = util2.RandSliceStr(browsers)

	// Apply to profile
	baseProfile := &BaseProfileInit{
		Name:    fmt.Sprint(ps.Name, "-Crawler-", util2.RandInt()),
		OS:      util2.RandSliceStr(oses),
		Browser: ps.BrowserName,
	}

	if ps.ProxyStatus == true {
		baseProfile.Network.Proxy = ps.Proxy
	}

	jsonData, err := json.Marshal(baseProfile)

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
	url := fmt.Sprint(viper.GetString("crawler.multilogin_url"), mlaUrl, ps.UUID)
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

func (bs *BaseSelenium) StartSelenium(profileName string, proxy util2.Proxy, proxyStatus bool) error {
	ps := &Profile{Name: profileName, Proxy: proxy, ProxyStatus: proxyStatus}
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
	selenium.SetDebug(viper.GetBool("crawler.selenium_debug"))
	caps := selenium.Capabilities{}

	//caps.AddFirefox(firefox.Capabilities{Args: []string{"--headless", "--no-sandbox"}})
	//caps.AddChrome(chrome.Capabilities{Prefs: map[string]interface{}{"profile.managed_default_content_settings.images": 2}})
	//caps.AddFirefox(firefox.Capabilities{Prefs: map[string]interface{}{"permissions.default.image": 2}})
	// Connect to Selenium
	wd, err := selenium.NewRemote(caps, ps.Value)
	if err != nil {
		return err
	}

	bs.WebDriver = wd
	bs.Profile = ps
	return nil
}

// FireFoxDisableImageLoading for FireFox browser
func (bs *BaseSelenium) FireFoxDisableImageLoading() error {
	err := bs.WebDriver.Get("about:config")
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 1)

	// Click Warning Button
	warningButton, err := bs.WebDriver.FindElement(selenium.ByCSSSelector, "#warningButton")
	if err != nil {
		return err
	}
	if err := warningButton.Click(); err != nil {
		return err
	}
	time.Sleep(time.Second * 1)

	// Find the Config Search Box and fill setting
	configSearchEl, err := bs.WebDriver.FindElement(selenium.ByCSSSelector, "#about-config-search")
	if err != nil {
		return err
	}
	err = configSearchEl.SendKeys("permissions.default.image")
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 1)

	// Find Edit setting, edit and fill it
	editButton, err := bs.WebDriver.FindElement(selenium.ByXPATH, "//td[@class='cell-edit']/button[@class='button-edit semi-transparent']")
	if err != nil {
		return err
	}
	err = editButton.Click()
	if err != nil {
		return err
	}
	time.Sleep(time.Second * 1)
	formEdit, err := bs.WebDriver.FindElement(selenium.ByXPATH, "//form[@id='form-edit']/input")
	if err != nil {
		return err
	}

	if err := formEdit.Clear(); err != nil {
		return err
	}

	if err := formEdit.SendKeys("2"); err != nil {
		return err
	}
	saveButton, err := bs.WebDriver.FindElement(selenium.ByXPATH, "//td[@class='cell-edit']/button[@class='primary button-save semi-transparent']")
	if err != nil {
		return err
	}

	if err := saveButton.Click(); err != nil {
		return err
	}
	return nil
}
func (bs *BaseSelenium) StopSessionBrowser(browserQuit bool) error {
	bs.Logger.Info(fmt.Sprint("Stop crawler & delete profile on Multilogin App: ", bs.Profile.UUID))
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

func (bs *BaseSelenium) GetHttpCookies() ([]*http.Cookie, error) {
	seleniumCookies, err := bs.WebDriver.GetCookies()
	if err != nil {
		return nil, err
	}
	return util2.ConvertSeleniumToHttpCookies(seleniumCookies), nil
}
