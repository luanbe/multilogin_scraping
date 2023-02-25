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
	"multilogin_scraping/app/models/entity"
	"multilogin_scraping/app/registry"
	"multilogin_scraping/app/schemas"
	"multilogin_scraping/app/service"
	"multilogin_scraping/crawlers"
	util2 "multilogin_scraping/pkg/utils"
	"strconv"
	"strings"
	"time"
)

type MovotoCrawler struct {
	BaseSel         *crawlers.BaseSelenium
	CMovoto         *colly.Collector
	Logger          *zap.Logger
	CrawlerBlocked  bool
	BrowserTurnOff  bool
	CrawlerSchemas  *CrawlerSchemas
	CrawlerTables   *CrawlerTables
	CrawlerServices CrawlerServices
	Proxy           *util2.Proxy
	Doc             *html.Node
}
type CrawlerServices struct {
	movotoService service.MovotoService
}

type CrawlerSchemas struct {
	SearchReq *schemas.MovotoSearchPageReq
}
type CrawlerTables struct {
	MovotoData *entity.Movoto
}

const searchURL = "https://www.movoto.com/api/v/search/?%s"

func NewMovotoCrawler(
	db *gorm.DB,
	logger *zap.Logger,
	proxy *util2.Proxy,
) *MovotoCrawler {
	BaseSel := crawlers.NewBaseSelenium(logger)
	c := colly.NewCollector()
	userAgent := fake.UserAgent()
	c.UserAgent = userAgent
	if viper.GetBool("crawler.movoto_crawler.proxy_status") == true {
		var host string
		if proxy.Type == "sock5" {
			host = "sock5"
		} else if proxy.Type == "HTTPS" {
			host = "https"
		} else {
			host = "http"
		}
		proxy := fmt.Sprintf("%s://%s:%s@%s:%s", host, proxy.Username, proxy.Password, proxy.Host, proxy.Port)
		if err := c.SetProxy(proxy); err != nil {
			logger.Warn(fmt.Sprintf("Can not set proxy on Colly with error: %s", err.Error()))
		}
	}
	return &MovotoCrawler{
		BaseSel:        BaseSel,
		CMovoto:        c,
		Logger:         logger,
		CrawlerBlocked: false,
		BrowserTurnOff: false,
		Proxy:          proxy,
		CrawlerSchemas: &CrawlerSchemas{
			SearchReq: &schemas.MovotoSearchPageReq{},
		},
		CrawlerTables: &CrawlerTables{
			&entity.Movoto{},
		},
		CrawlerServices: CrawlerServices{registry.RegisterMovotoService(db)},
	}
}

// NewBrowser to start new selenium
func (mc *MovotoCrawler) NewBrowser() error {
	if err := mc.BaseSel.StartSelenium("movoto", mc.Proxy, viper.GetBool("crawler.movoto_crawler.proxy_status"), []string{"mimic"}); err != nil {
		return err
	}
	// Disable image loading
	if viper.GetBool("crawler.disable_load_images") == true {
		if mc.BaseSel.Profile.BrowserName == "stealthfox" {
			if err := mc.BaseSel.FireFoxDisableImageLoading(); err != nil {
				return err
			}
		}
	}
	return nil
}

// UserAgentBrowserToColly for coping useragent from browser to colly
func (mc *MovotoCrawler) UserAgentBrowserToColly() error {
	userAgent, err := mc.BaseSel.WebDriver.ExecuteScript("return navigator.userAgent", nil)

	if err != nil {
		return err
	}

	if userAgent == nil {
		userAgent = fake.UserAgent()
	}

	mc.CMovoto.UserAgent = userAgent.(string)

	return nil
}

// RunMovotoCrawlerAPI with a loop to run crawler
func (mc *MovotoCrawler) RunMovotoCrawlerAPI(searchRes *schemas.MovotoSearchDataRes) error {
	mc.Logger.Info("Movoto Data is crawling...")
	mc.CrawlerTables.MovotoData.PropertyType = searchRes.PropertyType
	if searchRes.PropertyType != "" {
		mc.CrawlerTables.MovotoData.PropertyStatus = true
	}
	mc.CrawlerTables.MovotoData.Bath = float64(searchRes.Bath)
	mc.CrawlerTables.MovotoData.Bed = float64(searchRes.Bed)
	mc.CrawlerTables.MovotoData.LotSizeSF = string(rune(searchRes.LotSize))
	mc.CrawlerTables.MovotoData.HOAFee = string(rune(searchRes.Hoafee))
	mc.CrawlerTables.MovotoData.Pictures = searchRes.TnImgPath
	mc.CrawlerTables.MovotoData.URL = fmt.Sprint(viper.GetString("crawler.movoto_crawler.url"), searchRes.Path)
	mc.CrawlerTables.MovotoData.Address = fmt.Sprintf(
		"%s, %s, %s %s",
		searchRes.Geo.Address,
		searchRes.Geo.City,
		searchRes.Geo.State,
		searchRes.Geo.Zipcode,
	)

	err := func(address string) error {
		if err := mc.BaseSel.WebDriver.Get(address); err != nil {
			return err
		}
		// NOTE: time to load source. Need to increase if data was not showing
		time.Sleep(viper.GetDuration("crawler.movoto_crawler.time_load_source") * time.Second)

		pageSource, err := mc.BaseSel.WebDriver.PageSource()
		if err != nil {
			mc.BrowserTurnOff = true
			return err
		}

		if err := mc.ByPassVerifyHuman(pageSource, address); err != nil {
			return err
		}

		// TODO: Add Parse Data
		if err := mc.ParseData(pageSource); err != nil {
			return err
		}

		return nil

	}(mc.CrawlerTables.MovotoData.URL)

	if err != nil {
		mc.Logger.Error(err.Error())
		mc.Logger.Error("Failed to crawl data")
		return err
		// TODO: Update error for crawling here
	}
	mc.Logger.Info("Completed to crawl data")
	if err := mc.CrawlerServices.movotoService.AddMovoto(mc.CrawlerTables.MovotoData); err != nil {
		return err
	}
	mc.Logger.Info("Added/Updated record to Movoto Table")

	return nil
}

// ByPassVerifyHuman to bypass verify from Movoto website
func (mc *MovotoCrawler) ByPassVerifyHuman(pageSource string, url string) error {
	if mc.IsVerifyHuman(pageSource) == true {
		mc.CrawlerBlocked = true

	}
	if mc.CrawlerBlocked == true {
		for i := 0; i < 3; i++ {
			err := mc.BaseSel.WebDriver.Get(url)
			if err != nil {
				return err
			}
			pageSource, _ = mc.BaseSel.WebDriver.PageSource()
			if mc.IsVerifyHuman(pageSource) == false {
				mc.CrawlerBlocked = false
				return nil
			}
		}
		return fmt.Errorf("Crawler blocked for checking verify hunman")
	}
	return nil
}

// IsVerifyHuman to check website is blocking
func (mc *MovotoCrawler) IsVerifyHuman(pageSource string) bool {
	if strings.Contains(pageSource, "Please verify you're a human to continue") || strings.Contains(pageSource, "Let's confirm you are human") {
		return true
	}
	return false
}

func (mc *MovotoCrawler) CrawlSearchData(crawlerSearchRes *schemas.CrawlerSearchRes) (*schemas.MovotoSearchDataRes, error) {
	searchRes := &schemas.MovotoSearchPageRes{}
	movotoSearchData := &schemas.MovotoSearchDataRes{}
	path := fmt.Sprintf(
		"address %s %s %s %s",
		crawlerSearchRes.CrawlerRequest.Search.Address,
		crawlerSearchRes.CrawlerRequest.Search.City,
		crawlerSearchRes.CrawlerRequest.Search.State,
		crawlerSearchRes.CrawlerRequest.Search.Zipcode,
	)
	mc.CrawlerSchemas.SearchReq = &schemas.MovotoSearchPageReq{
		Path:              path,
		Trigger:           "mvtHeader",
		IncludeAllAddress: true,
		NewGeoSearch:      true,
	}
	searchPageQuery, err := query.Values(mc.CrawlerSchemas.SearchReq)
	if err != nil {
		return movotoSearchData, err
	}

	urlRun := fmt.Sprintf(searchURL, searchPageQuery.Encode())
	if err := mc.BaseSel.WebDriver.Get(urlRun); err != nil {
		return movotoSearchData, err
	}

	time.Sleep(time.Second * 2)
	pageSource, err := mc.BaseSel.WebDriver.PageSource()
	if err != nil {
		return movotoSearchData, err
	}

	if err := mc.ByPassVerifyHuman(pageSource, urlRun); err != nil {
		return movotoSearchData, err
	}
	doc, err := htmlquery.Parse(strings.NewReader(pageSource))
	if err != nil {
		return movotoSearchData, err
	}
	el := htmlquery.FindOne(doc, "//pre")

	jsonText := htmlquery.InnerText(el)
	if err := json.Unmarshal([]byte(jsonText), searchRes); err != nil {
		return movotoSearchData, err
	}

	if searchRes.Data.Listings == nil {
		return movotoSearchData, fmt.Errorf("not found data from address requested")
	}

	for _, v := range searchRes.Data.Listings {
		if v.Geo.Address == crawlerSearchRes.CrawlerRequest.Search.Address &&
			v.Geo.City == crawlerSearchRes.CrawlerRequest.Search.City &&
			v.Geo.State == crawlerSearchRes.CrawlerRequest.Search.State &&
			v.Geo.Zipcode == crawlerSearchRes.CrawlerRequest.Search.Zipcode {
			return &v, nil
		}
	}

	return movotoSearchData, fmt.Errorf("not found data from address requested")
}

func (mc *MovotoCrawler) ParseData(source string) error {
	var err error
	if mc.Doc, err = htmlquery.Parse(strings.NewReader(source)); err != nil {
		return err
	}
	movotoJsonRes := mc.ParseJsonData()

	if movotoJsonRes.PageData.Features != nil {
		for _, feature := range movotoJsonRes.PageData.Features {
			// Interior value
			if feature.Name == "Interior" {
				for _, interior := range feature.Value {
					if interior.Name == "Bathrooms" {
						mc.CrawlerTables.MovotoData.FullBathrooms, err = strconv.ParseFloat(interior.Value[0].Value, 64)
						if err != nil {
							mc.Logger.Warn(fmt.Sprintf("Parse Full Bathrooms:%s", err.Error()))
						}

						if mc.CrawlerTables.MovotoData.FullBathrooms != 0 {
							mc.CrawlerTables.MovotoData.HalfBathrooms = mc.CrawlerTables.MovotoData.FullBathrooms / 2
						}

					}
				}
			}
			// Exterior value
			if feature.Name == "Exterior" {
				for _, exterior := range feature.Value {
					if exterior.Name == "Parking" {
						for _, parking := range exterior.Value {
							if parking.Name == "# Covered Spaces" {
								mc.CrawlerTables.MovotoData.TotalParkingSpaces = parking.Value
							}
						}
					}
				}
			}
		}
	}
	mc.CrawlerTables.MovotoData.SF = float64(movotoJsonRes.PageData.SqftTotal)
	mc.CrawlerTables.MovotoData.SalesPrice = float64(movotoJsonRes.PageData.ListPrice)

	var pictures []string
	if len(movotoJsonRes.PageData.CategorizedPhotos) > 0 {
		for _, pic := range movotoJsonRes.PageData.CategorizedPhotos {
			if len(pic.Photos) > 0 {
				for _, photo := range pic.Photos {
					pictures = append(pictures, photo.URL)
				}
			}
		}
	}
	mc.CrawlerTables.MovotoData.Pictures = strings.Join(pictures, ", ")
	mc.CrawlerTables.MovotoData.Overview = movotoJsonRes.PageData.Description
	//mc.ParseFullBathrooms()
	//mc.ParseSF()
	//mc.ParseSalePrice()
	mc.ParseEstPayment()
	mc.ParsePrincipalInterest()
	mc.ParseMortgageInsurance()
	mc.ParsePropertyTaxes()
	mc.ParseHomeInsurance()
	mc.ParseEstimatedSalesRange()
	mc.ParseMLS()
	mc.ParseYearBuilt()
	mc.ParseLotSizeAcres()
	mc.ParseAppliances()
	mc.ParsePropertySubtype()
	mc.ParseFoundation()
	mc.ParseNewConstruction()
	return nil
}

func (mc *MovotoCrawler) ParseJsonData() *schemas.MovotoJsonRes {
	var movotoJsonRes *schemas.MovotoJsonRes
	jsDoc := htmlquery.FindOne(mc.Doc, "//script[contains(text(), \"__INITIAL_STATE__ \")]")
	jsText := htmlquery.InnerText(jsDoc)
	// Remove Javascript code
	if jsText != "" {
		jsText = strings.Split(jsText, "__INITIAL_STATE__ = ")[1]
		jsText = strings.Split(jsText, ";\n\t\t\t\twindow.startTime")[0]
	}

	// Convert string to json
	if err := json.Unmarshal([]byte(jsText), &movotoJsonRes); err != nil {
		mc.Logger.Warn(fmt.Sprintf("\"Parse Json Data: Got errors %s\"", err.Error()))
	}

	return movotoJsonRes
}

// ParseEstPayment for parsing Estimate Payment
func (mc *MovotoCrawler) ParseEstPayment() {
	estPaymentDoc := htmlquery.FindOne(mc.Doc, "//div[@comp=\"propertyTitle\"]//span[contains(text(), \"Estimate\")]/a")
	if estPaymentDoc == nil {
		mc.Logger.Warn("Parse Est Payment: Not found element.")
		return
	}
	mc.CrawlerTables.MovotoData.EstPayment = htmlquery.InnerText(estPaymentDoc)
}

// ParsePrincipalInterest for parsing Principal Interest
// NOTE: This parsing will be loading page source again.
func (mc *MovotoCrawler) ParsePrincipalInterest() {
	principalInterestDoc := htmlquery.FindOne(mc.Doc, "//div[contains(text(), \"Principal & Interest\")]/following-sibling::div/span")

	if principalInterestDoc == nil {
		mc.Logger.Warn("Parse Principal Interest: Not found element.")
		return
	}

	mc.CrawlerTables.MovotoData.PrincipalInterest = htmlquery.InnerText(principalInterestDoc)
}

// ParseMortgageInsurance for parsing Mortgage Insurance
func (mc *MovotoCrawler) ParseMortgageInsurance() {
	mortgageInsuranceDoc := htmlquery.FindOne(mc.Doc, "//div[@comp=\"propertyTitle\"]//span[contains(text(), \"Mortgage\")]/a")
	if mortgageInsuranceDoc == nil {
		mc.Logger.Warn("Parse Mortgage Insurance: Not found element.")
		return
	}
	mc.CrawlerTables.MovotoData.MortgageInsurance = htmlquery.InnerText(mortgageInsuranceDoc)
}

// ParsePropertyTaxes for parsing Property Taxes
func (mc *MovotoCrawler) ParsePropertyTaxes() {
	propertyTaxesDoc := htmlquery.FindOne(mc.Doc, "//section[@id=\"propertyMortgagePanel\"]//div[contains(text(), \"Taxes\")]/following-sibling::div/span")
	if propertyTaxesDoc == nil {
		mc.Logger.Warn("Parse Property Taxes: Not found element.")
		return
	}
	mc.CrawlerTables.MovotoData.PropertyTaxes = htmlquery.InnerText(propertyTaxesDoc)
}

// ParseHomeInsurance for parsing Home Insurance
// TODO: Still not found data element and will find later
func (mc *MovotoCrawler) ParseHomeInsurance() {
	homeInsuranceDoc := htmlquery.FindOne(mc.Doc, "//section[@id=\"propertyMortgagePanel\"]//a[@data-id=\"homeInsurance\"]/div[2]/span")
	if homeInsuranceDoc == nil {
		mc.Logger.Warn("Parse Home Insurance: Not found element.")
		return
	}
	mc.CrawlerTables.MovotoData.HomeInsurance = htmlquery.InnerText(homeInsuranceDoc)
}

// ParseEstimatedSalesRange for parsing Estimated Sale Range Minimum and Maximum
func (mc *MovotoCrawler) ParseEstimatedSalesRange() {
	estSaleRangeDoc := htmlquery.FindOne(mc.Doc, "//div[contains(text(), \"Estimated List Price\")]/following-sibling::div")
	if estSaleRangeDoc == nil {
		mc.Logger.Warn("Parse Estimated Sale Range: Not found element.")
		return
	}
	estSaleRangeText := htmlquery.InnerText(estSaleRangeDoc)
	estSaleRange := strings.Split(estSaleRangeText, "-")
	if len(estSaleRange) == 2 {
		mc.CrawlerTables.MovotoData.EstimatedSalesRangeMinimum = strings.TrimSpace(estSaleRange[0])
		mc.CrawlerTables.MovotoData.EstimatedSalesRangeMax = strings.TrimSpace(estSaleRange[1])
	}
}

// ParseMLS for parsing # MLS
func (mc *MovotoCrawler) ParseMLS() {
	mlsDoc := htmlquery.FindOne(mc.Doc, "//section[@id=\"propertyDetailPanel\"]//span[contains(text(), \"MLS #\")]/following-sibling::div")
	if mlsDoc == nil {
		mc.Logger.Warn("Parse #MLS: Not found element.")
		return
	}
	mc.CrawlerTables.MovotoData.MLS = htmlquery.InnerText(mlsDoc)
}

// ParseYearBuilt for parsing Year Built
func (mc *MovotoCrawler) ParseYearBuilt() {
	yearBuiltDoc := htmlquery.FindOne(mc.Doc, "//section[@id=\"propertyDetailPanel\"]//*[contains(text(), \"Year Built\")]/following-sibling::div")
	if yearBuiltDoc == nil {
		mc.Logger.Warn("Parse Year Built: Not found element.")
		return
	}
	mc.CrawlerTables.MovotoData.YearBuilt = htmlquery.InnerText(yearBuiltDoc)
}

// ParseLotSizeAcres for parsing Lot Size Acres
func (mc *MovotoCrawler) ParseLotSizeAcres() {
	lotSizeAcresDoc := htmlquery.FindOne(mc.Doc, "//ul[@class=\"feature-text-list\"]//li[contains(text(), \"Lot Size Acres\")]")
	if lotSizeAcresDoc == nil {
		mc.Logger.Warn("Parse Lot Size Acres: Not found element.")
		return
	}
	mc.CrawlerTables.MovotoData.LotSizeAcres = strings.TrimSpace(strings.Replace(htmlquery.InnerText(lotSizeAcresDoc), "Lot Size Acres:", "", -1))
}

// ParseAppliances for parsing Appliances
func (mc *MovotoCrawler) ParseAppliances() {
	appliancesDoc := htmlquery.FindOne(mc.Doc, "//section[@id=\"propertyDetailPanel\"]//*[contains(text(), \"Appliances\")]/following-sibling::div")
	if appliancesDoc == nil {
		mc.Logger.Warn("Parse Appliances: Not found element.")
		return
	}
	mc.CrawlerTables.MovotoData.Appliances = htmlquery.InnerText(appliancesDoc)
}

// ParsePropertySubtype for parsing Property Subtype
func (mc *MovotoCrawler) ParsePropertySubtype() {
	propertySubTypesDoc := htmlquery.FindOne(mc.Doc, "//section[@id=\"propertyDetailPanel\"]//*[contains(text(), \"SubType\")]/following-sibling::div")
	if propertySubTypesDoc == nil {
		mc.Logger.Warn("Parse Property Subtype: Not found element.")
		return
	}
	mc.CrawlerTables.MovotoData.ProperySubType = htmlquery.InnerText(propertySubTypesDoc)
}

// ParseFoundation for parsing Foundation
func (mc *MovotoCrawler) ParseFoundation() {
	foundationDoc := htmlquery.FindOne(mc.Doc, "//section[@id=\"propertyDetailPanel\"]//*[contains(text(), \"Foundation\")]/following-sibling::div")
	if foundationDoc == nil {
		mc.Logger.Warn("Parse Foundation: Not found element.")
		return
	}
	mc.CrawlerTables.MovotoData.Foundation = htmlquery.InnerText(foundationDoc)
}

// ParseNewConstruction for parsing New Construction
func (mc *MovotoCrawler) ParseNewConstruction() {
	newConstructionDoc := htmlquery.FindOne(mc.Doc, "//section[@id=\"propertyDetailPanel\"]//*[contains(text(), \"Construction\")]/following-sibling::div")
	if newConstructionDoc == nil {
		mc.Logger.Warn("Parse New Construction: Not found element.")
		return
	}
	mc.CrawlerTables.MovotoData.NewConstruction = htmlquery.InnerText(newConstructionDoc)
}
