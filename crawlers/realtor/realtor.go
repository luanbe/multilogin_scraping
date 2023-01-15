package realtor

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/icrowley/fake"
	"github.com/spf13/viper"
	"github.com/tebeka/selenium"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"multilogin_scraping/app/schemas"
	"multilogin_scraping/crawlers"
	util2 "multilogin_scraping/pkg/utils"
	"strings"
	"time"
)

type RealtorCrawler struct {
	WebDriver               selenium.WebDriver
	BaseSel                 *crawlers.BaseSelenium
	Profile                 *crawlers.Profile
	CZillow                 *colly.Collector
	Logger                  *zap.Logger
	CrawlerBlocked          bool
	BrowserTurnOff          bool
	SearchDataByCollyStatus bool
	CrawlerTables           *CrawlerTables
}

type CrawlerTables struct {
	RealtorData *schemas.RealtorData
}

func NewRealtorCrawler(
	db *gorm.DB,
	logger *zap.Logger,
	proxy util2.Proxy,
) (*RealtorCrawler, error) {
	BaseSel := crawlers.NewBaseSelenium(logger)
	if err := BaseSel.StartSelenium("realtor", proxy, viper.GetBool("crawler.realtor_crawler.proxy_status")); err != nil {
		return nil, err
	}

	// Disable image loading
	if viper.GetBool("crawler.disable_load_images") == true {
		if BaseSel.Profile.BrowserName == "stealthfox" {
			if err := BaseSel.FireFoxDisableImageLoading(); err != nil {
				return nil, err
			}
		}
	}

	userAgent, err := BaseSel.WebDriver.ExecuteScript("return navigator.userAgent", nil)
	if err != nil {
		return nil, err
	}

	if userAgent == nil {
		userAgent = fake.UserAgent()
	}
	c := colly.NewCollector()
	c.UserAgent = userAgent.(string)

	return &RealtorCrawler{
		WebDriver:      BaseSel.WebDriver,
		BaseSel:        BaseSel,
		Profile:        BaseSel.Profile,
		CZillow:        c,
		Logger:         logger,
		CrawlerBlocked: false,
		BrowserTurnOff: false,
	}, nil
}

func (rc *RealtorCrawler) RunRealtorCrawlerAPI(address string) error {
	for {
		err := func() error {
			rc.Logger.Info("Zillow Data is crawling...")
			rc.CrawlerTables.RealtorData.URL = viper.GetString("crawler.realtor_crawler.url")
			if err := rc.CrawlAddressAPI(rc.CrawlerTables.RealtorData.URL); err != nil {
				rc.Logger.Error(err.Error())
				rc.Logger.Error("Failed to crawl data")
				if rc.CrawlerBlocked == true {
					return err
				}
				// TODO: Update error for crawling here
			}
			rc.Logger.Info("Completed to crawl data")
			return nil
		}()
		if err != nil {
			return err
		}
		return nil
	}
}

func (rc *RealtorCrawler) CrawlAddressAPI(address string) error {
	if err := rc.WebDriver.Get(address); err != nil {
		return err
	}
	// NOTE: time to load source. Need to increase if data was not showing
	time.Sleep(viper.GetDuration("crawler.realtor_crawler.time_load_source") * time.Second)
	pageSource, err := rc.WebDriver.PageSource()
	if err != nil {
		rc.BrowserTurnOff = true
		return err
	}
	if err := rc.ByPassVerifyHuman(pageSource, address); err != nil {
		return err
	}
	// TODO: Add Parse Data
	//if err := rc.ParseData(pageSource); err != nil {
	//	return err
	//}
	return nil
}

func (rc *RealtorCrawler) ByPassVerifyHuman(pageSource string, url string) error {
	if rc.IsVerifyHuman(pageSource) == true {
		rc.CrawlerBlocked = true

	}
	if rc.CrawlerBlocked == true {
		for i := 0; i < 3; i++ {
			err := rc.WebDriver.Get(url)
			if err != nil {
				return err
			}
			pageSource, _ = rc.WebDriver.PageSource()
			if rc.IsVerifyHuman(pageSource) == false {
				rc.CrawlerBlocked = false
				return nil
			}
		}
		return fmt.Errorf("Crawler blocked for checking verify hunman")
	}
	return nil
}
func (rc *RealtorCrawler) IsVerifyHuman(pageSource string) bool {
	if strings.Contains(pageSource, "Please verify you're a human to continue") || strings.Contains(pageSource, "Let's confirm you are human") {
		return true
	}
	return false
}
