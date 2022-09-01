package zillow

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/net/html"
	"log"
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
	WebDriver       selenium.WebDriver
	BaseSel         *crawlers.BaseSelenium
	Profile         *crawlers.Profile
	CZillow         *colly.Collector
	SearchPageReq   *SearchPageReq
	CrawlerTables   *CrawlerTables
	CrawlerServices CrawlerServices
	Logger          *zap.Logger
}

type CrawlerTables struct {
	Maindb3    *entity.Maindb3
	ZillowData *entity.ZillowData
	//ZillowData *entry.ZillowData
}

type CrawlerServices struct {
	ZillowService  service.ZillowService
	Maindb3Service service.Maindb3Service
}

const searchURL = "https://www.zillow.com/search/GetSearchPageState.htm?searchQueryState=%s&wants={\"cat1\":[\"listResults\",\"mapResults\"],\"cat2\":[\"total\"],\"regionResults\":[\"total\"]}&requestId=5"

func NewZillowCrawler(
	c *colly.Collector,
	maindb3 *entity.Maindb3,
	zillowService service.ZillowService,
	maindb3Service service.Maindb3Service,
	logger *zap.Logger,
) (*ZillowCrawler, error) {
	BaseSel := crawlers.NewBaseSelenium()
	if err := BaseSel.StartSelenium("zillow"); err != nil {
		return nil, err
	}
	userAgent, err := BaseSel.WebDriver.ExecuteScript("return navigator.userAgent", nil)
	if err != nil {
		return nil, err
	}
	if userAgent == nil {
		userAgent = fake.UserAgent()
	}
	c.UserAgent = userAgent.(string)

	return &ZillowCrawler{
		WebDriver:       BaseSel.WebDriver,
		BaseSel:         BaseSel,
		Profile:         BaseSel.Profile,
		CZillow:         c,
		CrawlerTables:   &CrawlerTables{maindb3, &entity.ZillowData{}},
		CrawlerServices: CrawlerServices{zillowService, maindb3Service},
		Logger:          logger,
	}, nil
}

func (zc *ZillowCrawler) GetURLCrawling() string {
	address := fmt.Sprint(strings.TrimSpace(zc.CrawlerTables.Maindb3.OwnerAddress), ", ", zc.CrawlerTables.Maindb3.OwnerCityState)
	address = strings.Replace(address, " ", "-", -1)
	return fmt.Sprint("https://www.zillow.com/homes/", address, "_rb/")
}

func (zc *ZillowCrawler) ShowLogError(mes string) {
	zc.Logger.Error(mes, zap.Uint64("mainDBID", zc.CrawlerTables.Maindb3.ID), zap.String("URL", zc.GetURLCrawling()))
}

func (zc *ZillowCrawler) ShowLogInfo(mes string) {
	zc.Logger.Info(mes, zap.Uint64("mainDBID", zc.CrawlerTables.Maindb3.ID), zap.String("URL", zc.GetURLCrawling()))
}

func (zc *ZillowCrawler) RunZillowCrawler(exactAddress bool) error {
	defer zc.BaseSel.StopSelenium()
	// Crawling exact address
	if exactAddress == true {
		zc.CrawlerTables.ZillowData.URL = zc.GetURLCrawling()
		if err := zc.CrawlAddress(zc.CrawlerTables.ZillowData.URL); err != nil {
			return err
		}
	}
	// defer zc.BaseSel.StopSelenium()
	//if err := zc.WebDriver.Get(url); err != nil {
	//	log.Fatalln(err)
	//}
	//pageSource, err := zc.WebDriver.PageSource()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//zc.CheckVerifyHuman(pageSource)
	//zc.CrawlData(zc.CrawlMapBounds(pageSource))
	//fmt.Println(zc.WebDriver.GetCookies())
	return nil
}

func (zc *ZillowCrawler) CheckVerifyHuman(pageSource string) error {
	if strings.Contains(pageSource, "Please verify you're a human to continue") {
		return fmt.Errorf("the website blocked crawler")
	}
	return nil
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
				zc.CollectionData(result)
			}

			// Crawling data on Next Page
			zc.SearchPageReq.Pagination.CurrentPage += 1
			searchNextPageJson, err := json.Marshal(zc.SearchPageReq)
			if err != nil {
				log.Fatalln(err)
			}
			urlNextPage := fmt.Sprintf(searchURL, string(searchNextPageJson))
			fmt.Println("Crawling Next Page: ", urlNextPage)
			r.Request.Visit(urlNextPage)
		}
		//If len = 0 => crawl done!
		return
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

func (zc *ZillowCrawler) CollectionData(result SearchPageResResult) {
	propertyStatus := false
	if result.Beds > 0 || result.Baths > 0 {
		propertyStatus = true
	}
	halfBathRooms := result.HdpData.HomeInfo.Bedrooms / 2
	fullBathRooms := result.HdpData.HomeInfo.Bathrooms - halfBathRooms
	zc.CrawlerTables.ZillowData = &entity.ZillowData{
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
	zc.CrawlAddress(zc.CrawlerTables.ZillowData.URL)
}

func (zc *ZillowCrawler) CrawlAddress(address string) error {
	if err := zc.WebDriver.Get(address); err != nil {
		return err
	}
	// NOTE: time to load source. Need to increase if data was not showing
	time.Sleep(10 * time.Second)
	pageSource, err := zc.WebDriver.PageSource()
	if err != nil {
		return err
	}
	if err := zc.CheckVerifyHuman(pageSource); err != nil {
		return err
	}
	if err := zc.ParseData(pageSource); err != nil {
		return err
	}
	if err := zc.UpdateDB(); err != nil {
		return err
	}
	return nil
}

func (zc *ZillowCrawler) UpdateDB() error {

	if err := zc.CrawlerServices.ZillowService.AddZillow(zc.CrawlerTables.ZillowData); err != nil {
		return err
	}

	if err := zc.CrawlerServices.Maindb3Service.UpdateStatus(zc.CrawlerTables.Maindb3, viper.GetString("crawler.crawler_status.succeeded")); err != nil {
		return err
	}
	return nil
}

func (zc *ZillowCrawler) ParseData(source string) error {
	//htmlquery.DisableSelectorCache = true
	doc, err := htmlquery.Parse(strings.NewReader(source))

	if err != nil {
		return err
	}

	zc.ParseURL()
	zc.ParseBedBathSF(doc)
	zc.ParseAddress(doc)
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

	//// Year Built
	//yearBuilt := htmlquery.FindOne(doc, "//span[contains(text(), \"Year built\")]/text()")
	//if yearBuilt != nil {
	//	zillowData.YearBuilt = strings.TrimSpace(strings.Replace(yearBuilt.Data, "Year built:", "", -1))
	//}
	//
	//// Natural Gas
	//naturalGas := htmlquery.FindOne(doc, "//*[contains(text(), \"Natural Gas\") or contains(text(), \"natural gas\")]")
	//if naturalGas != nil {
	//	zillowData.NaturalGas = true
	//}
	//
	//// Central Air
	//centralAir := htmlquery.FindOne(doc, "//*[contains(text(), \"Central Air\") or contains(text(), \"central air\")]")
	//if centralAir != nil {
	//	zillowData.CentralAir = true
	//}
	//
	//// # of Garage Spaces
	//garageSpaces := htmlquery.FindOne(doc, "//*[contains(text(), \"garage spaces\")]/text()")
	//if garageSpaces != nil {
	//	zillowData.OfGarageSpaces = strings.TrimSpace(strings.Replace(garageSpaces.Data, " garage spaces", "", -1))
	//}
	//
	//// HOA Amount
	//hoaAmount := htmlquery.FindOne(doc, "//*[contains(text(), \"annually HOA fee\")]/text()")
	//if hoaAmount != nil {
	//	zillowData.HOAAmount = strings.TrimSpace(strings.Replace(hoaAmount.Data, " annually HOA fee", "", -1))
	//}
	//
	//// Lot Size
	//lotSizes := htmlquery.Find(doc, "//*[contains(text(), \"Lot size\")]/text()")
	//if lotSizes != nil {
	//	for _, v := range lotSizes {
	//		lotSizeData := strings.TrimSpace(strings.Replace(v.Data, "Lot size:", "", -1))
	//		if strings.Contains(lotSizeData, "sqft") == true {
	//			zillowData.LotSizeSF = lotSizeData
	//		}
	//		if strings.Contains(lotSizeData, "Acres") == true {
	//			zillowData.LotSizeAcres = lotSizeData
	//		}
	//	}
	//
	//}
	//
	//// Buyer's agent fee
	//buyerAgentFee := htmlquery.FindOne(doc, "//*[contains(text(), \"buyer's agent fee\")]/text()")
	//if buyerAgentFee != nil {
	//	zillowData.BuyerAgentFee = strings.TrimSpace(strings.Replace(buyerAgentFee.Data, " buyer's agent fee", "", -1))
	//}
	//
	//// Appliances
	//applicances := htmlquery.FindOne(doc, "//*[contains(text(), \"Appliances included\")]")
	//if applicances != nil {
	//	zillowData.Appliances = strings.TrimSpace(strings.Replace(htmlquery.InnerText(applicances), "Appliances included:", "", -1))
	//}
	//
	//// Living Room
	//livingRooms := htmlquery.Find(doc, "//h6[contains(text(), \"LivingRoom\")]/following-sibling::ul")
	//for _, livingroom := range livingRooms {
	//	// Living Room Level
	//	livingRoomLevel := htmlquery.FindOne(livingroom, ".//span[contains(text(), \"Level\")]")
	//	if livingRoomLevel != nil {
	//		zillowData.LivingRoomLevel = strings.TrimSpace(strings.Replace(htmlquery.InnerText(livingRoomLevel), "Level:", "", -1))
	//	}
	//
	//	// Living Room Dimensions
	//	livingRoomDimensions := htmlquery.FindOne(livingroom, ".//span[contains(text(), \"Dimensions\")]")
	//	if livingRoomDimensions != nil {
	//		zillowData.LivingRoomDimensions = strings.TrimSpace(strings.Replace(htmlquery.InnerText(livingRoomDimensions), "Dimensions:", "", -1))
	//	}
	//}
	//
	//// Primary Bedroom
	//primaryBedRooms := htmlquery.Find(doc, "//h6[contains(text(), \"PrimaryBedroom\")]/following-sibling::ul")
	//for _, primaryBedRoom := range primaryBedRooms {
	//	// Primary Bedroom Level
	//	primaryBedRoomLevel := htmlquery.FindOne(primaryBedRoom, ".//span[contains(text(), \"Level\")]")
	//	if primaryBedRoomLevel != nil {
	//		zillowData.PrimaryBedroomLevel = strings.TrimSpace(strings.Replace(htmlquery.InnerText(primaryBedRoomLevel), "Level:", "", -1))
	//	}
	//
	//	// Primary Bedroom Dimensions
	//	primaryBedRoomDimensions := htmlquery.FindOne(primaryBedRoom, ".//span[contains(text(), \"Dimensions\")]")
	//	if primaryBedRoomDimensions != nil {
	//		zillowData.PrimaryBedroomDimensions = strings.TrimSpace(strings.Replace(htmlquery.InnerText(primaryBedRoomDimensions), "Dimensions:", "", -1))
	//	}
	//}
	//
	//// Interior Features
	//interiorFeatures := htmlquery.FindOne(doc, "//span[contains(text(), \"Interior features\")]")
	//if interiorFeatures != nil {
	//	zillowData.InteriorFeatures = strings.TrimSpace(strings.Replace(htmlquery.InnerText(interiorFeatures), "Interior features:", "", -1))
	//}
	//
	//// Basement
	//basement := htmlquery.FindOne(doc, "//span[contains(text(), \"Basement\")]")
	//if basement != nil {
	//	zillowData.Basement = strings.TrimSpace(strings.Replace(htmlquery.InnerText(basement), "Basement:", "", -1))
	//}
	//
	//// Total Interior Livable Area SF
	//totalInteriorLivableArea := htmlquery.FindOne(doc, "//span[contains(text(), \"Total interior livable area\")]")
	//if totalInteriorLivableArea != nil {
	//	zillowData.TotalInteriorLivableAreaSF = strings.TrimSpace(strings.Replace(htmlquery.InnerText(totalInteriorLivableArea), "Total interior livable area:", "", -1))
	//}
	//
	//// # of Fireplaces
	//offFireplaces := htmlquery.FindOne(doc, "//span[contains(text(), \"Total number of fireplaces\")]")
	//if offFireplaces != nil {
	//	zillowData.OfFireplaces = strings.TrimSpace(strings.Replace(htmlquery.InnerText(offFireplaces), "Total number of fireplaces:", "", -1))
	//}
	//
	//// Fireplace features
	//fireplaceFeatures := htmlquery.FindOne(doc, "//span[contains(text(), \"Fireplace features\")]")
	//if fireplaceFeatures != nil {
	//	zillowData.FireplaceFeatures = strings.TrimSpace(strings.Replace(htmlquery.InnerText(fireplaceFeatures), "Fireplace features:", "", -1))
	//}
	//
	//// Flooring Type
	//flooringType := htmlquery.FindOne(doc, "//h6[contains(text(), \"Flooring\")]/following-sibling::ul//span[contains(text(), \"Flooring\")]")
	//if flooringType != nil {
	//	zillowData.FlooringType = strings.TrimSpace(strings.Replace(htmlquery.InnerText(flooringType), "Flooring:", "", -1))
	//}
	//
	//// Heating Type
	//heatingType := htmlquery.FindOne(doc, "//h6[contains(text(), \"Heating\")]/following-sibling::ul//span[contains(text(), \"Heating features\")]")
	//if heatingType != nil {
	//	zillowData.HeatingType = strings.TrimSpace(strings.Replace(htmlquery.InnerText(heatingType), "Heating features:", "", -1))
	//}
	//
	//// Parking
	//parkings := htmlquery.Find(doc, "//h6[contains(text(), \"Parking\")]/following-sibling::ul")
	//if parkings != nil {
	//	for _, parking := range parkings {
	//		// Total Parking Spaces
	//		totalParkingSpaces := htmlquery.FindOne(parking, ".//span[contains(text(), \"Total spaces\")]")
	//		if totalParkingSpaces != nil {
	//			zillowData.TotalParkingSpaces = strings.TrimSpace(strings.Replace(htmlquery.InnerText(totalParkingSpaces), "Total spaces:", "", -1))
	//		}
	//
	//		// Parking Features
	//		parkingFeatures := htmlquery.FindOne(parking, ".//span[contains(text(), \"Parking features\")]")
	//		if parkingFeatures != nil {
	//			zillowData.ParkingFeatures = strings.TrimSpace(strings.Replace(htmlquery.InnerText(parkingFeatures), "Parking features:", "", -1))
	//		}
	//
	//		// Covered Spaces
	//		coveredSpaces := htmlquery.FindOne(parking, ".//span[contains(text(), \"Covered spaces\")]")
	//		if coveredSpaces != nil {
	//			zillowData.CoveredSpaces = strings.TrimSpace(strings.Replace(htmlquery.InnerText(coveredSpaces), "Covered spaces:", "", -1))
	//		}
	//
	//	}
	//}
	//
	//// Lot Features
	//lotFeatures := htmlquery.FindOne(doc, "//h6[contains(text(), \"Lot\")]/following-sibling::ul//span[contains(text(), \"Lot features\")]")
	//if lotFeatures != nil {
	//	zillowData.LotFeatures = strings.TrimSpace(strings.Replace(htmlquery.InnerText(lotFeatures), "Lot features:", "", -1))
	//}
	//
	//// Parcel number
	//parcelNumber := htmlquery.FindOne(doc, "//h6[contains(text(), \"Other property information\")]/following-sibling::ul//span[contains(text(), \"Parcel number\")]")
	//if parcelNumber != nil {
	//	zillowData.ParcelNumber = strings.TrimSpace(strings.Replace(htmlquery.InnerText(parcelNumber), "Parcel number:", "", -1))
	//}
	//
	//// Property details - Property
	//propertydetails := htmlquery.Find(doc, "//h5[contains(text(), \"Property details\")]/following-sibling::div//h6[contains(text(), \"Property\")]/following-sibling::ul")
	//if propertydetails != nil {
	//	for _, property := range propertydetails {
	//		// # Levels (Stories/Floors)
	//		levelsStoriesFloors := htmlquery.FindOne(property, ".//span[contains(text(), \"Levels\")]")
	//		if levelsStoriesFloors != nil {
	//			zillowData.LevelsStoriesFloors = strings.TrimSpace(strings.Replace(htmlquery.InnerText(levelsStoriesFloors), "Levels:", "", -1))
	//		}
	//
	//		// Patio and Porch Details
	//		patioAndPorchDetails := htmlquery.FindOne(property, ".//span[contains(text(), \"Patio and porch details\")]")
	//		if patioAndPorchDetails != nil {
	//			zillowData.PatioAndPorchDetails = strings.TrimSpace(strings.Replace(htmlquery.InnerText(patioAndPorchDetails), "Patio and porch details:", "", -1))
	//		}
	//
	//	}
	//}
	//
	//// Construction details
	//constructionDetails := htmlquery.Find(doc, "//h5[contains(text(), \"Construction details\")]/following-sibling::div//h6/following-sibling::ul")
	//if constructionDetails != nil {
	//	for _, constructionDetail := range constructionDetails {
	//		// HomeType
	//		homeType := htmlquery.FindOne(constructionDetail, ".//span[contains(text(), \"Home type\")]")
	//		if homeType != nil {
	//			zillowData.HomeType = strings.TrimSpace(strings.Replace(htmlquery.InnerText(homeType), "Home type:", "", -1))
	//		}
	//		// Propery SubType
	//		propertySubType := htmlquery.FindOne(constructionDetail, ".//span[contains(text(), \"Property subType\")]")
	//		if propertySubType != nil {
	//			zillowData.ProperySubType = strings.TrimSpace(strings.Replace(htmlquery.InnerText(propertySubType), "Property subType:", "", -1))
	//		}
	//
	//		// Construction Materials
	//		constructionMaterials := htmlquery.FindOne(constructionDetail, ".//span[contains(text(), \"Construction materials\")]")
	//		if constructionMaterials != nil {
	//			zillowData.ConstructionMaterials = strings.TrimSpace(strings.Replace(htmlquery.InnerText(constructionMaterials), "Construction materials:", "", -1))
	//		}
	//
	//		// Foundation
	//		foundation := htmlquery.FindOne(constructionDetail, ".//span[contains(text(), \"Foundation\")]")
	//		if foundation != nil {
	//			zillowData.Foundation = strings.TrimSpace(strings.Replace(htmlquery.InnerText(foundation), "Foundation:", "", -1))
	//		}
	//
	//		// Roof
	//		roof := htmlquery.FindOne(constructionDetail, ".//span[contains(text(), \"Roof\")]")
	//		if roof != nil {
	//			zillowData.Roof = strings.TrimSpace(strings.Replace(htmlquery.InnerText(roof), "Roof:", "", -1))
	//		}
	//
	//		// New Construction
	//		newConstruction := htmlquery.FindOne(constructionDetail, ".//span[contains(text(), \"New construction\")]")
	//		if newConstruction != nil {
	//			zillowData.NewConstruction = strings.TrimSpace(strings.Replace(htmlquery.InnerText(newConstruction), "New construction:", "", -1))
	//		}
	//	}
	//}
	//
	//// Utilities / Green Energy Details
	//utiGreenEnergyDetails := htmlquery.Find(doc, "//h5[contains(text(), \"Utilities / Green Energy Details\")]/following-sibling::div//h6/following-sibling::ul")
	//if utiGreenEnergyDetails != nil {
	//	for _, utiGreenEnergyDetail := range utiGreenEnergyDetails {
	//		// Sewer Information
	//		sewerInformation := htmlquery.FindOne(utiGreenEnergyDetail, ".//span[contains(text(), \"Sewer information\")]")
	//		if sewerInformation != nil {
	//			zillowData.SewerInformation = strings.TrimSpace(strings.Replace(htmlquery.InnerText(sewerInformation), "Sewer information:", "", -1))
	//		}
	//
	//		// Water Information
	//		waterInformation := htmlquery.FindOne(utiGreenEnergyDetail, ".//span[contains(text(), \"Water information\")]")
	//		if waterInformation != nil {
	//			zillowData.WaterInformation = strings.TrimSpace(strings.Replace(htmlquery.InnerText(waterInformation), "Water information:", "", -1))
	//		}
	//	}
	//}
	//
	//// Community and Neighborhood Details
	//comNeiDetails := htmlquery.Find(doc, "//h5[contains(text(), \"Community and Neighborhood Details\")]/following-sibling::div//h6/following-sibling::ul")
	//if comNeiDetails != nil {
	//	for _, comNeiDetail := range comNeiDetails {
	//		// Region Location
	//		regionLocation := htmlquery.FindOne(comNeiDetail, ".//span[contains(text(), \"Region\")]")
	//		if regionLocation != nil {
	//			zillowData.RegionLocation = strings.TrimSpace(strings.Replace(htmlquery.InnerText(regionLocation), "Region:", "", -1))
	//		}
	//
	//		// Subdivision
	//		subdivision := htmlquery.FindOne(comNeiDetail, ".//span[contains(text(), \"Subdivision\")]")
	//		if subdivision != nil {
	//			zillowData.Subdivision = strings.TrimSpace(strings.Replace(htmlquery.InnerText(subdivision), "Subdivision:", "", -1))
	//		}
	//	}
	//}
	//
	//// HOA and financial details
	//hoaFinancialDetails := htmlquery.Find(doc, "//h5[contains(text(), \"HOA and financial details\")]/following-sibling::div//h6/following-sibling::ul")
	//if hoaFinancialDetails != nil {
	//	for _, hoaFinancialDetail := range hoaFinancialDetails {
	//		// Has HOA
	//		hasHoa := htmlquery.FindOne(hoaFinancialDetail, ".//span[contains(text(), \"Has HOA\")]")
	//		if hasHoa != nil {
	//			zillowData.HasHOA = strings.TrimSpace(strings.Replace(htmlquery.InnerText(hasHoa), "Has HOA:", "", -1))
	//		}
	//
	//		// HOA Fee detail
	//		hoaFeeDetail := htmlquery.FindOne(hoaFinancialDetail, ".//span[contains(text(), \"HOA fee\")]")
	//		if hoaFeeDetail != nil {
	//			zillowData.HOAFeeDetail = strings.TrimSpace(strings.Replace(htmlquery.InnerText(hoaFeeDetail), "HOA fee:", "", -1))
	//		}
	//
	//		// Services included
	//		servicesIncluded := htmlquery.FindOne(hoaFinancialDetail, ".//span[contains(text(), \"Services included\")]")
	//		if servicesIncluded != nil {
	//			zillowData.ServicesIncluded = strings.TrimSpace(strings.Replace(htmlquery.InnerText(servicesIncluded), "Services included:", "", -1))
	//		}
	//
	//		// Association Name
	//		associationName := htmlquery.FindOne(hoaFinancialDetail, ".//span[contains(text(), \"Association name\")]")
	//		if associationName != nil {
	//			zillowData.AssociationName = strings.TrimSpace(strings.Replace(htmlquery.InnerText(associationName), "Association name:", "", -1))
	//		}
	//
	//		// Association phone
	//		associationPhone := htmlquery.FindOne(hoaFinancialDetail, ".//span[contains(text(), \"Association phone\")]")
	//		if associationPhone != nil {
	//			zillowData.AssociationPhone = strings.TrimSpace(strings.Replace(htmlquery.InnerText(associationPhone), "Association phone:", "", -1))
	//		}
	//
	//		//Annual tax amount
	//		annualTaxAmount := htmlquery.FindOne(hoaFinancialDetail, ".//span[contains(text(), \"Annual tax amount\")]")
	//		if annualTaxAmount != nil {
	//			zillowData.AnnualTaxAmount = strings.TrimSpace(strings.Replace(htmlquery.InnerText(annualTaxAmount), "Annual tax amount:", "", -1))
	//		}
	//	}
	//}
	//
	//// GreatSchools rating
	//greatSchoolsRating := htmlquery.Find(doc, "//*[@id=\"ds-nearby-schools-list\"]/li")
	//if greatSchoolsRating != nil {
	//	for _, school := range greatSchoolsRating {
	//		// Elementary School
	//		elementarySchool := htmlquery.FindOne(school, ".//a[contains(text(), \"Elementary School\")]/following-sibling::span")
	//		if elementarySchool != nil {
	//			zillowData.ElementarySchool = strings.Replace(htmlquery.InnerText(elementarySchool), "Distance", ", Distance", -1)
	//		}
	//
	//		// Middle School
	//		middleSchool := htmlquery.FindOne(school, ".//a[contains(text(), \"Middle School\")]/following-sibling::span")
	//		if middleSchool != nil {
	//			zillowData.MiddleSchool = strings.Replace(htmlquery.InnerText(middleSchool), "Distance", ", Distance", -1)
	//		}
	//
	//		// High School
	//		highSchool := htmlquery.FindOne(school, ".//a[contains(text(), \"High School\")]/following-sibling::span")
	//		if highSchool != nil {
	//			zillowData.HighSchool = strings.Replace(htmlquery.InnerText(highSchool), "Distance", ", Distance", -1)
	//		}
	//	}
	//}
	//
	//// District
	//district := htmlquery.FindOne(doc, "//h5[contains(text(), \"Schools provided by the listing agent\")]/following-sibling::div/div[contains(text(), \"District\")]")
	//if district != nil {
	//	zillowData.District = strings.TrimSpace(strings.Replace(htmlquery.InnerText(district), "District:", "", -1))
	//}
	//
	//// Data Source
	//dataSource := htmlquery.FindOne(doc, "//*[contains(text(), \"Find assessor info on the\")]/a/@href")
	//if dataSource != nil {
	//	zillowData.DataSource = htmlquery.SelectAttr(dataSource, "href")
	//}
	zc.CrawlerTables.ZillowData.Maindb3ID = zc.CrawlerTables.Maindb3.ID

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
						zc.CrawlerTables.ZillowData.Bed = bed
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
						zc.CrawlerTables.ZillowData.Bath = bath
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
					zc.Logger.Error(err.Error())
				} else {
					zc.CrawlerTables.ZillowData.FullBathrooms = fullBathRoomValue
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
					zc.Logger.Error(err.Error())
				} else {
					zc.CrawlerTables.ZillowData.HalfBathrooms = halfBathRoomValue
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
					zc.Logger.Error(err.Error())
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
					zc.Logger.Error(err.Error())
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
						zc.Logger.Error(err.Error())
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
				zc.Logger.Error(err.Error())
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
				zc.Logger.Error(err.Error())
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
