package zillow

import (
	"encoding/json"
	"fmt"
	"log"
	"multilogin_scraping/crawlers"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
	"github.com/tebeka/selenium"
)

type ZillowCrawler struct {
	WebDriver     selenium.WebDriver
	BaseSel       *crawlers.BaseSelenium
	Profile       *crawlers.Profile
	CZillow       *colly.Collector
	SearchPageReq *SearchPageReq
}

const searchURL = "https://www.zillow.com/search/GetSearchPageState.htm?searchQueryState=%s&wants={\"cat1\":[\"listResults\",\"mapResults\"],\"cat2\":[\"total\"],\"regionResults\":[\"total\"]}&requestId=5"

func NewZillowCrawler(c *colly.Collector) *ZillowCrawler {
	BaseSel := crawlers.NewBaseSelenium()
	BaseSel.StartSelenium("zillow")
	userAgent, err := BaseSel.WebDriver.ExecuteScript("return navigator.userAgent", nil)
	if err != nil {
		log.Fatalln(err)
	}
	c.UserAgent = userAgent.(string)
	return &ZillowCrawler{WebDriver: BaseSel.WebDriver, BaseSel: BaseSel, Profile: BaseSel.Profile, CZillow: c}
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

	zc.CheckVerifyHuman(pageSource)
	zc.CrawlData(zc.CrawlMapBounds(pageSource))
	//fmt.Println(zc.WebDriver.GetCookies())
}

func (zc *ZillowCrawler) CheckVerifyHuman(pageSource string) {
	if strings.Contains(pageSource, "Please verify you're a human to continue") {
		zc.BaseSel.StopSelenium()
		log.Fatalln("the website blocked Zillow Crawler")
	}
}

func (zc *ZillowCrawler) CrawlData(mapBounds MapBounds) {
	temSearch := `
		{
		"isMapVisible": true,
		"filterState": {
			"sortSelection": {
				"value": "days"
			},
			"isAllHomes": {
				"value": true
			}
		},
		"isListVisible": true,
		"mapZoom": 13,
		"pagination": {
			"currentPage": 1
		}
	}
	`
	if err := json.Unmarshal([]byte(temSearch), &zc.SearchPageReq); err != nil {
		log.Fatalln(err)
	}
	zc.SearchPageReq.MapBounds = mapBounds

	searchPageJson, err := json.Marshal(zc.SearchPageReq)
	if err != nil {
		log.Fatalln(err)
	}

	zc.CZillow.OnResponse(func(r *colly.Response) {
		data := &SearchPageRes{}
		if err := json.Unmarshal(r.Body, data); err != nil {
			log.Fatalln(err)
		}
		if len(data.Cat1.SearchResults.ListResults) > 0 {
			for _, result := range data.Cat1.SearchResults.ListResults {
				zc.ExtractData(result)
			}
		}
	})

	zc.CZillow.OnError(func(r *colly.Response, err error) {
		log.Fatalln("HTTP Status code:", r.StatusCode, "|URL:", r.Request.URL, "|Errors:", err)
	})
	zc.CZillow.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Content-Type", "application/json")
	})
	urlRun := fmt.Sprintf(searchURL, string(searchPageJson))
	zc.CZillow.Visit(urlRun)
}

func (zc *ZillowCrawler) CrawlMapBounds(pageSource string) MapBounds {
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

	return MapBounds{
		West:  mapBoundsMap["west"].(float64),
		East:  mapBoundsMap["east"].(float64),
		South: mapBoundsMap["south"].(float64),
		North: mapBoundsMap["north"].(float64),
	}
}

func (zc *ZillowCrawler) ExtractData(result SearchPageResResult) {
	propertyStatus := false
	if result.Beds > 0 || result.Baths > 0 {
		propertyStatus = true
	}
	halfBathRooms := result.HdpData.HomeInfo.Bedrooms / 2
	fullBathRooms := result.HdpData.HomeInfo.Bathrooms - halfBathRooms
	zillowData := &ZillowData{
		URL:            result.DetailURL,
		Address:        result.Address,
		PropertyStatus: propertyStatus,
		Bed:            result.Beds,
		Bath:           result.Baths,
		FullBathrooms:  fullBathRooms,
		HalfBathrooms:  halfBathRooms,
		SalesPrice:     result.HdpData.HomeInfo.Price,
		RentZestimate:  result.HdpData.HomeInfo.RentZestimate,
		Zestimate:      result.HdpData.HomeInfo.Zestimate,
	}
	if err := zc.WebDriver.Get(zillowData.URL); err != nil {
		log.Fatalln(err)
	}
	// NOTE: time to load source. Need to increase if data was not showing
	time.Sleep(10 * time.Second)
	pageSource, err := zc.WebDriver.PageSource()
	if err != nil {
		log.Fatalln(err)
	}
	zc.CheckVerifyHuman(pageSource)
	zc.ParseData(pageSource, zillowData)

	time.Sleep(3 * time.Second)
}

func (zc *ZillowCrawler) ParseData(source string, zillowData *ZillowData) {
	//htmlquery.DisableSelectorCache = true
	doc, err := htmlquery.Parse(strings.NewReader(source))

	if err != nil {
		log.Fatalln(err)
	}

	// Address
	if zillowData.Address == "" {
		addresses := htmlquery.Find(doc, "//h1/text()")
		for _, v := range addresses {
			zillowData.Address += v.Data
		}
	}

	// SF
	sf := htmlquery.FindOne(doc, "//span[@data-testid='bed-bath-item']/span[normalize-space(text())='sqft']/preceding-sibling::strong/text()")
	if sf != nil {
		if zillowData.SF, err = strconv.ParseFloat(strings.Replace(sf.Data, ",", "", -1), 64); err != nil {
			log.Fatalln(err)
		}
	}

	// Est. Payment
	estPayment := htmlquery.FindOne(doc, "//div[@class='summary-container']//span[contains(text(), 'Est. payment')]/following-sibling::span/text()")
	if estPayment != nil {
		zillowData.EstPayment = strings.TrimSpace(estPayment.Data)
	}

	// Principal & Interest $
	principalInterest := htmlquery.FindOne(doc, "//h5[normalize-space(text())='Principal & interest']/following-sibling::span/text()")
	if principalInterest != nil {
		zillowData.PrincipalInterest = strings.TrimSpace(principalInterest.Data)
	}
	// Mortgage Insurance $
	mortgageInsurance := htmlquery.FindOne(doc, "//h5[normalize-space(text())='Mortgage insurance']/following-sibling::span/text()")
	if mortgageInsurance != nil {
		zillowData.MortgageInsurance = strings.TrimSpace(mortgageInsurance.Data)
	}

	// Property Taxes $
	propertyTaxes := htmlquery.FindOne(doc, "//h5[normalize-space(text())='Property taxes']/following-sibling::span/text()")
	if propertyTaxes != nil {
		zillowData.PropertyTaxes = strings.TrimSpace(propertyTaxes.Data)
	}

	// Home Insurance $
	homeInsurance := htmlquery.FindOne(doc, "//h5[contains(text(), 'Home insurance')]/following-sibling::span/text()")
	if homeInsurance != nil {
		zillowData.HomeInsurance = strings.TrimSpace(homeInsurance.Data)
	}

	// HOA Fees $
	hoaFees := htmlquery.FindOne(doc, "//h5[contains(text(), 'HOA fee')]/following-sibling::span/text()")
	if hoaFees != nil {
		zillowData.HOAFee = strings.TrimSpace(hoaFees.Data)
	}

	// Utilities $
	utilities := htmlquery.FindOne(doc, "//h5[contains(text(), \"Utilities\")]/following-sibling::span/text()")
	if utilities != nil {
		zillowData.Utilities = strings.TrimSpace(utilities.Data)
	}

	// Estimated Sales Range
	estimatedSalesRange := htmlquery.FindOne(doc, "//span[contains(text(), 'Estimated sales range')]/span/text()")
	if estimatedSalesRange != nil {
		estimatedSalesRangeList := strings.Split(strings.TrimSpace(estimatedSalesRange.Data), "-")
		zillowData.EstimatedSalesRangeMinimum = strings.TrimSpace(estimatedSalesRangeList[0])
		zillowData.EstimatedSalesRangeMax = strings.TrimSpace(estimatedSalesRangeList[1])
	}
	fmt.Println(zillowData)
}
