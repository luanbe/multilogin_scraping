package movoto

import (
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
	"github.com/google/go-querystring/query"
	"github.com/icrowley/fake"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/net/html"
	"gorm.io/gorm"
	"multilogin_scraping/app/schemas"
	"multilogin_scraping/crawlers"
	util2 "multilogin_scraping/pkg/utils"
	"strings"
	"time"
)

type MovotoCrawler struct {
	BaseSel        *crawlers.BaseSelenium
	CMovoto        *colly.Collector
	Logger         *zap.Logger
	CrawlerBlocked bool
	BrowserTurnOff bool
	CrawlerSchemas *CrawlerSchemas
	Proxy          *util2.Proxy
	Doc            *html.Node
}

type CrawlerSchemas struct {
	MovotoData *schemas.MovotoData
	SearchReq  *schemas.MovotoSearchPageReq
}

const searchURL = "https://parser-external.geo.moveaws.com/suggest?%s"

func NewMovotoCrawler(
	db *gorm.DB,
	logger *zap.Logger,
	proxy *util2.Proxy,
) *MovotoCrawler {
	BaseSel := crawlers.NewBaseSelenium(logger)
	c := colly.NewCollector()
	userAgent := fake.UserAgent()
	c.UserAgent = userAgent

	return &MovotoCrawler{
		BaseSel:        BaseSel,
		CMovoto:        c,
		Logger:         logger,
		CrawlerBlocked: false,
		BrowserTurnOff: false,
		Proxy:          proxy,
		CrawlerSchemas: &CrawlerSchemas{
			MovotoData: &schemas.MovotoData{},
			SearchReq:  &schemas.MovotoSearchPageReq{},
		},
	}
}

// NewBrowser to start new selenium
func (rc *MovotoCrawler) NewBrowser() error {
	if err := rc.BaseSel.StartSelenium("movoto", rc.Proxy, viper.GetBool("crawler.movoto_crawler.proxy_status")); err != nil {
		return err
	}
	// Disable image loading
	if viper.GetBool("crawler.disable_load_images") == true {
		if rc.BaseSel.Profile.BrowserName == "stealthfox" {
			if err := rc.BaseSel.FireFoxDisableImageLoading(); err != nil {
				return err
			}
		}
	}
	return nil
}

// UserAgentBrowserToColly for coping useragent from browser to colly
func (rc *MovotoCrawler) UserAgentBrowserToColly() error {
	userAgent, err := rc.BaseSel.WebDriver.ExecuteScript("return navigator.userAgent", nil)

	if err != nil {
		return err
	}

	if userAgent == nil {
		userAgent = fake.UserAgent()
	}

	rc.CMovoto.UserAgent = userAgent.(string)

	return nil
}

// RunMovotoCrawlerAPI with a loop to run crawler
func (rc *MovotoCrawler) RunMovotoCrawlerAPI(mprID string) error {
	rc.Logger.Info("Zillow Data is crawling...")
	rc.CrawlerSchemas.MovotoData.URL = fmt.Sprint(viper.GetString("crawler.movoto_crawler.url"), "realestateandhomes-detail/M", mprID)

	err := func(address string) error {
		if err := rc.BaseSel.WebDriver.Get(address); err != nil {
			return err
		}
		// NOTE: time to load source. Need to increase if data was not showing
		time.Sleep(viper.GetDuration("crawler.movoto_crawler.time_load_source") * time.Second)

		pageSource, err := rc.BaseSel.WebDriver.PageSource()
		if err != nil {
			rc.BrowserTurnOff = true
			return err
		}

		if err := rc.ByPassVerifyHuman(pageSource, address); err != nil {
			return err
		}

		// TODO: Add Parse Data
		if err := rc.ParseData(pageSource); err != nil {
			return err
		}

		return nil

	}(rc.CrawlerSchemas.MovotoData.URL)

	if err != nil {
		rc.Logger.Error(err.Error())
		rc.Logger.Error("Failed to crawl data")
		return err
		// TODO: Update error for crawling here
	}
	rc.Logger.Info("Completed to crawl data")
	return nil
}

// ByPassVerifyHuman to bypass verify from Movoto website
func (rc *MovotoCrawler) ByPassVerifyHuman(pageSource string, url string) error {
	if rc.IsVerifyHuman(pageSource) == true {
		rc.CrawlerBlocked = true

	}
	if rc.CrawlerBlocked == true {
		for i := 0; i < 3; i++ {
			err := rc.BaseSel.WebDriver.Get(url)
			if err != nil {
				return err
			}
			pageSource, _ = rc.BaseSel.WebDriver.PageSource()
			if rc.IsVerifyHuman(pageSource) == false {
				rc.CrawlerBlocked = false
				return nil
			}
		}
		return fmt.Errorf("Crawler blocked for checking verify hunman")
	}
	return nil
}

// IsVerifyHuman to check website is blocking
func (rc *MovotoCrawler) IsVerifyHuman(pageSource string) bool {
	if strings.Contains(pageSource, "Please verify you're a human to continue") || strings.Contains(pageSource, "Let's confirm you are human") {
		return true
	}
	return false
}

func (rc *MovotoCrawler) CrawlSearchData(search string) (string, error) {
	// NOTE: We only take browser cookies when getting block from movoto website
	//cookies, err := rc.BaseSel.GetHttpCookies()
	//if err != nil {
	//	rc.Logger.Error(err.Error())
	//	return
	//}
	data := &schemas.MovotoSearchPageRes{}
	rc.CrawlerSchemas.SearchReq = &schemas.MovotoSearchPageReq{
		Input:     search,
		ClientID:  "rdc-home",
		Limit:     10,
		AreaTypes: "address",
	}
	searchPageQuery, err := query.Values(rc.CrawlerSchemas.SearchReq)
	if err != nil {
		return "nil", err
	}

	rc.CMovoto.OnError(func(r *colly.Response, err error) {
		rc.Logger.Error(fmt.Sprint("HTTP Status code:", r.StatusCode, "|URL:", r.Request.URL, "|Errors:", err))
		return
	})

	rc.CMovoto.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Content-Type", "application/json")
	})

	rc.CMovoto.OnResponse(func(r *colly.Response) {
		if err := json.Unmarshal(r.Body, data); err != nil {
			rc.Logger.Error(err.Error())
			return
		}
	})

	urlRun := fmt.Sprintf(searchURL, searchPageQuery.Encode())

	// NOTE: We only take browser cookies when getting block from movoto website
	//err = zc.CZillow.SetCookies(urlRun, cookies)
	//if err != nil {
	//	zc.SearchDataByCollyStatus = false
	//	zc.ShowLogError(err.Error())
	//	return
	//}
	if err := rc.CMovoto.Visit(urlRun); err != nil {
		return "", err
	}
	for _, result := range data.Autocomplete {
		if result.FullAddress[0] == search {
			return result.MprID, nil
		}
	}
	return "", nil
}

func (rc *MovotoCrawler) ParseData(source string) error {
	var err error
	if rc.Doc, err = htmlquery.Parse(strings.NewReader(source)); err != nil {
		return err
	}

	//if err := rc.ClickElements(); err != nil {
	//	return err
	//}

	rc.ParseBed()
	rc.ParseBath()
	rc.ParsePropertyStatus()
	rc.ParseFullBathrooms()
	rc.ParseSF()
	rc.ParseSalePrice()
	rc.ParseEstPayment()
	rc.ParsePrincipalInterest()
	rc.ParseMortgageInsurance()
	rc.ParsePropertyTaxes()
	rc.ParseHomeInsurance()
	rc.ParseHOAFees()
	rc.ParseOverview()
	rc.ParseSource()
	rc.ParsePropertyType()
	rc.ParseYearBuilt()
	rc.ParseNaturalGas()
	rc.ParseCentralAir()
	rc.ParseGarageSpaces()
	rc.ParseLotSize()
	rc.ParseLotSizeAcres()
	rc.ParseInteriorFeatures()
	rc.ParsePrimaryBedroomLevel()
	rc.ParseFlooringType()
	rc.ParseHeatingType()
	rc.ParseTotalParkingSpaces()
	rc.ParseLotFeatures()
	rc.ParseProperySubType()
	rc.ParseConstructionMaterials()
	rc.ParseFoundation()
	rc.ParseSubdivision()
	rc.ParseElementarySchool()
	rc.ParseMiddleSchool()
	rc.ParseHighSchool()
	rc.ParseDataSource()
	rc.ParsePictures()
	return nil
}

//func (rc *MovotoCrawler) ClickElements() error {
//	//
//	//if err := rc.ClickPaymentCalculator(); err != nil {
//	//	return err
//	//}
//	//
//	//if err := rc.ClickPropertyHistory(); err != nil {
//	//	return err
//	//}
//	//
//	//pageSource, err := rc.BaseSel.WebDriver.PageSource()
//	//
//	//if rc.Doc, err = htmlquery.Parse(strings.NewReader(pageSource)); err != nil {
//	//	rc.Logger.Warn("Parse Principal Interest: Can not load page source")
//	//	return err
//	//}
//	//return nil
//}
//
//func (rc *MovotoCrawler) ClickPaymentCalculator() error {
//	//paymentCalculator, err := rc.BaseSel.WebDriver.FindElement(selenium.ByXPATH, "//section[@id=\"payment_calculator\"]")
//	//
//	//if err != nil {
//	//	rc.Logger.Warn("Click Payment Calculator: Not found section payment_calculator")
//	//	return err
//	//}
//	//
//	//if err := paymentCalculator.Click(); err != nil {
//	//	rc.Logger.Warn("Click Payment Calculator: Can not click payment_calculator")
//	//	return err
//	//}
//	//_, err = rc.BaseSel.WebDriver.ExecuteScript("document.getElementById(\"payment_calculator\").scrollIntoView();", nil)
//	//
//	//if err != nil {
//	//	rc.Logger.Warn("Click Payment Calculator: Error on scoll to payment_calculator")
//	//	return err
//	//}
//	//time.Sleep(2 * time.Second)
//	//return nil
//}
//
//func (rc *MovotoCrawler) ClickPropertyHistory() error {
//	//propertyHistory, err := rc.BaseSel.WebDriver.FindElement(selenium.ByXPATH, "//section[@id=\"section_property_history\"]")
//	//
//	//if err != nil {
//	//	rc.Logger.Warn("Click Property History: Not found section_property_history")
//	//	return err
//	//}
//	//
//	//if err := propertyHistory.Click(); err != nil {
//	//	rc.Logger.Warn("Click Property History: Can not click section_property_history")
//	//	return err
//	//}
//	//_, err = rc.BaseSel.WebDriver.ExecuteScript("document.getElementById(\"section_property_history\").scrollIntoView();", nil)
//	//
//	//if err != nil {
//	//	rc.Logger.Warn("Click Property History: Error on scoll to section_property_history")
//	//	return err
//	//}
//	//time.Sleep(2 * time.Second)
//	//return nil
//}

// ParseBed for crawling Bed data
func (rc *MovotoCrawler) ParseBed() {
	//bedDoc := htmlquery.FindOne(
	//	rc.Doc,
	//	"//li[contains(@data-testid,\"property-meta-beds\")]/span/text()",
	//)
	//if bedDoc == nil {
	//	rc.Logger.Warn("Parse Bed: Not found element.")
	//	return
	//}
	//bedText := htmlquery.InnerText(bedDoc)
	//if bedInt, err := strconv.Atoi(bedText); err != nil {
	//	rc.Logger.Error(fmt.Sprintf("Parse Bed: %v", err.Error()))
	//} else {
	//	rc.CrawlerSchemas.MovotoData.Bed = bedInt
	//}
}

// ParseBath for crawling Bath data
func (rc *MovotoCrawler) ParseBath() {
	//bathDoc := htmlquery.FindOne(
	//	rc.Doc,
	//	"//li[contains(@data-testid,\"property-meta-baths\")]/span/text()",
	//)
	//if bathDoc == nil {
	//	rc.Logger.Warn("Parse Bath: Not found element.")
	//	return
	//}
	//bathText := htmlquery.InnerText(bathDoc)
	//if bathInt, err := strconv.Atoi(bathText); err != nil {
	//	rc.Logger.Error(fmt.Sprintf("Parse Bath: %v", err.Error()))
	//} else {
	//	rc.CrawlerSchemas.MovotoData.Bath = bathInt
	//}
}

// ParsePropertyStatus for checking Property Status
func (rc *MovotoCrawler) ParsePropertyStatus() {
	//if rc.CrawlerSchemas.MovotoData.Bed > 0 || rc.CrawlerSchemas.MovotoData.Bath > 0 {
	//	rc.CrawlerSchemas.MovotoData.PropertyStatus = true
	//}
}

// ParseFullBathrooms for crawling full bathrooms
func (rc *MovotoCrawler) ParseFullBathrooms() {
	//fullBathroomsDoc := htmlquery.FindOne(rc.Doc, "//li[contains(text(), \"Full Bathrooms\")]/text()")
	//if fullBathroomsDoc == nil {
	//	rc.Logger.Warn("Parse full bathrooms: Not found element.")
	//	return
	//}
	//fullBathroomsText := htmlquery.InnerText(fullBathroomsDoc)
	//fullBathroomsSlice := strings.Split(fullBathroomsText, ":")
	//fullBathroomsText = fullBathroomsSlice[1]
	//
	//fullBathroomsFloat, err := util2.ConvertToFloat(fullBathroomsText)
	//if err != nil {
	//	rc.Logger.Error(fmt.Sprintf("Parse full bathrooms: %v", err.Error()))
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.FullBathrooms = fullBathroomsFloat
}

// ParseSF for parsing SF
func (rc *MovotoCrawler) ParseSF() {
	//sfDoc := htmlquery.FindOne(rc.Doc, "//li[contains(@data-testid,\"property-meta-sqft\")]//span[@class=\"meta-value\"]/text()")
	//if sfDoc == nil {
	//	rc.Logger.Warn("Parse SF: Not found element.")
	//	return
	//}
	//sfFloat, err := util2.ConvertToFloat(htmlquery.InnerText(sfDoc))
	//
	//if err != nil {
	//	rc.Logger.Error(fmt.Sprintf("Parse SF: %v", err.Error()))
	//	return
	//}
	//
	//rc.CrawlerSchemas.MovotoData.SF = sfFloat
}

// ParseSalePrice for parsing Sale Price
func (rc *MovotoCrawler) ParseSalePrice() {
	//salePriceDoc := htmlquery.FindOne(rc.Doc, "//div[@data-testid=\"list-price\"]//text()")
	//if salePriceDoc == nil {
	//	rc.Logger.Warn("Parse Sale Price: Not found element.")
	//	return
	//}
	//
	//salePriceFloat, err := util2.ConvertToFloat(htmlquery.InnerText(salePriceDoc))
	//
	//if err != nil {
	//	rc.Logger.Warn(fmt.Sprintf("Parse Sale Price: %v", err.Error()))
	//	return
	//}
	//
	//rc.CrawlerSchemas.MovotoData.SalesPrice = salePriceFloat
}

// ParseEstPayment for parsing Estimate Payment
func (rc *MovotoCrawler) ParseEstPayment() {
	//estPaymentDoc := htmlquery.FindOne(rc.Doc, "//*[@data-testid=\"est-payment\"]//*[contains(text(), \"$\")]")
	//if estPaymentDoc == nil {
	//	rc.Logger.Warn("Parse Est Payment: Not found element.")
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.EstPayment = htmlquery.InnerText(estPaymentDoc)
}

// ParsePrincipalInterest for parsing Principal Interest
// NOTE: This parsing will be loading page source again.
func (rc *MovotoCrawler) ParsePrincipalInterest() {
	//principalInterestDoc := htmlquery.FindOne(rc.Doc, "//div[contains(text(), \"Principal & Interest\")]/following-sibling::div")
	//
	//if principalInterestDoc == nil {
	//	rc.Logger.Warn("Parse Principal Interest: Not found element.")
	//	return
	//}
	//
	//rc.CrawlerSchemas.MovotoData.PrincipalInterest = htmlquery.InnerText(principalInterestDoc)
}

// ParseMortgageInsurance for parsing Mortgage Insurance
func (rc *MovotoCrawler) ParseMortgageInsurance() {
	//mortgageInsuranceDoc := htmlquery.FindOne(rc.Doc, "//div[contains(text(), \"Mortgage Insurance\")]/following-sibling::div")
	//if mortgageInsuranceDoc == nil {
	//	rc.Logger.Warn("Parse Mortgage Insurance: Not found element.")
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.MortgageInsurance = htmlquery.InnerText(mortgageInsuranceDoc)
}

// ParsePropertyTaxes for parsing Property Taxes
func (rc *MovotoCrawler) ParsePropertyTaxes() {
	//propertyTaxesDoc := htmlquery.FindOne(rc.Doc, "//div[contains(text(), \"Property tax\")]/following-sibling::div")
	//if propertyTaxesDoc == nil {
	//	rc.Logger.Warn("Parse Property Taxes: Not found element.")
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.PropertyTaxes = htmlquery.InnerText(propertyTaxesDoc)
}

// ParseHomeInsurance for parsing Home Insurance
func (rc *MovotoCrawler) ParseHomeInsurance() {
	//homeInsuranceDoc := htmlquery.FindOne(rc.Doc, "//div[contains(text(), \"Home Insurance\")]/following-sibling::div")
	//if homeInsuranceDoc == nil {
	//	rc.Logger.Warn("Parse Home Insurance: Not found element.")
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.HomeInsurance = htmlquery.InnerText(homeInsuranceDoc)
}

// ParseHOAFees for parsing HOA Fees
func (rc *MovotoCrawler) ParseHOAFees() {
	//hoaFeesDoc := htmlquery.FindOne(rc.Doc, "//div[contains(text(), \"HOA fees\")]/following-sibling::div")
	//if hoaFeesDoc == nil {
	//	rc.Logger.Warn("Parse HOA Fees: Not found element.")
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.HOAFee = htmlquery.InnerText(hoaFeesDoc)
}

// ParseOverview for parsing overview
func (rc *MovotoCrawler) ParseOverview() {
	//overviewDoc := htmlquery.FindOne(rc.Doc, "//div[@id=\"section_property_details\"]")
	//if overviewDoc == nil {
	//	rc.Logger.Warn("Parse Overview: Not found element.")
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.HOAFee = htmlquery.InnerText(overviewDoc)
}

// ParseSource for parsing source
func (rc *MovotoCrawler) ParseSource() {
	//sourceDoc := htmlquery.FindOne(rc.Doc, "(//div[@id=\"content-property_history\"]//table)[1]/tbody//tr[position()=1]/td[last()]")
	//if sourceDoc == nil {
	//	rc.Logger.Warn("Parse Source: Not found element.")
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.Source = htmlquery.InnerText(sourceDoc)
}

// ParsePropertyType for parsing Property Type
func (rc *MovotoCrawler) ParsePropertyType() {
	//propertyTypeDoc := htmlquery.FindOne(rc.Doc, "//div[@id=\"section_summary\"]//span[contains(text(), \"Property Type\")]/parent::div/following-sibling::div")
	//if propertyTypeDoc == nil {
	//	rc.Logger.Warn("Parse Property Type: Not found element.")
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.PropertyType = htmlquery.InnerText(propertyTypeDoc)
}

// ParseYearBuilt for parsing Year Built
func (rc *MovotoCrawler) ParseYearBuilt() {
	//yearBuiltDoc := htmlquery.FindOne(rc.Doc, "//div[@id=\"section_summary\"]//span[contains(text(), \"Year Built\")]/parent::div/following-sibling::div")
	//if yearBuiltDoc == nil {
	//	rc.Logger.Warn("Parse Year Built: Not found element.")
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.YearBuilt = htmlquery.InnerText(yearBuiltDoc)
}

// ParseNaturalGas for parsing Natural Gas
func (rc *MovotoCrawler) ParseNaturalGas() {
	//naturalGasDoc := htmlquery.FindOne(rc.Doc, "//*[contains(text(), \"Natural Gas\") or contains(text(), \"natural gas\")]")
	//if naturalGasDoc == nil {
	//	rc.Logger.Warn("Parse Natural Gas: Not found element.")
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.NaturalGas = true
}

// ParseCentralAir for parsing Central Air
func (rc *MovotoCrawler) ParseCentralAir() {
	//centerAirDoc := htmlquery.FindOne(rc.Doc, "//*[contains(text(), \"Central Air\") or contains(text(), \"central air\")]")
	//if centerAirDoc == nil {
	//	rc.Logger.Warn("Parse Central Air: Not found element.")
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.CentralAir = true
}

// ParseGarageSpaces for parsing Garage Spaces
func (rc *MovotoCrawler) ParseGarageSpaces() {
	//garageSpacesDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"Garage Spaces\")]")
	//if garageSpacesDoc == nil {
	//	rc.Logger.Warn("Parse Garage Spaces: Not found element.")
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.OfGarageSpaces = strings.TrimSpace(strings.Replace(htmlquery.InnerText(garageSpacesDoc), "Garage Spaces:", "", -1))
}

// ParseLotSize for parsing Lot Size
func (rc *MovotoCrawler) ParseLotSize() {
	//lotSizeDoc := htmlquery.FindOne(rc.Doc, "//li[contains(@data-testid,\"property-meta-lot-size\")]//span[@class=\"meta-value\"]/text()")
	//if lotSizeDoc == nil {
	//	rc.Logger.Warn("Parse Lot Size: Not found element.")
	//	return
	//}
	//lotSizeFloat, err := util2.ConvertToFloat(htmlquery.InnerText(lotSizeDoc))
	//
	//if err != nil {
	//	rc.Logger.Error(fmt.Sprintf("Parse Lot Size: %v", err.Error()))
	//	return
	//}
	//
	//rc.CrawlerSchemas.MovotoData.LotSizeSF = lotSizeFloat
}

// ParseLotSizeAcres for parsing Lot Size Acres
func (rc *MovotoCrawler) ParseLotSizeAcres() {
	//lotSizeAcresDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"Lot Size Acres\")]")
	//if lotSizeAcresDoc == nil {
	//	rc.Logger.Warn("Parse Lot Size Acres: Not found element.")
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.OfGarageSpaces = strings.TrimSpace(strings.Replace(htmlquery.InnerText(lotSizeAcresDoc), "Lot Size Acres:", "", -1))
}

// ParseInteriorFeatures for parsing Interior Features
func (rc *MovotoCrawler) ParseInteriorFeatures() {
	//var interiorFeaturesList []string
	//interiorFeaturesDocs := htmlquery.Find(rc.Doc, "//h4[contains(text(), \"Interior Features\")]/following-sibling::ul[1]/li/text()")
	//if interiorFeaturesDocs == nil {
	//	rc.Logger.Warn("Parse Interior Features: Not found element.")
	//	return
	//}
	//for _, v := range interiorFeaturesDocs {
	//	interiorFeaturesList = append(interiorFeaturesList, htmlquery.InnerText(v))
	//}
	//rc.CrawlerSchemas.MovotoData.InteriorFeatures = strings.Join(interiorFeaturesList, ", ")
}

// ParsePrimaryBedroomLevel for parsing Primary Bedroom Level
func (rc *MovotoCrawler) ParsePrimaryBedroomLevel() {
	//primaryBedroomLevelDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"Bedrooms\")]")
	//if primaryBedroomLevelDoc == nil {
	//	rc.Logger.Warn("Parse Primary Bedroom Level: Not found element.")
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.PrimaryBedroomLevel = strings.TrimSpace(strings.Replace(htmlquery.InnerText(primaryBedroomLevelDoc), "Bedrooms:", "", -1))
}

// ParseFlooringType for parsing Flooring Type
func (rc *MovotoCrawler) ParseFlooringType() {
	//flooringTypeDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"Flooring:\")]")
	//if flooringTypeDoc == nil {
	//	rc.Logger.Warn("Parse Flooring Type: Not found element.")
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.FlooringType = strings.TrimSpace(strings.Replace(htmlquery.InnerText(flooringTypeDoc), "Flooring:", "", -1))
}

// ParseHeatingType for parsing Heating Type
func (rc *MovotoCrawler) ParseHeatingType() {
	//heatingTypeDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"Heating Features:\")]")
	//if heatingTypeDoc == nil {
	//	rc.Logger.Warn("Parse Heating Type: Not found element.")
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.HeatingType = strings.TrimSpace(strings.Replace(htmlquery.InnerText(heatingTypeDoc), "Heating Features:", "", -1))
}

// ParseTotalParkingSpaces for parsing Total Parking Spaces
func (rc *MovotoCrawler) ParseTotalParkingSpaces() {
	//totalParkingSpacesDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"Parking Total:\")]")
	//if totalParkingSpacesDoc == nil {
	//	rc.Logger.Warn("Parse Total Parking Spaces: Not found element.")
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.TotalParkingSpaces = strings.TrimSpace(strings.Replace(htmlquery.InnerText(totalParkingSpacesDoc), "Parking Total:", "", -1))
}

// ParseLotFeatures for parsing Lot Features
func (rc *MovotoCrawler) ParseLotFeatures() {
	//var lotFeaturesList []string
	//lotFeaturesDocs := htmlquery.Find(rc.Doc, "//h4[contains(text(), \"Exterior and Lot Features\")]/following-sibling::ul[1]/li/text()")
	//if lotFeaturesDocs == nil {
	//	rc.Logger.Warn("Parse Lot Features: Not found element.")
	//	return
	//}
	//for _, v := range lotFeaturesDocs {
	//	lotFeaturesList = append(lotFeaturesList, htmlquery.InnerText(v))
	//}
	//rc.CrawlerSchemas.MovotoData.LotFeatures = strings.Join(lotFeaturesList, ", ")
}

// ParseProperySubType for parsing Propery SubType
func (rc *MovotoCrawler) ParseProperySubType() {
	//properySubTypeDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"Property Subtype:\")]")
	//if properySubTypeDoc == nil {
	//	rc.Logger.Warn("Parse Propery SubType: Not found element.")
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.ProperySubType = strings.TrimSpace(strings.Replace(htmlquery.InnerText(properySubTypeDoc), "Property Subtype:", "", -1))
}

// ParseConstructionMaterials for parsing Construction Materials
func (rc *MovotoCrawler) ParseConstructionMaterials() {
	//constructionMaterialsDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"Construction Materials:\")]")
	//if constructionMaterialsDoc == nil {
	//	rc.Logger.Warn("Parse Construction Materials: Not found element.")
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.ConstructionMaterials = strings.TrimSpace(strings.Replace(htmlquery.InnerText(constructionMaterialsDoc), "Construction Materials:", "", -1))
}

// ParseFoundation for parsing Foundation
func (rc *MovotoCrawler) ParseFoundation() {
	//foundationDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"Foundation Details:\")]")
	//if foundationDoc == nil {
	//	rc.Logger.Warn("Parse Foundation: Not found element.")
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.Foundation = strings.TrimSpace(strings.Replace(htmlquery.InnerText(foundationDoc), "Foundation Details:", "", -1))
}

// ParseSubdivision for parsing Subdivision
func (rc *MovotoCrawler) ParseSubdivision() {
	//SubdivisionDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"Subdivision:\")]")
	//if SubdivisionDoc == nil {
	//	rc.Logger.Warn("Parse Subdivision: Not found element.")
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.Subdivision = strings.TrimSpace(strings.Replace(htmlquery.InnerText(SubdivisionDoc), "Subdivision:", "", -1))
}

// ParseElementarySchool for parsing Elementary School
func (rc *MovotoCrawler) ParseElementarySchool() {
	//elementarySchoolDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"Elementary School:\")]")
	//if elementarySchoolDoc == nil {
	//	rc.Logger.Warn("Parse Elementary School: Not found element.")
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.ElementarySchool = strings.TrimSpace(strings.Replace(htmlquery.InnerText(elementarySchoolDoc), "Elementary School:", "", -1))
}

// ParseMiddleSchool for parsing Middle School
func (rc *MovotoCrawler) ParseMiddleSchool() {
	//middleSchoolDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"Middle School:\")]")
	//if middleSchoolDoc == nil {
	//	rc.Logger.Warn("Parse Middle School: Not found element.")
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.MiddleSchool = strings.TrimSpace(strings.Replace(htmlquery.InnerText(middleSchoolDoc), "Middle School:", "", -1))
}

// ParseHighSchool for parsing High School
func (rc *MovotoCrawler) ParseHighSchool() {
	//highSchoolDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"High School:\")]")
	//if highSchoolDoc == nil {
	//	rc.Logger.Warn("Parse High School: Not found element.")
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.HighSchool = strings.TrimSpace(strings.Replace(htmlquery.InnerText(highSchoolDoc), "High School:", "", -1))
}

// ParseDataSource for parsing Data Source
func (rc *MovotoCrawler) ParseDataSource() {
	//dataSourceDoc := htmlquery.FindOne(rc.Doc, "//div[contains(text(), \"Data Source:\")]/following-sibling::div")
	//if dataSourceDoc == nil {
	//	rc.Logger.Warn("Parse Data Source: Not found element.")
	//	return
	//}
	//rc.CrawlerSchemas.MovotoData.YearBuilt = htmlquery.InnerText(dataSourceDoc)
}

// ParsePictures for parsing pictures
func (rc *MovotoCrawler) ParsePictures() {
	//pictures := htmlquery.Find(rc.Doc, "//*[@class=\"main-carousel\"]//picture/img")
	//
	//if pictures == nil {
	//	rc.Logger.Warn("Parse Picture: Not found element.")
	//	return
	//}
	//
	//var picSlice []string
	//for _, pic := range pictures {
	//	picSlice = append(picSlice, htmlquery.SelectAttr(pic, "src"))
	//}
	//rc.CrawlerSchemas.MovotoData.Pictures = picSlice

}
