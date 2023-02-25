package tasks

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"multilogin_scraping/app/schemas"
	"multilogin_scraping/crawlers/realtor"
	"multilogin_scraping/helper"
	util2 "multilogin_scraping/pkg/utils"
	"time"
)

type RealtorProcessor struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

// NewRealtorApiTask begin to start a new task
func (rp RealtorProcessor) NewRealtorApiTask(
	crawlerSearchRes *schemas.CrawlerSearchRes,
	proxy *util2.Proxy,
	redis helper.RedisCache,
) {

	realtorCrawler := realtor.NewRealtorCrawler(rp.DB, rp.Logger, proxy)
	mprID, err := realtorCrawler.CrawlSearchData(crawlerSearchRes)
	if err != nil {
		crawlerSearchRes.Realtor.Error = err.Error()
		crawlerSearchRes.Realtor.Status = viper.GetString("crawler.crawler_status.failed")
	} else {
		for {
			// We will start new browser here
			// If browser creating is fail, use continue for creating new browser again
			if err := realtorCrawler.NewBrowser(); err != nil {
				rp.Logger.Error(err.Error())
				continue
			}

			rp.Logger.Info(fmt.Sprint("Start crawler on Multilogin App: ", realtorCrawler.BaseSel.Profile.UUID))
			if err := realtorCrawler.RunRealtorCrawlerAPI(mprID); err != nil {
				realtorCrawler.Logger.Error(err.Error())
				realtorCrawler.BaseSel.StopSessionBrowser(true)

				// if a browser is blocked or stopped, we will re-run it from a loop
				if realtorCrawler.BrowserTurnOff == true || realtorCrawler.CrawlerBlocked == true {
					continue
				}

				crawlerSearchRes.Realtor.Error = err.Error()
				crawlerSearchRes.Realtor.Status = viper.GetString("crawler.crawler_status.failed")

			} else {
				crawlerSearchRes.Realtor.Status = viper.GetString("crawler.crawler_status.succeeded")
				crawlerSearchRes.Realtor.Data = realtorCrawler.CrawlerTables.RealtorData
			}
			realtorCrawler.BaseSel.StopSessionBrowser(true)
			break
		}
	}
	if err = redis.SetRedis(crawlerSearchRes.TaskID, crawlerSearchRes, time.Hour*1); err != nil {
		rp.Logger.Fatal(err.Error())
	}
}
