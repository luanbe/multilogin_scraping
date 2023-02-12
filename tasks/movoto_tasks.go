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
	address string,
	proxy *util2.Proxy,
	crawlerTask *schemas.MovotoCrawlerTask,
	redis helper.RedisCache,
) {
	movotoCrawler := movoto.NewMovotoCrawler(rp.DB, rp.Logger, proxy)
	mprID, err := movotoCrawler.CrawlSearchData(address)
	if err != nil {
		rp.Logger.Fatal(err.Error())
	}

	for {
		// We will start new browser here
		// If browser creating is fail, use continue for creating new browser again
		if err := movotoCrawler.NewBrowser(); err != nil {
			rp.Logger.Error(err.Error())
			continue
		}

		rp.Logger.Info(fmt.Sprint("Start crawler on Multilogin App: ", movotoCrawler.BaseSel.Profile.UUID))
		if err := movotoCrawler.RunMovotoCrawlerAPI(mprID); err != nil {
			movotoCrawler.Logger.Error(err.Error())
			movotoCrawler.BaseSel.StopSessionBrowser(true)

			// if a browser is blocked or stopped, we will re-run it from a loop
			if movotoCrawler.BrowserTurnOff == true || movotoCrawler.CrawlerBlocked == true {
				continue
			}

			crawlerTask.Error = err.Error()
			crawlerTask.Status = viper.GetString("crawler.crawler_status.failed")
			if err = redis.SetRedis(crawlerTask.TaskID, crawlerTask, time.Hour*1); err != nil {
				rp.Logger.Fatal(err.Error())
			}
		} else {
			crawlerTask.Status = viper.GetString("crawler.crawler_status.succeeded")
			crawlerTask.MovotoDetail = movotoCrawler.CrawlerSchemas.MovotoData
			if err := redis.SetRedis(crawlerTask.TaskID, crawlerTask, time.Hour*1); err != nil {
				rp.Logger.Fatal(err.Error())
			}
		}
		movotoCrawler.BaseSel.StopSessionBrowser(true)
		break
	}
}
