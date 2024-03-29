package realtor

import (
	"encoding/json"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
	"github.com/google/go-querystring/query"
	"github.com/icrowley/fake"
	"github.com/spf13/viper"
	"github.com/tebeka/selenium"
	"go.uber.org/zap"
	"golang.org/x/net/html"
	"gorm.io/gorm"
	"multilogin_scraping/app/schemas"
	"multilogin_scraping/crawlers"
	util2 "multilogin_scraping/pkg/utils"
	"strconv"
	"strings"
	"time"
)

type RealtorCrawler struct {
	BaseSel        *crawlers.BaseSelenium
	CRealtor       *colly.Collector
	Logger         *zap.Logger
	CrawlerBlocked bool
	BrowserTurnOff bool
	CrawlerSchemas *CrawlerSchemas
	Proxy          *util2.Proxy
	Doc            *html.Node
}

type CrawlerSchemas struct {
	RealtorData *schemas.RealtorData
	SearchReq   *schemas.RealtorSearchPageReq
}

const searchURL = "https://parser-external.geo.moveaws.com/suggest?%s"

func NewRealtorCrawler(
	db *gorm.DB,
	logger *zap.Logger,
	proxy *util2.Proxy,
) *RealtorCrawler {
	BaseSel := crawlers.NewBaseSelenium(logger)
	c := colly.NewCollector()
	userAgent := fake.UserAgent()
	c.UserAgent = userAgent

	return &RealtorCrawler{
		BaseSel:        BaseSel,
		CRealtor:       c,
		Logger:         logger,
		CrawlerBlocked: false,
		BrowserTurnOff: false,
		Proxy:          proxy,
		CrawlerSchemas: &CrawlerSchemas{
			RealtorData: &schemas.RealtorData{},
			SearchReq:   &schemas.RealtorSearchPageReq{},
		},
	}
}

// NewBrowser to start new selenium
func (rc *RealtorCrawler) NewBrowser() error {
	if err := rc.BaseSel.StartSelenium("realtor", rc.Proxy, viper.GetBool("crawler.realtor_crawler.proxy_status"), []string{"stealthfox"}); err != nil {
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
func (rc *RealtorCrawler) UserAgentBrowserToColly() error {
	userAgent, err := rc.BaseSel.WebDriver.ExecuteScript("return navigator.userAgent", nil)

	if err != nil {
		return err
	}

	if userAgent == nil {
		userAgent = fake.UserAgent()
	}

	rc.CRealtor.UserAgent = userAgent.(string)

	return nil
}

// RunRealtorCrawlerAPI with a loop to run crawler
func (rc *RealtorCrawler) RunRealtorCrawlerAPI(mprID string) error {
	rc.Logger.Info("Zillow Data is crawling...")
	rc.CrawlerSchemas.RealtorData.URL = fmt.Sprint(viper.GetString("crawler.realtor_crawler.url"), "realestateandhomes-detail/M", mprID)

	err := func(address string) error {
		if err := rc.BaseSel.WebDriver.Get(address); err != nil {
			return err
		}
		// NOTE: time to load source. Need to increase if data was not showing
		time.Sleep(viper.GetDuration("crawler.realtor_crawler.time_load_source") * time.Second)

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

	}(rc.CrawlerSchemas.RealtorData.URL)

	if err != nil {
		rc.Logger.Error(err.Error())
		rc.Logger.Error("Failed to crawl data")
		return err
		// TODO: Update error for crawling here
	}
	rc.Logger.Info("Completed to crawl data")
	return nil
}

// ByPassVerifyHuman to bypass verify from Realtor website
func (rc *RealtorCrawler) ByPassVerifyHuman(pageSource string, url string) error {
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
func (rc *RealtorCrawler) IsVerifyHuman(pageSource string) bool {
	if strings.Contains(pageSource, "Please verify you're a human to continue") || strings.Contains(pageSource, "Let's confirm you are human") {
		return true
	}
	return false
}

func (rc *RealtorCrawler) CrawlSearchData(crawlerSearchRes *schemas.CrawlerSearchRes) (string, error) {
	data := &schemas.RealtorSearchPageRes{}
	search := fmt.Sprintf(
		"%s %s %s %s",
		crawlerSearchRes.Search.Address,
		crawlerSearchRes.Search.City,
		crawlerSearchRes.Search.State,
		crawlerSearchRes.Search.Zipcode,
	)
	rc.CrawlerSchemas.SearchReq = &schemas.RealtorSearchPageReq{
		Input:     search,
		ClientID:  "rdc-home",
		Limit:     10,
		AreaTypes: "address",
	}
	searchPageQuery, err := query.Values(rc.CrawlerSchemas.SearchReq)
	if err != nil {
		return "nil", err
	}

	rc.CRealtor.OnError(func(r *colly.Response, err error) {
		rc.Logger.Error(fmt.Sprint("HTTP Status code:", r.StatusCode, "|URL:", r.Request.URL, "|Errors:", err))
		return
	})

	rc.CRealtor.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Content-Type", "application/json")
	})

	rc.CRealtor.OnResponse(func(r *colly.Response) {
		if err := json.Unmarshal(r.Body, data); err != nil {
			rc.Logger.Error(err.Error())
			return
		}
	})

	urlRun := fmt.Sprintf(searchURL, searchPageQuery.Encode())

	// NOTE: We only take browser cookies when getting block from realtor website
	//err = zc.CZillow.SetCookies(urlRun, cookies)
	//if err != nil {
	//	zc.SearchDataByCollyStatus = false
	//	zc.ShowLogError(err.Error())
	//	return
	//}
	if err := rc.CRealtor.Visit(urlRun); err != nil {
		return "", err
	}
	for _, result := range data.Autocomplete {
		if result.Line == crawlerSearchRes.Search.Address &&
			result.City == crawlerSearchRes.Search.City &&
			result.StateCode == crawlerSearchRes.Search.State &&
			result.PostalCode == crawlerSearchRes.Search.Zipcode {
			return result.MprID, nil
		}
	}
	return "", fmt.Errorf("not found data from address requested")
}

func (rc *RealtorCrawler) ParseData(source string) error {
	var err error
	if rc.Doc, err = htmlquery.Parse(strings.NewReader(source)); err != nil {
		return err
	}
	// Need to sure the data is existing
	sectionSummary := htmlquery.FindOne(rc.Doc, "//div[@id=\"section_summary\"]")
	if sectionSummary == nil {
		return fmt.Errorf("not found data from address requested")
	}

	if err := rc.ClickElements(); err != nil {
		return err
	}

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

func (rc *RealtorCrawler) ClickElements() error {

	rc.ClickPaymentCalculator()
	rc.ClickPropertyHistory()

	pageSource, err := rc.BaseSel.WebDriver.PageSource()

	if rc.Doc, err = htmlquery.Parse(strings.NewReader(pageSource)); err != nil {
		rc.Logger.Warn("Parse Principal Interest: Can not load page source")
		return err
	}
	return nil
}

func (rc *RealtorCrawler) ClickPaymentCalculator() {
	paymentCalculator, err := rc.BaseSel.WebDriver.FindElement(selenium.ByXPATH, "//section[@id=\"payment_calculator\"]")

	if err != nil {
		rc.Logger.Warn("Click Payment Calculator: Not found section payment_calculator")
		return
	}

	if err := paymentCalculator.Click(); err != nil {
		rc.Logger.Warn("Click Payment Calculator: Can not click payment_calculator")
		return
	}
	_, err = rc.BaseSel.WebDriver.ExecuteScript("document.getElementById(\"payment_calculator\").scrollIntoView();", nil)

	if err != nil {
		rc.Logger.Warn("Click Payment Calculator: Error on scoll to payment_calculator")
		return
	}
	time.Sleep(2 * time.Second)
}

func (rc *RealtorCrawler) ClickPropertyHistory() {
	propertyHistory, err := rc.BaseSel.WebDriver.FindElement(selenium.ByXPATH, "//section[@id=\"section_property_history\"]")

	if err != nil {
		rc.Logger.Warn("Click Property History: Not found section_property_history")
		return
	}

	if err := propertyHistory.Click(); err != nil {
		rc.Logger.Warn("Click Property History: Can not click section_property_history")
		return
	}
	_, err = rc.BaseSel.WebDriver.ExecuteScript("document.getElementById(\"section_property_history\").scrollIntoView();", nil)

	if err != nil {
		rc.Logger.Warn("Click Property History: Error on scoll to section_property_history")
		return
	}
	time.Sleep(2 * time.Second)
}

// ParseBed for crawling Bed data
func (rc *RealtorCrawler) ParseBed() {
	bedDoc := htmlquery.FindOne(
		rc.Doc,
		"//li[contains(@data-testid,\"property-meta-beds\")]/span/text()",
	)
	if bedDoc == nil {
		rc.Logger.Warn("Parse Bed: Not found element.")
		return
	}
	bedText := htmlquery.InnerText(bedDoc)
	if bedInt, err := strconv.Atoi(bedText); err != nil {
		rc.Logger.Error(fmt.Sprintf("Parse Bed: %v", err.Error()))
	} else {
		rc.CrawlerSchemas.RealtorData.Bed = bedInt
	}
}

// ParseBath for crawling Bath data
func (rc *RealtorCrawler) ParseBath() {
	bathDoc := htmlquery.FindOne(
		rc.Doc,
		"//li[contains(@data-testid,\"property-meta-baths\")]/span/text()",
	)
	if bathDoc == nil {
		rc.Logger.Warn("Parse Bath: Not found element.")
		return
	}
	bathText := htmlquery.InnerText(bathDoc)
	if bathInt, err := strconv.Atoi(bathText); err != nil {
		rc.Logger.Error(fmt.Sprintf("Parse Bath: %v", err.Error()))
	} else {
		rc.CrawlerSchemas.RealtorData.Bath = bathInt
	}
}

// ParsePropertyStatus for checking Property Status
func (rc *RealtorCrawler) ParsePropertyStatus() {
	if rc.CrawlerSchemas.RealtorData.Bed > 0 || rc.CrawlerSchemas.RealtorData.Bath > 0 {
		rc.CrawlerSchemas.RealtorData.PropertyStatus = true
	}
}

// ParseFullBathrooms for crawling full bathrooms
func (rc *RealtorCrawler) ParseFullBathrooms() {
	fullBathroomsDoc := htmlquery.FindOne(rc.Doc, "//li[contains(text(), \"Full Bathrooms\")]/text()")
	if fullBathroomsDoc == nil {
		rc.Logger.Warn("Parse full bathrooms: Not found element.")
		return
	}
	fullBathroomsText := htmlquery.InnerText(fullBathroomsDoc)
	fullBathroomsSlice := strings.Split(fullBathroomsText, ":")
	fullBathroomsText = fullBathroomsSlice[1]

	fullBathroomsFloat, err := util2.ConvertToFloat(fullBathroomsText)
	if err != nil {
		rc.Logger.Error(fmt.Sprintf("Parse full bathrooms: %v", err.Error()))
		return
	}
	rc.CrawlerSchemas.RealtorData.FullBathrooms = fullBathroomsFloat
}

// ParseSF for parsing SF
func (rc *RealtorCrawler) ParseSF() {
	sfDoc := htmlquery.FindOne(rc.Doc, "//li[contains(@data-testid,\"property-meta-sqft\")]//span[@class=\"meta-value\"]/text()")
	if sfDoc == nil {
		rc.Logger.Warn("Parse SF: Not found element.")
		return
	}
	sfFloat, err := util2.ConvertToFloat(htmlquery.InnerText(sfDoc))

	if err != nil {
		rc.Logger.Error(fmt.Sprintf("Parse SF: %v", err.Error()))
		return
	}

	rc.CrawlerSchemas.RealtorData.SF = sfFloat
}

// ParseSalePrice for parsing Sale Price
func (rc *RealtorCrawler) ParseSalePrice() {
	salePriceDoc := htmlquery.FindOne(rc.Doc, "//div[@data-testid=\"list-price\"]//text()")
	if salePriceDoc == nil {
		rc.Logger.Warn("Parse Sale Price: Not found element.")
		return
	}

	salePriceFloat, err := util2.ConvertToFloat(htmlquery.InnerText(salePriceDoc))

	if err != nil {
		rc.Logger.Warn(fmt.Sprintf("Parse Sale Price: %v", err.Error()))
		return
	}

	rc.CrawlerSchemas.RealtorData.SalesPrice = salePriceFloat
}

// ParseEstPayment for parsing Estimate Payment
func (rc *RealtorCrawler) ParseEstPayment() {
	estPaymentDoc := htmlquery.FindOne(rc.Doc, "//*[@data-testid=\"est-payment\"]//*[contains(text(), \"$\")]")
	if estPaymentDoc == nil {
		rc.Logger.Warn("Parse Est Payment: Not found element.")
		return
	}
	rc.CrawlerSchemas.RealtorData.EstPayment = htmlquery.InnerText(estPaymentDoc)
}

// ParsePrincipalInterest for parsing Principal Interest
// NOTE: This parsing will be loading page source again.
func (rc *RealtorCrawler) ParsePrincipalInterest() {
	principalInterestDoc := htmlquery.FindOne(rc.Doc, "//div[contains(text(), \"Principal & Interest\")]/following-sibling::div")

	if principalInterestDoc == nil {
		rc.Logger.Warn("Parse Principal Interest: Not found element.")
		return
	}

	rc.CrawlerSchemas.RealtorData.PrincipalInterest = htmlquery.InnerText(principalInterestDoc)
}

// ParseMortgageInsurance for parsing Mortgage Insurance
func (rc *RealtorCrawler) ParseMortgageInsurance() {
	mortgageInsuranceDoc := htmlquery.FindOne(rc.Doc, "//div[contains(text(), \"Mortgage Insurance\")]/following-sibling::div")
	if mortgageInsuranceDoc == nil {
		rc.Logger.Warn("Parse Mortgage Insurance: Not found element.")
		return
	}
	rc.CrawlerSchemas.RealtorData.MortgageInsurance = htmlquery.InnerText(mortgageInsuranceDoc)
}

// ParsePropertyTaxes for parsing Property Taxes
func (rc *RealtorCrawler) ParsePropertyTaxes() {
	propertyTaxesDoc := htmlquery.FindOne(rc.Doc, "//div[contains(text(), \"Property tax\")]/following-sibling::div")
	if propertyTaxesDoc == nil {
		rc.Logger.Warn("Parse Property Taxes: Not found element.")
		return
	}
	rc.CrawlerSchemas.RealtorData.PropertyTaxes = htmlquery.InnerText(propertyTaxesDoc)
}

// ParseHomeInsurance for parsing Home Insurance
func (rc *RealtorCrawler) ParseHomeInsurance() {
	homeInsuranceDoc := htmlquery.FindOne(rc.Doc, "//div[contains(text(), \"Home Insurance\")]/following-sibling::div")
	if homeInsuranceDoc == nil {
		rc.Logger.Warn("Parse Home Insurance: Not found element.")
		return
	}
	rc.CrawlerSchemas.RealtorData.HomeInsurance = htmlquery.InnerText(homeInsuranceDoc)
}

// ParseHOAFees for parsing HOA Fees
func (rc *RealtorCrawler) ParseHOAFees() {
	hoaFeesDoc := htmlquery.FindOne(rc.Doc, "//div[contains(text(), \"HOA fees\")]/following-sibling::div")
	if hoaFeesDoc == nil {
		rc.Logger.Warn("Parse HOA Fees: Not found element.")
		return
	}
	rc.CrawlerSchemas.RealtorData.HOAFee = htmlquery.InnerText(hoaFeesDoc)
}

// ParseOverview for parsing overview
func (rc *RealtorCrawler) ParseOverview() {
	overviewDoc := htmlquery.FindOne(rc.Doc, "//div[@id=\"section_property_details\"]")
	if overviewDoc == nil {
		rc.Logger.Warn("Parse Overview: Not found element.")
		return
	}
	rc.CrawlerSchemas.RealtorData.HOAFee = htmlquery.InnerText(overviewDoc)
}

// ParseSource for parsing source
func (rc *RealtorCrawler) ParseSource() {
	sourceDoc := htmlquery.FindOne(rc.Doc, "(//div[@id=\"content-property_history\"]//table)[1]/tbody//tr[position()=1]/td[last()]")
	if sourceDoc == nil {
		rc.Logger.Warn("Parse Source: Not found element.")
		return
	}
	rc.CrawlerSchemas.RealtorData.Source = htmlquery.InnerText(sourceDoc)
}

// ParsePropertyType for parsing Property Type
func (rc *RealtorCrawler) ParsePropertyType() {
	propertyTypeDoc := htmlquery.FindOne(rc.Doc, "//div[@id=\"section_summary\"]//span[contains(text(), \"Property Type\")]/parent::div/following-sibling::div")
	if propertyTypeDoc == nil {
		rc.Logger.Warn("Parse Property Type: Not found element.")
		return
	}
	rc.CrawlerSchemas.RealtorData.PropertyType = htmlquery.InnerText(propertyTypeDoc)
}

// ParseYearBuilt for parsing Year Built
func (rc *RealtorCrawler) ParseYearBuilt() {
	yearBuiltDoc := htmlquery.FindOne(rc.Doc, "//div[@id=\"section_summary\"]//span[contains(text(), \"Year Built\")]/parent::div/following-sibling::div")
	if yearBuiltDoc == nil {
		rc.Logger.Warn("Parse Year Built: Not found element.")
		return
	}
	rc.CrawlerSchemas.RealtorData.YearBuilt = htmlquery.InnerText(yearBuiltDoc)
}

// ParseNaturalGas for parsing Natural Gas
func (rc *RealtorCrawler) ParseNaturalGas() {
	naturalGasDoc := htmlquery.FindOne(rc.Doc, "//*[contains(text(), \"Natural Gas\") or contains(text(), \"natural gas\")]")
	if naturalGasDoc == nil {
		rc.Logger.Warn("Parse Natural Gas: Not found element.")
		return
	}
	rc.CrawlerSchemas.RealtorData.NaturalGas = true
}

// ParseCentralAir for parsing Central Air
func (rc *RealtorCrawler) ParseCentralAir() {
	centerAirDoc := htmlquery.FindOne(rc.Doc, "//*[contains(text(), \"Central Air\") or contains(text(), \"central air\")]")
	if centerAirDoc == nil {
		rc.Logger.Warn("Parse Central Air: Not found element.")
		return
	}
	rc.CrawlerSchemas.RealtorData.CentralAir = true
}

// ParseGarageSpaces for parsing Garage Spaces
func (rc *RealtorCrawler) ParseGarageSpaces() {
	garageSpacesDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"Garage Spaces\")]")
	if garageSpacesDoc == nil {
		rc.Logger.Warn("Parse Garage Spaces: Not found element.")
		return
	}
	rc.CrawlerSchemas.RealtorData.OfGarageSpaces = strings.TrimSpace(strings.Replace(htmlquery.InnerText(garageSpacesDoc), "Garage Spaces:", "", -1))
}

// ParseLotSize for parsing Lot Size
func (rc *RealtorCrawler) ParseLotSize() {
	lotSizeDoc := htmlquery.FindOne(rc.Doc, "//li[contains(@data-testid,\"property-meta-lot-size\")]//span[@class=\"meta-value\"]/text()")
	if lotSizeDoc == nil {
		rc.Logger.Warn("Parse Lot Size: Not found element.")
		return
	}
	lotSizeFloat, err := util2.ConvertToFloat(htmlquery.InnerText(lotSizeDoc))

	if err != nil {
		rc.Logger.Error(fmt.Sprintf("Parse Lot Size: %v", err.Error()))
		return
	}

	rc.CrawlerSchemas.RealtorData.LotSizeSF = lotSizeFloat
}

// ParseLotSizeAcres for parsing Lot Size Acres
func (rc *RealtorCrawler) ParseLotSizeAcres() {
	lotSizeAcresDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"Lot Size Acres\")]")
	if lotSizeAcresDoc == nil {
		rc.Logger.Warn("Parse Lot Size Acres: Not found element.")
		return
	}
	rc.CrawlerSchemas.RealtorData.OfGarageSpaces = strings.TrimSpace(strings.Replace(htmlquery.InnerText(lotSizeAcresDoc), "Lot Size Acres:", "", -1))
}

// ParseInteriorFeatures for parsing Interior Features
func (rc *RealtorCrawler) ParseInteriorFeatures() {
	var interiorFeaturesList []string
	interiorFeaturesDocs := htmlquery.Find(rc.Doc, "//h4[contains(text(), \"Interior Features\")]/following-sibling::ul[1]/li/text()")
	if interiorFeaturesDocs == nil {
		rc.Logger.Warn("Parse Interior Features: Not found element.")
		return
	}
	for _, v := range interiorFeaturesDocs {
		interiorFeaturesList = append(interiorFeaturesList, htmlquery.InnerText(v))
	}
	rc.CrawlerSchemas.RealtorData.InteriorFeatures = strings.Join(interiorFeaturesList, ", ")
}

// ParsePrimaryBedroomLevel for parsing Primary Bedroom Level
func (rc *RealtorCrawler) ParsePrimaryBedroomLevel() {
	primaryBedroomLevelDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"Bedrooms\")]")
	if primaryBedroomLevelDoc == nil {
		rc.Logger.Warn("Parse Primary Bedroom Level: Not found element.")
		return
	}
	rc.CrawlerSchemas.RealtorData.PrimaryBedroomLevel = strings.TrimSpace(strings.Replace(htmlquery.InnerText(primaryBedroomLevelDoc), "Bedrooms:", "", -1))
}

// ParseFlooringType for parsing Flooring Type
func (rc *RealtorCrawler) ParseFlooringType() {
	flooringTypeDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"Flooring:\")]")
	if flooringTypeDoc == nil {
		rc.Logger.Warn("Parse Flooring Type: Not found element.")
		return
	}
	rc.CrawlerSchemas.RealtorData.FlooringType = strings.TrimSpace(strings.Replace(htmlquery.InnerText(flooringTypeDoc), "Flooring:", "", -1))
}

// ParseHeatingType for parsing Heating Type
func (rc *RealtorCrawler) ParseHeatingType() {
	heatingTypeDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"Heating Features:\")]")
	if heatingTypeDoc == nil {
		rc.Logger.Warn("Parse Heating Type: Not found element.")
		return
	}
	rc.CrawlerSchemas.RealtorData.HeatingType = strings.TrimSpace(strings.Replace(htmlquery.InnerText(heatingTypeDoc), "Heating Features:", "", -1))
}

// ParseTotalParkingSpaces for parsing Total Parking Spaces
func (rc *RealtorCrawler) ParseTotalParkingSpaces() {
	totalParkingSpacesDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"Parking Total:\")]")
	if totalParkingSpacesDoc == nil {
		rc.Logger.Warn("Parse Total Parking Spaces: Not found element.")
		return
	}
	rc.CrawlerSchemas.RealtorData.TotalParkingSpaces = strings.TrimSpace(strings.Replace(htmlquery.InnerText(totalParkingSpacesDoc), "Parking Total:", "", -1))
}

// ParseLotFeatures for parsing Lot Features
func (rc *RealtorCrawler) ParseLotFeatures() {
	var lotFeaturesList []string
	lotFeaturesDocs := htmlquery.Find(rc.Doc, "//h4[contains(text(), \"Exterior and Lot Features\")]/following-sibling::ul[1]/li/text()")
	if lotFeaturesDocs == nil {
		rc.Logger.Warn("Parse Lot Features: Not found element.")
		return
	}
	for _, v := range lotFeaturesDocs {
		lotFeaturesList = append(lotFeaturesList, htmlquery.InnerText(v))
	}
	rc.CrawlerSchemas.RealtorData.LotFeatures = strings.Join(lotFeaturesList, ", ")
}

// ParseProperySubType for parsing Propery SubType
func (rc *RealtorCrawler) ParseProperySubType() {
	properySubTypeDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"Property Subtype:\")]")
	if properySubTypeDoc == nil {
		rc.Logger.Warn("Parse Propery SubType: Not found element.")
		return
	}
	rc.CrawlerSchemas.RealtorData.ProperySubType = strings.TrimSpace(strings.Replace(htmlquery.InnerText(properySubTypeDoc), "Property Subtype:", "", -1))
}

// ParseConstructionMaterials for parsing Construction Materials
func (rc *RealtorCrawler) ParseConstructionMaterials() {
	constructionMaterialsDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"Construction Materials:\")]")
	if constructionMaterialsDoc == nil {
		rc.Logger.Warn("Parse Construction Materials: Not found element.")
		return
	}
	rc.CrawlerSchemas.RealtorData.ConstructionMaterials = strings.TrimSpace(strings.Replace(htmlquery.InnerText(constructionMaterialsDoc), "Construction Materials:", "", -1))
}

// ParseFoundation for parsing Foundation
func (rc *RealtorCrawler) ParseFoundation() {
	foundationDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"Foundation Details:\")]")
	if foundationDoc == nil {
		rc.Logger.Warn("Parse Foundation: Not found element.")
		return
	}
	rc.CrawlerSchemas.RealtorData.Foundation = strings.TrimSpace(strings.Replace(htmlquery.InnerText(foundationDoc), "Foundation Details:", "", -1))
}

// ParseSubdivision for parsing Subdivision
func (rc *RealtorCrawler) ParseSubdivision() {
	SubdivisionDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"Subdivision:\")]")
	if SubdivisionDoc == nil {
		rc.Logger.Warn("Parse Subdivision: Not found element.")
		return
	}
	rc.CrawlerSchemas.RealtorData.Subdivision = strings.TrimSpace(strings.Replace(htmlquery.InnerText(SubdivisionDoc), "Subdivision:", "", -1))
}

// ParseElementarySchool for parsing Elementary School
func (rc *RealtorCrawler) ParseElementarySchool() {
	elementarySchoolDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"Elementary School:\")]")
	if elementarySchoolDoc == nil {
		rc.Logger.Warn("Parse Elementary School: Not found element.")
		return
	}
	rc.CrawlerSchemas.RealtorData.ElementarySchool = strings.TrimSpace(strings.Replace(htmlquery.InnerText(elementarySchoolDoc), "Elementary School:", "", -1))
}

// ParseMiddleSchool for parsing Middle School
func (rc *RealtorCrawler) ParseMiddleSchool() {
	middleSchoolDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"Middle School:\")]")
	if middleSchoolDoc == nil {
		rc.Logger.Warn("Parse Middle School: Not found element.")
		return
	}
	rc.CrawlerSchemas.RealtorData.MiddleSchool = strings.TrimSpace(strings.Replace(htmlquery.InnerText(middleSchoolDoc), "Middle School:", "", -1))
}

// ParseHighSchool for parsing High School
func (rc *RealtorCrawler) ParseHighSchool() {
	highSchoolDoc := htmlquery.FindOne(rc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"High School:\")]")
	if highSchoolDoc == nil {
		rc.Logger.Warn("Parse High School: Not found element.")
		return
	}
	rc.CrawlerSchemas.RealtorData.HighSchool = strings.TrimSpace(strings.Replace(htmlquery.InnerText(highSchoolDoc), "High School:", "", -1))
}

// ParseDataSource for parsing Data Source
func (rc *RealtorCrawler) ParseDataSource() {
	dataSourceDoc := htmlquery.FindOne(rc.Doc, "//div[contains(text(), \"Data Source:\")]/following-sibling::div")
	if dataSourceDoc == nil {
		rc.Logger.Warn("Parse Data Source: Not found element.")
		return
	}
	rc.CrawlerSchemas.RealtorData.YearBuilt = htmlquery.InnerText(dataSourceDoc)
}

// ParsePictures for parsing pictures
func (rc *RealtorCrawler) ParsePictures() {
	pictures := htmlquery.Find(rc.Doc, "//*[@class=\"main-carousel\"]//picture/img")

	if pictures == nil {
		rc.Logger.Warn("Parse Picture: Not found element.")
		return
	}

	var picSlice []string
	for _, pic := range pictures {
		picSlice = append(picSlice, htmlquery.SelectAttr(pic, "src"))
	}
	rc.CrawlerSchemas.RealtorData.Pictures = picSlice

}
