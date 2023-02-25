package tasks

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"multilogin_scraping/app/schemas"
	"multilogin_scraping/crawlers/movoto"
	"multilogin_scraping/helper"
	util2 "multilogin_scraping/pkg/utils"
	"time"
)

type MovotoProcessor struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

// NewMovotoApiTask begin to start a new task
func (rp MovotoProcessor) NewMovotoApiTask(
	crawlerSearchRes *schemas.CrawlerSearchRes,
	proxy *util2.Proxy,
	redis helper.RedisCache,
) {
	movotoCrawler := movoto.NewMovotoCrawler(rp.DB, rp.Logger, proxy)

	for {
		// We will start new browser here
		// If browser creating is fail, use continue for creating new browser again
		if err := movotoCrawler.NewBrowser(); err != nil {
			rp.Logger.Error(err.Error())
			continue
		}
		rp.Logger.Info(fmt.Sprint("Start crawler on Multilogin App: ", movotoCrawler.BaseSel.Profile.UUID))
		searchRes, err := movotoCrawler.CrawlSearchData(crawlerSearchRes)

		if err != nil {
			crawlerSearchRes.Movoto.Error = err.Error()
			crawlerSearchRes.Movoto.Status = viper.GetString("crawler.crawler_status.failed")

		} else {
			if err := movotoCrawler.RunMovotoCrawlerAPI(searchRes); err != nil {
				movotoCrawler.Logger.Error(err.Error())
				movotoCrawler.BaseSel.StopSessionBrowser(true)

				// if a browser is blocked or stopped, we will re-run it from a loop
				if movotoCrawler.BrowserTurnOff == true || movotoCrawler.CrawlerBlocked == true {
					continue
				}

				crawlerSearchRes.Movoto.Error = err.Error()
				crawlerSearchRes.Movoto.Status = viper.GetString("crawler.crawler_status.failed")
			} else {
				crawlerSearchRes.Movoto.Status = viper.GetString("crawler.crawler_status.succeeded")
				crawlerSearchRes.Movoto.Data = movotoCrawler.CrawlerTables.MovotoData

			}
		}
		movotoCrawler.BaseSel.StopSessionBrowser(true)
		break
	}
	if err := redis.SetRedis(crawlerSearchRes.TaskID, crawlerSearchRes, time.Hour*1); err != nil {
		rp.Logger.Fatal(err.Error())
	}
}
