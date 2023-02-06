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
	address string,
	proxy *util2.Proxy,
	crawlerTask *schemas.RealtorCrawlerTask,
	redis helper.RedisCache,
) {
	realtorCrawler := realtor.NewRealtorCrawler(rp.DB, rp.Logger, proxy)
	mprID, err := realtorCrawler.CrawlSearchData(address)
	if err != nil {
		rp.Logger.Fatal(err.Error())
	}

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

			crawlerTask.Error = err.Error()
			crawlerTask.Status = viper.GetString("crawler.crawler_status.failed")
			if err = redis.SetRedis(crawlerTask.TaskID, crawlerTask, time.Hour*1); err != nil {
				rp.Logger.Fatal(err.Error())
			}
		} else {
			crawlerTask.Status = viper.GetString("crawler.crawler_status.succeeded")
			crawlerTask.RealtorDetail = realtorCrawler.CrawlerSchemas.RealtorData
			if err := redis.SetRedis(crawlerTask.TaskID, crawlerTask, time.Hour*1); err != nil {
				rp.Logger.Fatal(err.Error())
			}
		}
		realtorCrawler.BaseSel.StopSessionBrowser(true)
		break
	}
}
