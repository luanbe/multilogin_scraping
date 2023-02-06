package realtor

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
	if err := rc.BaseSel.StartSelenium("realtor", rc.Proxy, viper.GetBool("crawler.realtor_crawler.proxy_status")); err != nil {
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

		pageSource, err := rc.BaseSel.WebDriver.PageSource()
		if err != nil {
			rc.BrowserTurnOff = true
			return err
		}

		if err := rc.ByPassVerifyHuman(pageSource, address); err != nil {
			return err
		}

		// NOTE: time to load source. Need to increase if data was not showing
		time.Sleep(viper.GetDuration("crawler.realtor_crawler.time_load_source") * time.Second)

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

func (rc *RealtorCrawler) CrawlSearchData(search string) (string, error) {
	// NOTE: We only take browser cookies when getting block from realtor website
	//cookies, err := rc.BaseSel.GetHttpCookies()
	//if err != nil {
	//	rc.Logger.Error(err.Error())
	//	return
	//}
	data := &schemas.RealtorSearchPageRes{}
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
		if result.FullAddress[0] == search {
			return result.MprID, nil
		}
	}
	return "", nil
}

func (rc *RealtorCrawler) ParseData(source string) error {
	doc, err := htmlquery.Parse(strings.NewReader(source))

	if err != nil {
		return err
	}
	rc.ParseBed(doc)
	rc.ParseBath(doc)
	rc.ParsePropertyStatus()
	rc.ParseFullBathrooms(doc)
	rc.ParseSF(doc)
	rc.ParseSalePrice(doc)
	rc.ParseEstPayment(doc)
	return nil
}

// ParseBed for crawling Bed data
func (rc *RealtorCrawler) ParseBed(doc *html.Node) {
	bedDoc := htmlquery.FindOne(
		doc,
		"//li[contains(@data-testid,\"property-meta-beds\")]/span/text()",
	)
	if bedDoc == nil {
		rc.Logger.Warn("Parse Bed: Not found Bed data")
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
func (rc *RealtorCrawler) ParseBath(doc *html.Node) {
	bedDoc := htmlquery.FindOne(
		doc,
		"//li[contains(@data-testid,\"property-meta-baths\")]/span/text()",
	)
	if bedDoc == nil {
		rc.Logger.Warn("Parse Bath: Not found Bath data")
		return
	}
	bedText := htmlquery.InnerText(bedDoc)
	if bedInt, err := strconv.Atoi(bedText); err != nil {
		rc.Logger.Error(fmt.Sprintf("Parse Bath: %v", err.Error()))
	} else {
		rc.CrawlerSchemas.RealtorData.Bed = bedInt
	}
}

// ParsePropertyStatus for checking Property Status
func (rc *RealtorCrawler) ParsePropertyStatus() {
	if rc.CrawlerSchemas.RealtorData.Bed > 0 || rc.CrawlerSchemas.RealtorData.Bath > 0 {
		rc.CrawlerSchemas.RealtorData.PropertyStatus = true
	}
}

// ParseFullBathrooms for crawling full bathrooms
func (rc *RealtorCrawler) ParseFullBathrooms(doc *html.Node) {
	fullBathroomsDoc := htmlquery.FindOne(doc, "//li[contains(text(), \"Full Bathrooms\")]/text()")
	if fullBathroomsDoc == nil {
		rc.Logger.Warn("Parse full bathrooms: Not found bathroom")
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

func (rc *RealtorCrawler) ParseSF(doc *html.Node) {
	sfDoc := htmlquery.FindOne(doc, "//li[contains(@data-testid,\"property-meta-lot-size\")]//span[@class=\"meta-value\"]/text()")
	if sfDoc == nil {
		rc.Logger.Warn("Parse SF: Not found SF")
		return
	}
	sfFloat, err := util2.ConvertToFloat(htmlquery.InnerText(sfDoc))

	if err != nil {
		rc.Logger.Error(fmt.Sprintf("Parse SF: %v", err.Error()))
		return
	}

	rc.CrawlerSchemas.RealtorData.FullBathrooms = sfFloat
}

func (rc *RealtorCrawler) ParseSalePrice(doc *html.Node) {
	salePriceDoc := htmlquery.FindOne(doc, "//div[@data-testid=\"list-price\"]//text()")
	if salePriceDoc == nil {
		rc.Logger.Warn("Parse Sale Price: Not found Sale Price")
		return
	}

	salePriceFloat, err := util2.ConvertToFloat(htmlquery.InnerText(salePriceDoc))

	if err != nil {
		rc.Logger.Warn(fmt.Sprintf("Parse Sale Price: %v", err.Error()))
		return
	}

	rc.CrawlerSchemas.RealtorData.SalesPrice = salePriceFloat
}

func (rc *RealtorCrawler) ParseEstPayment(doc *html.Node) {
	estPaymentDoc := htmlquery.FindOne(doc, "//div[@data-testid=\"list-estimate\"]/text()")
	if estPaymentDoc == nil {
		rc.Logger.Warn("Parse Est Payment: Not found Est Payment")
		return
	}
	rc.CrawlerSchemas.RealtorData.EstPayment = htmlquery.InnerText(estPaymentDoc)
}
