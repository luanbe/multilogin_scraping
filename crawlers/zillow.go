package crawlers

import (
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"log"
	"strings"

	"github.com/tebeka/selenium"
)

type ZillowCrawler struct {
	WebDriver selenium.WebDriver
	BaseSel   *BaseSelenium
	Profile   *Profile
}

func NewZillowCrawler() *ZillowCrawler {
	BaseSel := NewBaseSelenium()
	BaseSel.StartSelenium("zillow")
	return &ZillowCrawler{WebDriver: BaseSel.WebDriver, BaseSel: BaseSel, Profile: BaseSel.Profile}
}

func (zc *ZillowCrawler) RunZillowCrawler() {
	address := "Urb Santa Elvira, PR"
	address = strings.Replace(address, " ", "-", -1)
	url := fmt.Sprint("https://www.zillow.com/homes/", address, "_rb/")
	defer zc.BaseSel.StopSelenium()
	if err := zc.WebDriver.Get(url); err != nil {
		log.Fatalln(err)
	}
	pageSource, err := zc.WebDriver.PageSource()
	if err != nil {
		log.Fatalln(err)
	}

	if err := zc.CheckVerifyHuman(pageSource); err != nil {
		log.Fatalln(err)
	}
	mapBounds := zc.ReadMapBounds(pageSource)
	fmt.Println(mapBounds)
	// fmt.Println(zc.WebDriver.GetCookies())
}

func (zc *ZillowCrawler) CheckVerifyHuman(pageSource string) error {
	if strings.Contains(pageSource, "Please verify you're a human to continue") {
		return fmt.Errorf("the website blocked Zillow Crawler")
	}
	return nil
}

func (zc *ZillowCrawler) ReadMapBounds(pageSource string) map[string]float64 {
	doc, err := htmlquery.Parse(strings.NewReader(pageSource))
	if err != nil {
		log.Fatalln(err)
	}
	script := htmlquery.FindOne(doc, "//script[@data-zrr-shared-data-key='mobileSearchPageStore']")
	if script == nil {
		log.Fatalln("Zillow Crawler didn't find the Map Bounds")
	}
	dataScript := strings.Replace(script.FirstChild.Data, "<!--", "", -1)
	jsonString := strings.Replace(dataScript, "-->", "", -1)

	// Declared an empty map interface
	var result map[string]interface{}

	// Unmarshal or Decode the JSON to the interface.
	if err := json.Unmarshal([]byte(jsonString), &result); err != nil {
		log.Fatalln(err)
	}

	queryState := result["queryState"]
	if queryState == nil {
		log.Fatalln("Not found queryState in Json data")
	}
	queryStateMap := queryState.(map[string]interface{})

	mapBounds := queryStateMap["mapBounds"]
	if mapBounds == nil {
		log.Fatalln("Not found mapBounds in Json data")
	}
	mapBoundsMap := mapBounds.(map[string]interface{})

	return map[string]float64{
		"west":  mapBoundsMap["west"].(float64),
		"east":  mapBoundsMap["east"].(float64),
		"south": mapBoundsMap["south"].(float64),
		"north": mapBoundsMap["north"].(float64),
	}
}
