package zillow

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/net/html"
	"gorm.io/gorm"
	"multilogin_scraping/app/registry"
	"multilogin_scraping/app/schemas"
	"multilogin_scraping/app/service"
	"multilogin_scraping/crawlers"
	util2 "multilogin_scraping/pkg/utils"
	"strconv"
	"strings"
	"time"

	"github.com/icrowley/fake"
	"github.com/spf13/viper"

	"multilogin_scraping/app/models/entity"

	"github.com/antchfx/htmlquery"
	"github.com/gocolly/colly"
	"github.com/tebeka/selenium"
)

type ZillowCrawler struct {
	WebDriver               selenium.WebDriver
	BaseSel                 *crawlers.BaseSelenium
	Profile                 *crawlers.Profile
	CZillow                 *colly.Collector
	SearchPageReq           *schemas.ZillowSearchPageReq
	CrawlerTables           *CrawlerTables
	CrawlerServices         CrawlerServices
	Logger                  *zap.Logger
	OnlyHistoryTable        bool
	CrawlerBlocked          bool
	BrowserTurnOff          bool
	Maindb3List             []*entity.ZillowMaindb3Address
	SearchDataByCollyStatus bool
	Proxy                   *util2.Proxy
}

type CrawlerTables struct {
	Maindb3                *entity.ZillowMaindb3Address
	ZillowData             *entity.Zillow
	ZillowSearchData       []*entity.Zillow
	ZillowPriceHistory     []*entity.ZillowPriceHistory
	ZillowPublicTaxHistory []*entity.ZillowPublicTaxHistory
	MapBounds              *schemas.MapBounds
	//Zillow *entry.Zillow
}

type CrawlerServices struct {
	ZillowService  service.ZillowService
	Maindb3Service service.Maindb3Service
}

const searchURL = "https://www.zillow.com/search/GetSearchPageState.htm?searchQueryState=%s&wants={\"cat1\":[\"listResults\",\"mapResults\"],\"cat2\":[\"total\"],\"regionResults\":[\"total\"]}&requestId=5"

func NewZillowCrawler(
	db *gorm.DB,
	maindb3List []*entity.ZillowMaindb3Address,
	logger *zap.Logger,
	onlyHistoryTable bool,
	proxy *util2.Proxy,
) *ZillowCrawler {
	BaseSel := crawlers.NewBaseSelenium(logger)
	c := colly.NewCollector()
	userAgent := fake.UserAgent()
	c.UserAgent = userAgent

	return &ZillowCrawler{
		WebDriver: BaseSel.WebDriver,
		BaseSel:   BaseSel,
		Profile:   BaseSel.Profile,
		CZillow:   c,
		CrawlerTables: &CrawlerTables{
			&entity.ZillowMaindb3Address{},
			&entity.Zillow{},
			[]*entity.Zillow{},
			[]*entity.ZillowPriceHistory{},
			[]*entity.ZillowPublicTaxHistory{},
			&schemas.MapBounds{},
		},
		CrawlerServices:  CrawlerServices{registry.RegisterZillowService(db), registry.RegisterMaindb3Service(db)},
		Logger:           logger,
		OnlyHistoryTable: onlyHistoryTable,
		CrawlerBlocked:   false,
		BrowserTurnOff:   false,
		Maindb3List:      maindb3List,
		Proxy:            proxy,
	}
}

// NewBrowser to start new selenium
func (zc *ZillowCrawler) NewBrowser() error {
	if err := zc.BaseSel.StartSelenium("zillow", zc.Proxy, viper.GetBool("crawler.realtor_crawler.proxy_status"), []string{"stealthfox"}); err != nil {
		return err
	}
	// Disable image loading
	if viper.GetBool("crawler.disable_load_images") == true {
		if zc.BaseSel.Profile.BrowserName == "stealthfox" {
			if err := zc.BaseSel.FireFoxDisableImageLoading(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (zc *ZillowCrawler) GetURLCrawling() string {
	var address string
	if zc.CrawlerTables.Maindb3.OwnerAddress != "" && zc.CrawlerTables.Maindb3.OwnerCityState != "" {
		address = fmt.Sprint(strings.TrimSpace(zc.CrawlerTables.Maindb3.OwnerAddress), ", ", zc.CrawlerTables.Maindb3.OwnerCityState)
	} else {
		ownerCityState := fmt.Sprint(zc.CrawlerTables.Maindb3.AddressCity, ", ", zc.CrawlerTables.Maindb3.AddressState)
		address = fmt.Sprint(strings.TrimSpace(zc.CrawlerTables.Maindb3.AddressStreet), ", ", ownerCityState)
	}

	address = strings.Replace(address, " ", "-", -1)
	address = strings.Replace(address, "--", "-", -1)
	return fmt.Sprint(viper.GetString("crawler.zillow_crawler.url"), address, "_rb/")

	//// NOTE: For testing data
	//return "https://www.zillow.com/homes/PO-BOX-2073,-LAKE-DALLAS,-TX_rb/"
}

func (zc *ZillowCrawler) ShowLogError(mes string) {
	zc.Logger.Error(mes, zap.Uint64("mainDBID", zc.CrawlerTables.Maindb3.ID), zap.String("URL", zc.GetURLCrawling()))
}

func (zc *ZillowCrawler) ShowLogInfo(mes string) {
	zc.Logger.Info(mes, zap.Uint64("mainDBID", zc.CrawlerTables.Maindb3.ID), zap.String("URL", zc.GetURLCrawling()))
}

func (zc *ZillowCrawler) RunZillowCrawler() error {
	for {
		err := func() error {
			for _, maindb3 := range zc.Maindb3List {
				zc.CrawlerTables.Maindb3 = maindb3
				zc.CrawlerTables.ZillowData = &entity.Zillow{}
				zc.ShowLogInfo("Zillow Data is crawling...")
				zc.CrawlerTables.ZillowData.URL = zc.GetURLCrawling()
				if err := zc.CrawlAddress(zc.CrawlerTables.ZillowData.URL); err != nil {
					zc.ShowLogError(err.Error())
					zc.ShowLogError("Failed to crawl data")
					if zc.CrawlerBlocked == true {
						return err
					}
					if err = zc.CrawlerServices.Maindb3Service.UpdateStatus(zc.CrawlerTables.Maindb3, viper.GetString("crawler.crawler_status.failed")); err != nil {
						return err
					}

				}
				zc.ShowLogInfo("Completed to crawl data")
				time.Sleep(time.Second * viper.GetDuration("crawler.zillow_crawler.crawl_next_time"))
			}
			return nil
		}()
		if err != nil {
			return err
		}
		return nil
	}
}

func (zc *ZillowCrawler) RunZillowCrawlerAPI(crawlerSearchRes *schemas.CrawlerSearchRes) error {
	zc.ShowLogInfo("Zillow Data is crawling...")
	zillowRootURL := viper.GetString("crawler.zillow_crawler.url")
	zc.CrawlerTables.ZillowData.URL = fmt.Sprintf(
		"%s%s-%s-%s-%s",
		zillowRootURL,
		strings.Replace(crawlerSearchRes.CrawlerRequest.Search.Address, " ", "-", -1),
		strings.Replace(crawlerSearchRes.CrawlerRequest.Search.City, " ", "-", -1),
		crawlerSearchRes.CrawlerRequest.Search.State,
		crawlerSearchRes.CrawlerRequest.Search.Zipcode,
	)

	err := func() error {
		if err := zc.BaseSel.WebDriver.Get(zc.CrawlerTables.ZillowData.URL); err != nil {
			return err
		}
		// NOTE: time to load source. Need to increase if data was not showing
		time.Sleep(viper.GetDuration("crawler.zillow_crawler.time_load_source") * time.Second)

		pageSource, err := zc.BaseSel.WebDriver.PageSource()
		if err != nil {
			zc.BrowserTurnOff = true
			return err
		}

		if err := zc.ByPassVerifyHuman(pageSource, zc.CrawlerTables.ZillowData.URL); err != nil {
			return err
		}

		// TODO: Add Parse Data
		if err := zc.ParseData(pageSource); err != nil {
			return err
		}

		return nil

	}()
	if err != nil {
		zc.Logger.Error(err.Error())
		zc.Logger.Error("Failed to crawl data")
		return err
		// TODO: Update error for crawling here
	}
	zc.Logger.Info("Completed to crawl data")
	if err := zc.CrawlerServices.ZillowService.AddZillow(zc.CrawlerTables.ZillowData); err != nil {
		return err
	}
	zc.ShowLogInfo("Added/Updated record to Zillow Table")

	return nil
}

func (zc *ZillowCrawler) ByPassVerifyHuman(pageSource string, url string) error {
	if zc.IsVerifyHuman(pageSource) == true {
		zc.CrawlerBlocked = true

	}
	if zc.CrawlerBlocked == true {
		for i := 0; i < 3; i++ {
			err := zc.WebDriver.Get(url)
			if err != nil {
				return err
			}
			pageSource, _ = zc.WebDriver.PageSource()
			if zc.IsVerifyHuman(pageSource) == false {
				zc.CrawlerBlocked = false
				return nil
			}
		}
		return fmt.Errorf("Crawler blocked for checking verify hunman")
	}
	return nil
}
func (zc *ZillowCrawler) IsVerifyHuman(pageSource string) bool {
	if strings.Contains(pageSource, "Please verify you're a human to continue") || strings.Contains(pageSource, "Let's confirm you are human") {
		return true
	}
	return false
}

func (zc *ZillowCrawler) CrawlSearchDataByColly() {
	zc.SearchDataByCollyStatus = true
	cookies, err := zc.BaseSel.GetHttpCookies()
	if err != nil {
		zc.ShowLogError(err.Error())
		zc.SearchDataByCollyStatus = false
		return
	}

	if zc.CrawlerTables.MapBounds == nil {
		zc.ShowLogError("not found map bounds data")
		zc.SearchDataByCollyStatus = false
		return
	}

	zc.CrawlerTables.ZillowSearchData = []*entity.Zillow{}

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
		zc.SearchDataByCollyStatus = false
		zc.ShowLogError(err.Error())
		return
	}
	zc.SearchPageReq.MapBounds = zc.CrawlerTables.MapBounds

	searchPageJson, err := json.Marshal(zc.SearchPageReq)
	if err != nil {
		zc.SearchDataByCollyStatus = false
		zc.ShowLogError(err.Error())
		return
	}

	zc.CZillow.OnResponse(func(r *colly.Response) {
		currentUrl := r.Ctx.Get("currentURL")
		data := &schemas.ZillowSearchPageRes{}
		if err := json.Unmarshal(r.Body, data); err != nil {
			zc.ShowLogInfo("Not found Json data when crawling by colly")
			zc.SearchDataByCollyStatus = false
			return
		}
		if len(data.Cat1.SearchResults.RelaxedResults) > 0 {
			for _, result := range data.Cat1.SearchResults.RelaxedResults {
				zillowSearchData := zc.CrawlZillowSearchRelaxedData(result)
				zc.CrawlerTables.ZillowSearchData = append(zc.CrawlerTables.ZillowSearchData, zillowSearchData)
			}
			zc.ShowLogInfo(fmt.Sprint("Found relaxedResults for crawling on Current Page: ", currentUrl))
		}

		if len(data.Cat1.SearchResults.ListResults) > 0 {
			for _, result := range data.Cat1.SearchResults.ListResults {
				zillowSearchData := zc.CrawlZillowSearchResultData(result)
				zc.CrawlerTables.ZillowSearchData = append(zc.CrawlerTables.ZillowSearchData, zillowSearchData)
			}

			// Crawling data on Next Page
			zc.SearchPageReq.Pagination.CurrentPage += 1
			searchNextPageJson, err := json.Marshal(zc.SearchPageReq)
			if err != nil {
				zc.SearchDataByCollyStatus = false
				zc.ShowLogError(err.Error())
				return
			}
			urlNextPage := fmt.Sprintf(searchURL, string(searchNextPageJson))
			zc.ShowLogInfo(fmt.Sprint("Found listResults for crawling on Current Page: ", currentUrl))

			err = zc.CZillow.SetCookies(urlNextPage, cookies)
			if err != nil {
				zc.SearchDataByCollyStatus = false
				zc.ShowLogError(err.Error())
				return
			}
			if err == nil {
				r.Request.Visit(urlNextPage)
			}
		}

	})

	zc.CZillow.OnError(func(r *colly.Response, err error) {
		zc.ShowLogError(fmt.Sprint("HTTP Status code:", r.StatusCode, "|URL:", r.Request.URL, "|Errors:", err))
		zc.SearchDataByCollyStatus = false
		return
	})
	zc.CZillow.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Content-Type", "application/json")
		r.Ctx.Put("currentURL", r.URL.String())
	})
	urlRun := fmt.Sprintf(searchURL, string(searchPageJson))
	err = zc.CZillow.SetCookies(urlRun, cookies)
	if err != nil {
		zc.SearchDataByCollyStatus = false
		zc.ShowLogError(err.Error())
		return
	}
	zc.CZillow.Visit(urlRun)
}
func (zc *ZillowCrawler) CrawlSearchData() error {

	if zc.CrawlerTables.MapBounds == nil {
		return fmt.Errorf("not found map bounds data")
	}
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
		return err
	}
	zc.SearchPageReq.MapBounds = zc.CrawlerTables.MapBounds

	searchPageJson, err := json.Marshal(zc.SearchPageReq)
	if err != nil {
		return err
	}

	urlRun := fmt.Sprintf(searchURL, string(searchPageJson))
	zc.CrawlerTables.ZillowSearchData = []*entity.Zillow{}
	var nextURL string
	var firstRun bool
	firstRun = true
	for {
		if nextURL == "" && firstRun == true {
			url, err := zc.CrawlNextSearchData(urlRun)
			if err != nil {
				return err
			}
			if url == "" {
				break
			}
			nextURL = url
			firstRun = false
			continue
		}
		if nextURL != "" {
			url, err := zc.CrawlNextSearchData(nextURL)
			if err != nil {
				return err
			}
			if url == "" {
				break
			}
			nextURL = url
			continue
		}
		break
	}
	return nil
}
func (zc *ZillowCrawler) CrawlNextSearchData(urlRun string) (string, error) {
	searchPageRes := &schemas.ZillowSearchPageRes{}
	if err := zc.WebDriver.Get(urlRun); err != nil {
		return "", err
	}
	time.Sleep(time.Second * 3)
	rawdataTab, err := zc.WebDriver.FindElement(selenium.ByID, "rawdata-tab")
	if err != nil {
		return "", err
	}
	if err := rawdataTab.Click(); err != nil {
		return "", err
	}
	time.Sleep(time.Second * 2)
	pageSource, err := zc.WebDriver.PageSource()
	if err != nil {
		return "", err
	}
	if err := zc.ByPassVerifyHuman(pageSource, urlRun); err != nil {
		return "", err
	}
	doc, err := htmlquery.Parse(strings.NewReader(pageSource))
	if err != nil {
		return "", err
	}
	el := htmlquery.FindOne(doc, "//pre[@class='data']")

	jsonText := htmlquery.InnerText(el)
	if err := json.Unmarshal([]byte(jsonText), searchPageRes); err != nil {
		return "", err
	}

	currentUrl, _ := zc.WebDriver.CurrentURL()

	if len(searchPageRes.Cat1.SearchResults.RelaxedResults) > 0 {
		for _, result := range searchPageRes.Cat1.SearchResults.RelaxedResults {
			zillowSearchData := zc.CrawlZillowSearchRelaxedData(result)
			zc.CrawlerTables.ZillowSearchData = append(zc.CrawlerTables.ZillowSearchData, zillowSearchData)
		}
		zc.ShowLogInfo(fmt.Sprint("Found relaxedResults for crawling on Current Page: ", currentUrl))
	}

	if len(searchPageRes.Cat1.SearchResults.ListResults) > 0 {
		for _, result := range searchPageRes.Cat1.SearchResults.ListResults {
			zillowSearchData := zc.CrawlZillowSearchResultData(result)
			zc.CrawlerTables.ZillowSearchData = append(zc.CrawlerTables.ZillowSearchData, zillowSearchData)
		}
		// Crawling data on Next Page
		zc.SearchPageReq.Pagination.CurrentPage += 1
		searchNextPageJson, err := json.Marshal(zc.SearchPageReq)
		if err != nil {
			return "", err
		}
		nextURL := fmt.Sprintf(searchURL, string(searchNextPageJson))

		zc.ShowLogInfo(fmt.Sprint("Found listResults for crawling on Current Page: ", currentUrl))
		return nextURL, nil
	}

	return "", nil
}

func (zc *ZillowCrawler) CrawlZillowSearchResultData(result schemas.ZillowSearchPageResResult) *entity.Zillow {

	propertyStatus := false
	if result.Beds > 0 || result.Baths > 0 {
		propertyStatus = true
	}
	halfBathRooms := result.HdpData.HomeInfo.Bedrooms / 2
	fullBathRooms := result.HdpData.HomeInfo.Bathrooms - halfBathRooms
	return &entity.Zillow{
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

}

func (zc *ZillowCrawler) CrawlZillowSearchRelaxedData(result schemas.ZillowSearchPageResRelaxedResult) *entity.Zillow {

	propertyStatus := false
	if result.Beds > 0 || result.Baths > 0 {
		propertyStatus = true
	}
	halfBathRooms := result.HdpData.HomeInfo.Bedrooms / 2
	fullBathRooms := result.HdpData.HomeInfo.Bathrooms - halfBathRooms
	return &entity.Zillow{
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

}

func (zc *ZillowCrawler) CrawlAddress(address string) error {
	if err := zc.WebDriver.Get(address); err != nil {
		return err
	}
	// NOTE: time to load source. Need to increase if data was not showing
	//time.Sleep(viper.GetDuration("crawler.zillow_crawler.time_load_source") * time.Second)

	if err := zc.WebDriver.WaitWithTimeoutAndInterval(func(wd selenium.WebDriver) (bool, error) {
		viewContainer, err := wd.FindElement(selenium.ByXPATH, "//div[@class=\"data-view-container\"]")

		if err != nil {
			return false, err
		}

		display, err := viewContainer.IsDisplayed()

		if err != nil {
			return false, err
		}

		return display, nil
	}, 1000, 1000); err != nil {
		return err
	}

	pageSource, err := zc.WebDriver.PageSource()
	if err != nil {
		zc.BrowserTurnOff = true
		return err
	}
	if err := zc.ByPassVerifyHuman(pageSource, address); err != nil {
		return err
	}
	if err := zc.ParseData(pageSource); err != nil {
		return err
	}
	if zc.OnlyHistoryTable == false {
		if err := zc.UpdateZillowDB(); err != nil {
			return err
		}
		zc.ShowLogInfo("Added/Updated record to Zillow Table")
	}

	if err := zc.UpdateZillowPriceHistoryDB(); err != nil {
		return err
	}
	zc.ShowLogInfo("Added/Updated record to Price History Table")
	if err := zc.UpdateZillowPublicTaxHistoryDB(); err != nil {
		return err
	}
	zc.ShowLogInfo("Added/Updated record to Public Taxt History Table")
	if err := zc.UpdateMaindb3DB(); err != nil {
		return err
	}

	// Begin to Crawl map data
	//zc.CrawlSearchData()
	//if zc.SearchDataByCollyStatus == false {
	//	if err := zc.CrawlSearchData(); err != nil {
	//		return err
	//	}
	//}
	//if len(zc.CrawlerSchemas.ZillowSearchData) > 0 {
	//	for _, v := range zc.CrawlerSchemas.ZillowSearchData {
	//		zillowData, err := zc.CrawlerServices.ZillowService.GetZillowByURL(v.URL)
	//		if err != nil {
	//			return err
	//		}
	//		if zillowData == nil {
	//			v.CrawlingStatus = viper.GetString("crawler.crawler_status.rerun")
	//			if err := zc.CrawlerServices.ZillowService.AddZillow(v); err != nil {
	//				return err
	//			}
	//			zc.ShowLogInfo(fmt.Sprint("Added Searching Data to Zillow Detail Table with ID: ", v.ID))
	//		}
	//	}
	//}
	return nil
}

func (zc *ZillowCrawler) CrawlAddressAPI(address string) error {
	if err := zc.WebDriver.Get(address); err != nil {
		return err
	}
	//// NOTE: time to load source. Need to increase if data was not showing
	//timeLoad := viper.GetDuration("crawler.zillow_crawler.time_load_source") * time.Second
	//if err := zc.WebDriver.WaitWithTimeoutAndInterval(func(wd selenium.WebDriver) (bool, error) {
	//	viewContainer, _ := wd.FindElement(selenium.ByXPATH, "//div[@class=\"data-view-container\"]//*[contains(text(), \"Lot size\")]")
	//	if viewContainer != nil {
	//		display, _ := viewContainer.IsDisplayed()
	//		return display, nil
	//	}
	//	return false, nil
	//}, timeLoad, timeLoad); err != nil {
	//	return err
	//}
	time.Sleep(viper.GetDuration("crawler.zillow_crawler.time_load_source") * time.Second)

	pageSource, err := zc.WebDriver.PageSource()
	if err != nil {
		zc.BrowserTurnOff = true
		return err
	}
	if err := zc.ByPassVerifyHuman(pageSource, address); err != nil {
		return err
	}
	if err := zc.ParseData(pageSource); err != nil {
		return err
	}
	return nil
}

func (zc *ZillowCrawler) UpdateZillowDB() error {
	zillowData, err := zc.CrawlerServices.ZillowService.GetZillowByID(zc.CrawlerTables.Maindb3.ID)
	if err != nil {
		return err
	}
	if zillowData == nil {
		if err := zc.CrawlerServices.ZillowService.AddZillow(zc.CrawlerTables.ZillowData); err != nil {
			return err
		}
	} else {
		if err := zc.CrawlerServices.ZillowService.UpdateZillow(zc.CrawlerTables.ZillowData, zillowData.ID); err != nil {
			return err
		}
	}
	return nil
}
func (zc *ZillowCrawler) UpdateMaindb3DB() error {
	if err := zc.CrawlerServices.Maindb3Service.UpdateStatus(zc.CrawlerTables.Maindb3, viper.GetString("crawler.crawler_status.succeeded")); err != nil {
		return err
	}
	return nil
}

func (zc *ZillowCrawler) UpdateZillowPriceHistoryDB() error {
	if err := zc.CrawlerServices.ZillowService.UpdateZillowPriceHistory(zc.CrawlerTables.ZillowPriceHistory); err != nil {
		return err
	}
	return nil
}

func (zc *ZillowCrawler) UpdateZillowPublicTaxHistoryDB() error {
	if err := zc.CrawlerServices.ZillowService.UpdateZillowPublicTaxHistory(zc.CrawlerTables.ZillowPublicTaxHistory); err != nil {
		return err
	}
	return nil
}

func (zc *ZillowCrawler) ParseData(source string) error {
	//htmlquery.DisableSelectorCache = true
	doc, err := htmlquery.Parse(strings.NewReader(source))

	// Need to sure the data is existing
	detailPageContainer := htmlquery.FindOne(doc, "//div[@id=\"details-page-container\"]")
	if detailPageContainer == nil {
		return fmt.Errorf("not found data from address requested")
	}
	if err != nil {
		return err
	}
	zc.ParseURL()
	zc.ParseBedBathSF(doc)
	zc.ParseAddress(doc)
	if zc.OnlyHistoryTable == false {
		zc.ParseFullBathroom(doc)
		zc.ParseHalfBathroom(doc)
		zc.ParseSalePrice(doc)
		zc.ParseZestimate(doc)
		zc.ParseEstPayment(doc)
		zc.ParsePrincipalInterest(doc)
		zc.ParseMortgageInsurance(doc)
		zc.ParsePropertyTaxes(doc)
		zc.ParseHomeInsurance(doc)
		zc.ParseHoaFee(doc)
		zc.ParseUtilities(doc)
		zc.ParseEstimatedSalesRange(doc)
		zc.ParsePictures(doc)
		zc.ParseTimeOnZillow(doc)
		zc.ParseViews(doc)
		zc.ParseSaves(doc)
		zc.ParseOverview(doc)
		zc.ParseMSL(doc)
		zc.ParseZillowCheckedDate(doc)
		zc.ParseDataUploadedDate(doc)
		zc.ParseListBy(doc)
		zc.ParseSourceZillow(doc)
		zc.ParseYearBuilt(doc)
		zc.ParseNaturalGas(doc)
		zc.ParseCentralAir(doc)
		zc.ParseGarageSpaces(doc)
		zc.ParseHoaAmount(doc)
		zc.ParseLotSizes(doc)
		zc.ParseBuyerAgentFee(doc)
		zc.ParseApplicances(doc)
		zc.ParseLivingRooms(doc)
		zc.ParsePrimaryBedRooms(doc)
		zc.ParseInteriorFeatures(doc)
		zc.ParseBasement(doc)
		zc.ParseTotalInteriorLivableArea(doc)
		zc.ParseOffFireplaces(doc)
		zc.ParseFireplaceFeatures(doc)
		zc.ParseFlooringType(doc)
		zc.ParseHeatingType(doc)
		zc.ParseParking(doc)
		zc.ParseLotFeatures(doc)
		zc.ParseParcelNumber(doc)
		zc.ParsePropertydetails(doc)
		zc.ParseConstructionDetails(doc)
		zc.ParseUtiGreenEnergyDetails(doc)
		zc.ParseComNeiDetails(doc)
		zc.ParseHoaFinancialDetails(doc)
		zc.ParseGreatSchools(doc)
		zc.ParseDistrict(doc)
		zc.ParseDataSource(doc)
	}
	zc.ParseMapBounds(doc)
	zc.ParseZillowPriceHistory(doc)
	zc.ParseZillowPublicTaxHistory(doc)

	return nil
}

func (zc *ZillowCrawler) ParseURL() {
	if zc.CrawlerTables.ZillowData.URL == "" {
		url, err := zc.BaseSel.WebDriver.CurrentURL()
		if err != nil {
			zc.ShowLogError(err.Error())
		}
		zc.CrawlerTables.ZillowData.URL = url
	}
}

func (zc *ZillowCrawler) ParseBedBathSF(doc *html.Node) {
	if zc.CrawlerTables.ZillowData.Bed == 0 || zc.CrawlerTables.ZillowData.Bath == 0 {
		bedPathItems := htmlquery.Find(doc, "//span[contains(@data-testid,\"bed-bath\")]/span | //span[contains(@data-testid,\"bed-bath\")]/button")
		for _, item := range bedPathItems {
			itemText := htmlquery.InnerText(item)
			if strings.Contains(itemText, "bd") {
				bedStr := strings.Replace(itemText, "bd", "", -1)
				bedStr = util2.RemoveSpecialCharacters(bedStr)
				bedStr = strings.TrimSpace(bedStr)
				if bedStr != "" {
					bed, err := strconv.Atoi(bedStr)
					if err != nil {
						zc.ShowLogError(err.Error())
					} else {
						zc.CrawlerTables.ZillowData.Bed = float64(bed)
					}

				}

			}
			if strings.Contains(itemText, "ba") {
				bathStr := strings.Replace(itemText, "ba", "", -1)
				bathStr = util2.RemoveSpecialCharacters(bathStr)
				bathStr = strings.TrimSpace(bathStr)
				if bathStr != "" {
					bath, err := strconv.Atoi(bathStr)
					if err != nil {
						zc.ShowLogError(err.Error())
					} else {
						zc.CrawlerTables.ZillowData.Bath = float64(bath)
					}
				}

			}
			if strings.Contains(itemText, "sqft") {
				sfStr := strings.Replace(itemText, "sqft", "", -1)
				sfStr = strings.Replace(sfStr, ",", ".", -1)
				sfStr = strings.Replace(sfStr, "-", "", -1)
				sfStr = util2.RemoveSpecialCharacters(sfStr)
				sfStr = strings.TrimSpace(sfStr)
				if sfStr != "" {
					if sf, err := strconv.ParseFloat(sfStr, 64); err != nil {
						zc.ShowLogError(err.Error())
					} else {
						zc.CrawlerTables.ZillowData.SF = sf
					}
				}
			}
		}
	}

	// Property Status
	if zc.CrawlerTables.ZillowData.Bed > 0 || zc.CrawlerTables.ZillowData.Bath > 0 {
		zc.CrawlerTables.ZillowData.PropertyStatus = true
	}
}

func (zc *ZillowCrawler) ParseAddress(doc *html.Node) {
	if zc.CrawlerTables.ZillowData.Address == "" {
		addresses := htmlquery.Find(doc, "//h1/text()")
		for _, v := range addresses {
			zc.CrawlerTables.ZillowData.Address += v.Data
		}
	}

}

func (zc *ZillowCrawler) ParseFullBathroom(doc *html.Node) {
	// Full Bathrooms
	if zc.CrawlerTables.ZillowData.FullBathrooms == 0 {
		fullPathRoom := htmlquery.FindOne(doc, "//span[contains(text(), \"Full bathrooms\")]")
		if fullPathRoom != nil {
			fullBathRoomText := htmlquery.InnerText(fullPathRoom)
			fullBathRoomText = strings.Replace(fullBathRoomText, "Full bathrooms", "", -1)
			fullBathRoomText = strings.Replace(fullBathRoomText, ":", "", -1)
			fullBathRoomText = util2.RemoveSpecialCharacters(fullBathRoomText)
			fullBathRoomText = strings.TrimSpace(fullBathRoomText)
			if fullBathRoomText != "" {
				fullBathRoomValue, err := strconv.Atoi(strings.TrimSpace(fullBathRoomText))
				if err != nil {
					zc.ShowLogError(err.Error())
				} else {
					zc.CrawlerTables.ZillowData.FullBathrooms = float64(fullBathRoomValue)
				}
			}
		}
	}
}

func (zc *ZillowCrawler) ParseHalfBathroom(doc *html.Node) {
	if zc.CrawlerTables.ZillowData.HalfBathrooms == 0 {
		halfPathRoom := htmlquery.FindOne(doc, "//h6[contains(text(), \"Bedrooms and bathrooms\")]/following-sibling::ul//span[contains(text(), \"Bathrooms\")]")
		if halfPathRoom != nil {
			halfPathRoomText := htmlquery.InnerText(halfPathRoom)
			halfPathRoomText = strings.Replace(halfPathRoomText, "Bathrooms", "", -1)
			halfPathRoomText = strings.Replace(halfPathRoomText, ":", "", -1)
			halfPathRoomText = util2.RemoveSpecialCharacters(halfPathRoomText)
			halfPathRoomText = strings.TrimSpace(halfPathRoomText)
			if halfPathRoomText != "" {
				halfBathRoomValue, err := strconv.Atoi(halfPathRoomText)
				if err != nil {
					zc.ShowLogError(err.Error())
				} else {
					zc.CrawlerTables.ZillowData.HalfBathrooms = float64(halfBathRoomValue)
				}

			}

		}
	}

}

func (zc *ZillowCrawler) ParseSalePrice(doc *html.Node) {
	if zc.CrawlerTables.ZillowData.SalesPrice == 0 {
		salePrice := htmlquery.FindOne(doc, "//span[@data-testid=\"price\"]/span/text()")
		if salePrice == nil {
			salePrice = htmlquery.FindOne(doc, "//*[contains(text(), \"Estimated sale price\")]/following-sibling::p/text()")
		}
		if salePrice != nil {
			salePriceStr := strings.Replace(salePrice.Data, "$", "", -1)
			salePriceStr = strings.Replace(salePriceStr, ",", ".", -1)
			salePriceStr = util2.RemoveSpecialCharacters(salePriceStr)
			salePriceStr = strings.TrimSpace(salePriceStr)
			if salePriceStr != "" {
				salePrice, err := strconv.ParseFloat(salePriceStr, 64)
				if err != nil {
					zc.ShowLogError(err.Error())
				} else {
					zc.CrawlerTables.ZillowData.SalesPrice = salePrice
				}
			}
		}

	}

}

func (zc *ZillowCrawler) ParseZestimate(doc *html.Node) {
	if zc.CrawlerTables.ZillowData.RentZestimate == 0 || zc.CrawlerTables.ZillowData.Zestimate == 0 {
		zestimates := htmlquery.Find(doc, "//*[contains(text(), \"Zestimate\")]/following-sibling::span/span/text()")
		if zestimates != nil {
			zestimateStr := strings.Replace(zestimates[0].Data, "$", "", -1)
			zestimateStr = strings.Replace(zestimateStr, ",", ".", -1)
			zestimateStr = util2.RemoveSpecialCharacters(zestimateStr)
			zestimateStr = strings.TrimSpace(zestimateStr)
			if zestimateStr != "" {
				if zestimateValue, err := strconv.ParseFloat(zestimateStr, 64); err != nil {
					zc.ShowLogError(err.Error())
				} else {
					zc.CrawlerTables.ZillowData.Zestimate = zestimateValue
				}
			}

			if len(zestimates) > 1 {
				rentZestimateStr := strings.Replace(zestimates[1].Data, "$", "", -1)
				rentZestimateStr = strings.Replace(rentZestimateStr, ",", ".", -1)
				rentZestimateStr = util2.RemoveSpecialCharacters(rentZestimateStr)
				rentZestimateStr = strings.TrimSpace(rentZestimateStr)
				if rentZestimateStr != "" {
					if rentZestimateValue, err := strconv.ParseFloat(rentZestimateStr, 64); err != nil {
						zc.ShowLogError(err.Error())
					} else {
						zc.CrawlerTables.ZillowData.RentZestimate = rentZestimateValue
					}
				}

			}

		}
	}

}

func (zc *ZillowCrawler) ParseEstPayment(doc *html.Node) {
	estPayment := htmlquery.FindOne(doc, "//div[@class='summary-container']//span[contains(text(), 'Est. payment')]/following-sibling::span/text()")
	if estPayment != nil {
		zc.CrawlerTables.ZillowData.EstPayment = strings.TrimSpace(estPayment.Data)
	}

}

func (zc *ZillowCrawler) ParsePrincipalInterest(doc *html.Node) {
	principalInterest := htmlquery.FindOne(doc, "//h5[normalize-space(text())='Principal & interest']/following-sibling::span/text()")
	if principalInterest != nil {
		zc.CrawlerTables.ZillowData.PrincipalInterest = strings.TrimSpace(principalInterest.Data)
	}
}

func (zc *ZillowCrawler) ParseMortgageInsurance(doc *html.Node) {
	mortgageInsurance := htmlquery.FindOne(doc, "//h5[normalize-space(text())='Mortgage insurance']/following-sibling::span/text()")
	if mortgageInsurance != nil {
		zc.CrawlerTables.ZillowData.MortgageInsurance = strings.TrimSpace(mortgageInsurance.Data)
	}

}

func (zc *ZillowCrawler) ParsePropertyTaxes(doc *html.Node) {
	propertyTaxes := htmlquery.FindOne(doc, "//h5[normalize-space(text())='Property taxes']/following-sibling::span/text()")
	if propertyTaxes != nil {
		zc.CrawlerTables.ZillowData.PropertyTaxes = strings.TrimSpace(propertyTaxes.Data)
	}

}

func (zc *ZillowCrawler) ParseHomeInsurance(doc *html.Node) {
	homeInsurance := htmlquery.FindOne(doc, "//h5[contains(text(), 'Home insurance')]/following-sibling::span/text()")
	if homeInsurance != nil {
		zc.CrawlerTables.ZillowData.HomeInsurance = strings.TrimSpace(homeInsurance.Data)
	}
}

func (zc *ZillowCrawler) ParseHoaFee(doc *html.Node) {
	hoaFees := htmlquery.FindOne(doc, "//h5[contains(text(), 'HOA fee')]/following-sibling::span/text()")
	if hoaFees != nil {
		zc.CrawlerTables.ZillowData.HOAFee = strings.TrimSpace(hoaFees.Data)
	}
}

func (zc *ZillowCrawler) ParseUtilities(doc *html.Node) {
	utilities := htmlquery.FindOne(doc, "//h5[contains(text(), \"Utilities\")]/following-sibling::span/text()")
	if utilities != nil {
		zc.CrawlerTables.ZillowData.Utilities = strings.TrimSpace(utilities.Data)
	}

}

func (zc *ZillowCrawler) ParseEstimatedSalesRange(doc *html.Node) {
	estimatedSalesRange := htmlquery.FindOne(doc, "//span[contains(text(), 'Estimated sales range')]/span/text()")
	if estimatedSalesRange != nil {
		estimatedSalesRangeList := strings.Split(strings.TrimSpace(estimatedSalesRange.Data), "-")
		zc.CrawlerTables.ZillowData.EstimatedSalesRangeMinimum = strings.TrimSpace(estimatedSalesRangeList[0])
		zc.CrawlerTables.ZillowData.EstimatedSalesRangeMax = strings.TrimSpace(estimatedSalesRangeList[1])
	}

}

func (zc *ZillowCrawler) ParsePictures(doc *html.Node) {
	pictures := htmlquery.Find(doc, "//*[contains(@class, \"media-stream-tile\")]//img")

	if pictures != nil {
		var picSlice []string
		for _, pic := range pictures {
			picSlice = append(picSlice, htmlquery.SelectAttr(pic, "src"))
		}
		zc.CrawlerTables.ZillowData.Pictures = strings.Join(picSlice, ", ")
	}

}

func (zc *ZillowCrawler) ParseTimeOnZillow(doc *html.Node) {
	timeOnZillow := htmlquery.FindOne(doc, "//dt[contains(text(), \"Time on Zillow\")]/following-sibling::dd/strong/text()")
	if timeOnZillow != nil {
		zc.CrawlerTables.ZillowData.TimeOnZillow = strings.TrimSpace(timeOnZillow.Data)
	}

}

func (zc *ZillowCrawler) ParseViews(doc *html.Node) {
	views := htmlquery.FindOne(doc, "//dt/button[contains(text(), \"Views\") ]/parent::dt/following-sibling::dd/strong/text()")
	if views != nil {
		viewsData := util2.RemoveSpecialCharacters(views.Data)
		viewsData = strings.TrimSpace(viewsData)
		if viewsData != "" {
			if viewsValue, err := strconv.Atoi(viewsData); err != nil {
				zc.ShowLogError(err.Error())
			} else {
				zc.CrawlerTables.ZillowData.Views = viewsValue
			}
		}

	}
}

func (zc *ZillowCrawler) ParseSaves(doc *html.Node) {
	saves := htmlquery.FindOne(doc, "//dt/button[contains(text(), \"Saves\") ]/parent::dt/following-sibling::dd/strong/text()")
	if saves != nil {
		savesData := util2.RemoveSpecialCharacters(saves.Data)
		savesData = strings.TrimSpace(savesData)
		if savesData != "" {
			if savesValue, err := strconv.Atoi(savesData); err != nil {
				zc.ShowLogError(err.Error())
			} else {
				zc.CrawlerTables.ZillowData.Saves = savesValue
			}
		}

	}

}

func (zc *ZillowCrawler) ParseOverview(doc *html.Node) {
	overview := htmlquery.FindOne(doc, "//h4[contains(text(), \"Overview\")]/following-sibling::div//div[contains(@class, \"Spacer\")]//div[contains(@class, \"Text\")]/text()")
	if overview != nil {
		zc.CrawlerTables.ZillowData.Overview = strings.TrimSpace(overview.Data)
	}
}

func (zc *ZillowCrawler) ParseMSL(doc *html.Node) {
	mls := htmlquery.FindOne(doc, "//span[contains(text(), \"MLS#:\")]/text()")
	if mls != nil {
		zc.CrawlerTables.ZillowData.MLS = strings.TrimSpace(strings.Replace(mls.Data, "MLS#:", "", -1))
	}

}

func (zc *ZillowCrawler) ParseZillowCheckedDate(doc *html.Node) {
	zillowCheckedDate := htmlquery.FindOne(doc, "//*[contains(text(), \"Zillow checked:\")]/text()")
	if zillowCheckedDate != nil {
		zc.CrawlerTables.ZillowData.ZillowCheckedDate = strings.TrimSpace(strings.Replace(zillowCheckedDate.Data, "Zillow checked:", "", -1))
	}

}

func (zc *ZillowCrawler) ParseDataUploadedDate(doc *html.Node) {
	dataUploadedDate := htmlquery.FindOne(doc, "//*[contains(text(), \"Data updated:\")]/text()")
	if dataUploadedDate != nil {
		zc.CrawlerTables.ZillowData.DataUploadedDate = strings.TrimSpace(strings.Replace(dataUploadedDate.Data, "Data updated:", "", -1))
	}
}

func (zc *ZillowCrawler) ParseListBy(doc *html.Node) {
	listBy := htmlquery.Find(doc, "//*[contains(text(), \"Listed by:\")]/following-sibling::span/p/text()")
	if listBy != nil {
		var listBySlice []string
		for _, listByValue := range listBy {
			if listByValue.Data != "" {
				listBySlice = append(listBySlice, listByValue.Data)
			}
		}
		zc.CrawlerTables.ZillowData.ListedBy = strings.Join(listBySlice, "| ")
	}
}

func (zc *ZillowCrawler) ParseSourceZillow(doc *html.Node) {
	sourceZillow := htmlquery.FindOne(doc, "//*[contains(text(), \"Source:\")]/text()")
	if sourceZillow != nil {
		zc.CrawlerTables.ZillowData.Source = strings.TrimSpace(strings.Replace(sourceZillow.Data, "Source:", "", -1))
	}
}

func (zc *ZillowCrawler) ParseYearBuilt(doc *html.Node) {
	yearBuilt := htmlquery.FindOne(doc, "//span[contains(text(), \"Year built\")]/text()")
	if yearBuilt != nil {
		zc.CrawlerTables.ZillowData.YearBuilt = strings.TrimSpace(strings.Replace(yearBuilt.Data, "Year built:", "", -1))
	}

}

func (zc *ZillowCrawler) ParseNaturalGas(doc *html.Node) {
	naturalGas := htmlquery.FindOne(doc, "//*[contains(text(), \"Natural Gas\") or contains(text(), \"natural gas\")]")
	if naturalGas != nil {
		zc.CrawlerTables.ZillowData.NaturalGas = true
	}

}

func (zc *ZillowCrawler) ParseCentralAir(doc *html.Node) {
	centralAir := htmlquery.FindOne(doc, "//*[contains(text(), \"Central Air\") or contains(text(), \"central air\")]")
	if centralAir != nil {
		zc.CrawlerTables.ZillowData.CentralAir = true
	}
}

func (zc *ZillowCrawler) ParseGarageSpaces(doc *html.Node) {
	garageSpaces := htmlquery.FindOne(doc, "//h4[contains(text(), \"Facts and features\")]/following-sibling::div//*[contains(text(), \"garage spaces\") or contains(text(), \"Garage spaces\") or contains(text(), \"Garage Spaces\")]/text()")
	if garageSpaces != nil {
		zc.CrawlerTables.ZillowData.OfGarageSpaces = strings.TrimSpace(strings.Replace(garageSpaces.Data, " garage spaces", "", -1))
	}
}

func (zc *ZillowCrawler) ParseHoaAmount(doc *html.Node) {
	hoaAmount := htmlquery.FindOne(doc, "//*[contains(text(), \"annually HOA fee\")]/text()")
	if hoaAmount != nil {
		zc.CrawlerTables.ZillowData.HOAAmount = strings.TrimSpace(strings.Replace(hoaAmount.Data, " annually HOA fee", "", -1))
	}
}

func (zc *ZillowCrawler) ParseLotSizes(doc *html.Node) {
	lotSizesDoc := htmlquery.FindOne(doc, "//div[@class=\"data-view-container\"]//*[contains(text(), \"Lot size\")]")
	if lotSizesDoc != nil {
		lotSizesText := htmlquery.InnerText(lotSizesDoc)
		lotSizesText = strings.TrimSpace(strings.Replace(lotSizesText, "Lot size:", "", -1))
		if strings.Contains(lotSizesText, "sqft") == true {
			zc.CrawlerTables.ZillowData.LotSizeSF = lotSizesText
		}
		if strings.Contains(lotSizesText, "Acres") == true {
			zc.CrawlerTables.ZillowData.LotSizeAcres = lotSizesText
		}

	}
}

func (zc *ZillowCrawler) ParseBuyerAgentFee(doc *html.Node) {
	buyerAgentFee := htmlquery.FindOne(doc, "//*[contains(text(), \"buyer's agent fee\")]/text()")
	if buyerAgentFee != nil {
		zc.CrawlerTables.ZillowData.BuyerAgentFee = strings.TrimSpace(strings.Replace(buyerAgentFee.Data, " buyer's agent fee", "", -1))
	}

}

func (zc *ZillowCrawler) ParseApplicances(doc *html.Node) {
	applicances := htmlquery.FindOne(doc, "//*[contains(text(), \"Appliances included\")]")
	if applicances != nil {
		zc.CrawlerTables.ZillowData.Appliances = strings.TrimSpace(strings.Replace(htmlquery.InnerText(applicances), "Appliances included:", "", -1))
	}
}

func (zc *ZillowCrawler) ParseLivingRooms(doc *html.Node) {
	// Living Room
	livingRooms := htmlquery.Find(doc, "//h6[contains(text(), \"LivingRoom\")]/following-sibling::ul")
	for _, livingroom := range livingRooms {
		// Living Room Level
		livingRoomLevel := htmlquery.FindOne(livingroom, ".//span[contains(text(), \"Level\")]")
		if livingRoomLevel != nil {
			zc.CrawlerTables.ZillowData.LivingRoomLevel = strings.TrimSpace(strings.Replace(htmlquery.InnerText(livingRoomLevel), "Level:", "", -1))
		}

		// Living Room Dimensions
		livingRoomDimensions := htmlquery.FindOne(livingroom, ".//span[contains(text(), \"Dimensions\")]")
		if livingRoomDimensions != nil {
			zc.CrawlerTables.ZillowData.LivingRoomDimensions = strings.TrimSpace(strings.Replace(htmlquery.InnerText(livingRoomDimensions), "Dimensions:", "", -1))
		}
	}

}

func (zc *ZillowCrawler) ParsePrimaryBedRooms(doc *html.Node) {
	// Primary Bedroom
	primaryBedRooms := htmlquery.Find(doc, "//h6[contains(text(), \"PrimaryBedroom\")]/following-sibling::ul")
	for _, primaryBedRoom := range primaryBedRooms {
		// Primary Bedroom Level
		primaryBedRoomLevel := htmlquery.FindOne(primaryBedRoom, ".//span[contains(text(), \"Level\")]")
		if primaryBedRoomLevel != nil {
			zc.CrawlerTables.ZillowData.PrimaryBedroomLevel = strings.TrimSpace(strings.Replace(htmlquery.InnerText(primaryBedRoomLevel), "Level:", "", -1))
		}

		// Primary Bedroom Dimensions
		primaryBedRoomDimensions := htmlquery.FindOne(primaryBedRoom, ".//span[contains(text(), \"Dimensions\")]")
		if primaryBedRoomDimensions != nil {
			zc.CrawlerTables.ZillowData.PrimaryBedroomDimensions = strings.TrimSpace(strings.Replace(htmlquery.InnerText(primaryBedRoomDimensions), "Dimensions:", "", -1))
		}
	}

}

func (zc *ZillowCrawler) ParseInteriorFeatures(doc *html.Node) {
	// Interior Features
	interiorFeatures := htmlquery.FindOne(doc, "//span[contains(text(), \"Interior features\")]")
	if interiorFeatures != nil {
		zc.CrawlerTables.ZillowData.InteriorFeatures = strings.TrimSpace(strings.Replace(htmlquery.InnerText(interiorFeatures), "Interior features:", "", -1))
	}

}

func (zc *ZillowCrawler) ParseBasement(doc *html.Node) {
	basement := htmlquery.FindOne(doc, "//span[contains(text(), \"Basement\")]")
	if basement != nil {
		zc.CrawlerTables.ZillowData.Basement = strings.TrimSpace(strings.Replace(htmlquery.InnerText(basement), "Basement:", "", -1))
	}
}

func (zc *ZillowCrawler) ParseTotalInteriorLivableArea(doc *html.Node) {
	// Total Interior Livable Area SF
	totalInteriorLivableArea := htmlquery.FindOne(doc, "//span[contains(text(), \"Total interior livable area\")]")
	if totalInteriorLivableArea != nil {
		zc.CrawlerTables.ZillowData.TotalInteriorLivableAreaSF = strings.TrimSpace(strings.Replace(htmlquery.InnerText(totalInteriorLivableArea), "Total interior livable area:", "", -1))
	}

}

func (zc *ZillowCrawler) ParseOffFireplaces(doc *html.Node) {
	// # of Fireplaces
	offFireplaces := htmlquery.FindOne(doc, "//span[contains(text(), \"Total number of fireplaces\")]")
	if offFireplaces != nil {
		zc.CrawlerTables.ZillowData.OfFireplaces = strings.TrimSpace(strings.Replace(htmlquery.InnerText(offFireplaces), "Total number of fireplaces:", "", -1))
	}

}

func (zc *ZillowCrawler) ParseFireplaceFeatures(doc *html.Node) {
	// Fireplace features
	fireplaceFeatures := htmlquery.FindOne(doc, "//span[contains(text(), \"Fireplace features\")]")
	if fireplaceFeatures != nil {
		zc.CrawlerTables.ZillowData.FireplaceFeatures = strings.TrimSpace(strings.Replace(htmlquery.InnerText(fireplaceFeatures), "Fireplace features:", "", -1))
	}
}

func (zc *ZillowCrawler) ParseFlooringType(doc *html.Node) {
	// Flooring Type
	flooringType := htmlquery.FindOne(doc, "//h6[contains(text(), \"Flooring\")]/following-sibling::ul//span[contains(text(), \"Flooring\")]")
	if flooringType != nil {
		zc.CrawlerTables.ZillowData.FlooringType = strings.TrimSpace(strings.Replace(htmlquery.InnerText(flooringType), "Flooring:", "", -1))
	}

}

func (zc *ZillowCrawler) ParseHeatingType(doc *html.Node) {
	// Heating Type
	heatingType := htmlquery.FindOne(doc, "//h6[contains(text(), \"Heating\")]/following-sibling::ul//span[contains(text(), \"Heating features\")]")
	if heatingType != nil {
		zc.CrawlerTables.ZillowData.HeatingType = strings.TrimSpace(strings.Replace(htmlquery.InnerText(heatingType), "Heating features:", "", -1))
	}

}

func (zc *ZillowCrawler) ParseParking(doc *html.Node) {
	// Parking
	parkings := htmlquery.Find(doc, "//h6[contains(text(), \"Parking\")]/following-sibling::ul")
	if parkings != nil {
		for _, parking := range parkings {
			// Total Parking Spaces
			totalParkingSpaces := htmlquery.FindOne(parking, ".//span[contains(text(), \"Total spaces\")]")
			if totalParkingSpaces != nil {
				zc.CrawlerTables.ZillowData.TotalParkingSpaces = strings.TrimSpace(strings.Replace(htmlquery.InnerText(totalParkingSpaces), "Total spaces:", "", -1))
			}

			// Parking Features
			parkingFeatures := htmlquery.FindOne(parking, ".//span[contains(text(), \"Parking features\")]")
			if parkingFeatures != nil {
				zc.CrawlerTables.ZillowData.ParkingFeatures = strings.TrimSpace(strings.Replace(htmlquery.InnerText(parkingFeatures), "Parking features:", "", -1))
			}

			// Covered Spaces
			coveredSpaces := htmlquery.FindOne(parking, ".//span[contains(text(), \"Covered spaces\")]")
			if coveredSpaces != nil {
				zc.CrawlerTables.ZillowData.CoveredSpaces = strings.TrimSpace(strings.Replace(htmlquery.InnerText(coveredSpaces), "Covered spaces:", "", -1))
			}

		}
	}

}

func (zc *ZillowCrawler) ParseLotFeatures(doc *html.Node) {
	// Lot Features
	lotFeatures := htmlquery.FindOne(doc, "//h6[contains(text(), \"Lot\")]/following-sibling::ul//span[contains(text(), \"Lot features\")]")
	if lotFeatures != nil {
		zc.CrawlerTables.ZillowData.LotFeatures = strings.TrimSpace(strings.Replace(htmlquery.InnerText(lotFeatures), "Lot features:", "", -1))
	}

}

func (zc *ZillowCrawler) ParseParcelNumber(doc *html.Node) {
	// Parcel number
	parcelNumber := htmlquery.FindOne(doc, "//h6[contains(text(), \"Other property information\")]/following-sibling::ul//span[contains(text(), \"Parcel number\")]")
	if parcelNumber != nil {
		zc.CrawlerTables.ZillowData.ParcelNumber = strings.TrimSpace(strings.Replace(htmlquery.InnerText(parcelNumber), "Parcel number:", "", -1))
	}

}

func (zc *ZillowCrawler) ParsePropertydetails(doc *html.Node) {
	// Property details - Property
	propertydetails := htmlquery.Find(doc, "//h5[contains(text(), \"Property details\")]/following-sibling::div//h6[contains(text(), \"Property\")]/following-sibling::ul")
	if propertydetails != nil {
		for _, property := range propertydetails {
			// # Levels (Stories/Floors)
			levelsStoriesFloors := htmlquery.FindOne(property, ".//span[contains(text(), \"Levels\")]")
			if levelsStoriesFloors != nil {
				zc.CrawlerTables.ZillowData.LevelsStoriesFloors = strings.TrimSpace(strings.Replace(htmlquery.InnerText(levelsStoriesFloors), "Levels:", "", -1))
			}

			// Patio and Porch Details
			patioAndPorchDetails := htmlquery.FindOne(property, ".//span[contains(text(), \"Patio and porch details\")]")
			if patioAndPorchDetails != nil {
				zc.CrawlerTables.ZillowData.PatioAndPorchDetails = strings.TrimSpace(strings.Replace(htmlquery.InnerText(patioAndPorchDetails), "Patio and porch details:", "", -1))
			}

		}
	}

}

func (zc *ZillowCrawler) ParseConstructionDetails(doc *html.Node) {
	// Construction details
	constructionDetails := htmlquery.Find(doc, "//h5[contains(text(), \"Construction details\")]/following-sibling::div//h6/following-sibling::ul")
	if constructionDetails != nil {
		for _, constructionDetail := range constructionDetails {
			// HomeType
			homeType := htmlquery.FindOne(constructionDetail, ".//span[contains(text(), \"Home type\")]")
			if homeType != nil {
				zc.CrawlerTables.ZillowData.HomeType = strings.TrimSpace(strings.Replace(htmlquery.InnerText(homeType), "Home type:", "", -1))
			}
			// Propery SubType
			propertySubType := htmlquery.FindOne(constructionDetail, ".//span[contains(text(), \"Property subType\")]")
			if propertySubType != nil {
				zc.CrawlerTables.ZillowData.ProperySubType = strings.TrimSpace(strings.Replace(htmlquery.InnerText(propertySubType), "Property subType:", "", -1))
			}

			// Construction Materials
			constructionMaterials := htmlquery.FindOne(constructionDetail, ".//span[contains(text(), \"Construction materials\")]")
			if constructionMaterials != nil {
				zc.CrawlerTables.ZillowData.ConstructionMaterials = strings.TrimSpace(strings.Replace(htmlquery.InnerText(constructionMaterials), "Construction materials:", "", -1))
			}

			// Foundation
			foundation := htmlquery.FindOne(constructionDetail, ".//span[contains(text(), \"Foundation\")]")
			if foundation != nil {
				zc.CrawlerTables.ZillowData.Foundation = strings.TrimSpace(strings.Replace(htmlquery.InnerText(foundation), "Foundation:", "", -1))
			}

			// Roof
			roof := htmlquery.FindOne(constructionDetail, ".//span[contains(text(), \"Roof\")]")
			if roof != nil {
				zc.CrawlerTables.ZillowData.Roof = strings.TrimSpace(strings.Replace(htmlquery.InnerText(roof), "Roof:", "", -1))
			}

			// New Construction
			newConstruction := htmlquery.FindOne(constructionDetail, ".//span[contains(text(), \"New construction\")]")
			if newConstruction != nil {
				zc.CrawlerTables.ZillowData.NewConstruction = strings.TrimSpace(strings.Replace(htmlquery.InnerText(newConstruction), "New construction:", "", -1))
			}
		}
	}
}

func (zc *ZillowCrawler) ParseUtiGreenEnergyDetails(doc *html.Node) {
	// Utilities / Green Energy Details
	utiGreenEnergyDetails := htmlquery.Find(doc, "//h5[contains(text(), \"Utilities / Green Energy Details\")]/following-sibling::div//h6/following-sibling::ul")
	if utiGreenEnergyDetails != nil {
		for _, utiGreenEnergyDetail := range utiGreenEnergyDetails {
			// Sewer Information
			sewerInformation := htmlquery.FindOne(utiGreenEnergyDetail, ".//span[contains(text(), \"Sewer information\")]")
			if sewerInformation != nil {
				zc.CrawlerTables.ZillowData.SewerInformation = strings.TrimSpace(strings.Replace(htmlquery.InnerText(sewerInformation), "Sewer information:", "", -1))
			}

			// Water Information
			waterInformation := htmlquery.FindOne(utiGreenEnergyDetail, ".//span[contains(text(), \"Water information\")]")
			if waterInformation != nil {
				zc.CrawlerTables.ZillowData.WaterInformation = strings.TrimSpace(strings.Replace(htmlquery.InnerText(waterInformation), "Water information:", "", -1))
			}
		}
	}

}

func (zc *ZillowCrawler) ParseComNeiDetails(doc *html.Node) {
	// Community and Neighborhood Details
	comNeiDetails := htmlquery.Find(doc, "//h5[contains(text(), \"Community and Neighborhood Details\")]/following-sibling::div//h6/following-sibling::ul")
	if comNeiDetails != nil {
		for _, comNeiDetail := range comNeiDetails {
			// Region Location
			regionLocation := htmlquery.FindOne(comNeiDetail, ".//span[contains(text(), \"Region\")]")
			if regionLocation != nil {
				zc.CrawlerTables.ZillowData.RegionLocation = strings.TrimSpace(strings.Replace(htmlquery.InnerText(regionLocation), "Region:", "", -1))
			}

			// Subdivision
			subdivision := htmlquery.FindOne(comNeiDetail, ".//span[contains(text(), \"Subdivision\")]")
			if subdivision != nil {
				zc.CrawlerTables.ZillowData.Subdivision = strings.TrimSpace(strings.Replace(htmlquery.InnerText(subdivision), "Subdivision:", "", -1))
			}
		}
	}
}

func (zc *ZillowCrawler) ParseHoaFinancialDetails(doc *html.Node) {
	// HOA and financial details
	hoaFinancialDetails := htmlquery.Find(doc, "//h5[contains(text(), \"HOA and financial details\")]/following-sibling::div//h6/following-sibling::ul")
	if hoaFinancialDetails != nil {
		for _, hoaFinancialDetail := range hoaFinancialDetails {
			// Has HOA
			hasHoa := htmlquery.FindOne(hoaFinancialDetail, ".//span[contains(text(), \"Has HOA\")]")
			if hasHoa != nil {
				zc.CrawlerTables.ZillowData.HasHOA = strings.TrimSpace(strings.Replace(htmlquery.InnerText(hasHoa), "Has HOA:", "", -1))
			}

			// HOA Fee detail
			hoaFeeDetail := htmlquery.FindOne(hoaFinancialDetail, ".//span[contains(text(), \"HOA fee\")]")
			if hoaFeeDetail != nil {
				zc.CrawlerTables.ZillowData.HOAFeeDetail = strings.TrimSpace(strings.Replace(htmlquery.InnerText(hoaFeeDetail), "HOA fee:", "", -1))
			}

			// Services included
			servicesIncluded := htmlquery.FindOne(hoaFinancialDetail, ".//span[contains(text(), \"Services included\")]")
			if servicesIncluded != nil {
				zc.CrawlerTables.ZillowData.ServicesIncluded = strings.TrimSpace(strings.Replace(htmlquery.InnerText(servicesIncluded), "Services included:", "", -1))
			}

			// Association Name
			associationName := htmlquery.FindOne(hoaFinancialDetail, ".//span[contains(text(), \"Association name\")]")
			if associationName != nil {
				zc.CrawlerTables.ZillowData.AssociationName = strings.TrimSpace(strings.Replace(htmlquery.InnerText(associationName), "Association name:", "", -1))
			}

			// Association phone
			associationPhone := htmlquery.FindOne(hoaFinancialDetail, ".//span[contains(text(), \"Association phone\")]")
			if associationPhone != nil {
				zc.CrawlerTables.ZillowData.AssociationPhone = strings.TrimSpace(strings.Replace(htmlquery.InnerText(associationPhone), "Association phone:", "", -1))
			}

			//Annual tax amount
			annualTaxAmount := htmlquery.FindOne(hoaFinancialDetail, ".//span[contains(text(), \"Annual tax amount\")]")
			if annualTaxAmount != nil {
				zc.CrawlerTables.ZillowData.AnnualTaxAmount = strings.TrimSpace(strings.Replace(htmlquery.InnerText(annualTaxAmount), "Annual tax amount:", "", -1))
			}
		}
	}
}

func (zc *ZillowCrawler) ParseGreatSchools(doc *html.Node) {
	// GreatSchools rating
	greatSchoolsRating := htmlquery.Find(doc, "//*[@id=\"ds-nearby-schools-list\"]/li")
	if greatSchoolsRating != nil {
		for _, school := range greatSchoolsRating {
			// Elementary School
			elementarySchool := htmlquery.FindOne(school, ".//a[contains(text(), \"Elementary School\")]/following-sibling::span")
			if elementarySchool != nil {
				zc.CrawlerTables.ZillowData.ElementarySchool = strings.Replace(htmlquery.InnerText(elementarySchool), "Distance", ", Distance", -1)
			}

			// Middle School
			middleSchool := htmlquery.FindOne(school, ".//a[contains(text(), \"Middle School\")]/following-sibling::span")
			if middleSchool != nil {
				zc.CrawlerTables.ZillowData.MiddleSchool = strings.Replace(htmlquery.InnerText(middleSchool), "Distance", ", Distance", -1)
			}

			// High School
			highSchool := htmlquery.FindOne(school, ".//a[contains(text(), \"High School\")]/following-sibling::span")
			if highSchool != nil {
				zc.CrawlerTables.ZillowData.HighSchool = strings.Replace(htmlquery.InnerText(highSchool), "Distance", ", Distance", -1)
			}
		}
	}

}

func (zc *ZillowCrawler) ParseDistrict(doc *html.Node) {
	// District
	district := htmlquery.FindOne(doc, "//h5[contains(text(), \"Schools provided by the listing agent\")]/following-sibling::div/div[contains(text(), \"District\")]")
	if district != nil {
		zc.CrawlerTables.ZillowData.District = strings.TrimSpace(strings.Replace(htmlquery.InnerText(district), "District:", "", -1))
	}
}

func (zc *ZillowCrawler) ParseDataSource(doc *html.Node) {
	// Data Source
	dataSource := htmlquery.FindOne(doc, "//*[contains(text(), \"Find assessor info on the\")]/a/@href")
	if dataSource != nil {
		zc.CrawlerTables.ZillowData.DataSource = htmlquery.SelectAttr(dataSource, "href")
	}
}

func (zc *ZillowCrawler) ParseZillowPriceHistory(doc *html.Node) {
	priceHistories := htmlquery.Find(doc, "//h5[contains(text(), \"Price history\")]/following-sibling::table/tbody/tr[@label]")
	if priceHistories != nil {
		for _, history := range priceHistories {
			priceHistory := &entity.ZillowPriceHistory{}
			date := htmlquery.FindOne(history, "./td[1]/span")
			if date != nil {
				priceHistory.Date = strings.TrimSpace(htmlquery.InnerText(date))
			}
			event := htmlquery.FindOne(history, "./td[2]/span")
			if event != nil {
				priceHistory.Event = strings.TrimSpace(htmlquery.InnerText(event))
			}
			price := htmlquery.FindOne(history, "./td[3]/span[1]")
			if price != nil {
				priceHistory.Price = strings.TrimSpace(htmlquery.InnerText(price))

			}
			source := htmlquery.FindOne(history, "./following-sibling::tr[1]//span[contains(text(), \"Source\")]")
			priceHistory.Source = htmlquery.InnerText(source)
			priceHistory.Maindb3ID = zc.CrawlerTables.Maindb3.ID
			priceHistory.Address = zc.CrawlerTables.ZillowData.Address
			zc.CrawlerTables.ZillowPriceHistory = append(zc.CrawlerTables.ZillowPriceHistory, priceHistory)
		}
	}
}

func (zc *ZillowCrawler) ParseZillowPublicTaxHistory(doc *html.Node) {
	publicTaxHistories := htmlquery.Find(doc, "//h5[contains(text(), \"Public tax history\")]/following-sibling::table/tbody/tr")
	if publicTaxHistories != nil {
		for _, history := range publicTaxHistories {
			publicTaxHistory := &entity.ZillowPublicTaxHistory{}
			year := htmlquery.FindOne(history, "./td[1]")
			if year != nil {
				publicTaxHistory.Year = strings.TrimSpace(htmlquery.InnerText(year))
			}
			propertyTax := htmlquery.FindOne(history, "./td[2]")
			if propertyTax != nil {
				publicTaxHistory.PropertyTaxes = strings.TrimSpace(htmlquery.InnerText(propertyTax))
			}
			taxAssessment := htmlquery.FindOne(history, "./td[3]")
			if taxAssessment != nil {
				publicTaxHistory.TaxAssessment = strings.TrimSpace(htmlquery.InnerText(taxAssessment))
			}
			publicTaxHistory.Maindb3ID = zc.CrawlerTables.Maindb3.ID
			publicTaxHistory.Address = zc.CrawlerTables.ZillowData.Address
			zc.CrawlerTables.ZillowPublicTaxHistory = append(zc.CrawlerTables.ZillowPublicTaxHistory, publicTaxHistory)
		}
	}
}

func (zc *ZillowCrawler) ParseMapBounds(doc *html.Node) {
	script := htmlquery.FindOne(doc, "//script[@data-zrr-shared-data-key='mobileSearchPageStore']")
	if script == nil {
		zc.ShowLogError("Zillow Crawler didn't find the Map Bounds")
		return
	}
	dataScript := strings.Replace(script.FirstChild.Data, "<!--", "", -1)
	jsonString := strings.Replace(dataScript, "-->", "", -1)

	// Declared an empty map interface
	var result map[string]interface{}

	// Unmarshal or Decode the JSON to the interface.
	if err := json.Unmarshal([]byte(jsonString), &result); err != nil {
		zc.ShowLogError(err.Error())
		return
	}

	queryState := result["queryState"]
	if queryState == nil {
		zc.ShowLogError("Not found queryState in Json data")
		return
	}
	queryStateMap := queryState.(map[string]interface{})

	mapBounds := queryStateMap["mapBounds"]
	if mapBounds == nil {
		zc.ShowLogError("Not found mapBounds in Json data")
		return
	}
	mapBoundsMap := mapBounds.(map[string]interface{})

	zc.CrawlerTables.MapBounds = &schemas.MapBounds{
		West:  mapBoundsMap["west"].(float64),
		East:  mapBoundsMap["east"].(float64),
		South: mapBoundsMap["south"].(float64),
		North: mapBoundsMap["north"].(float64),
	}
}
