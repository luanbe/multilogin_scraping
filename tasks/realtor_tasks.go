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

func (rp RealtorProcessor) CrawlRealtorDataByAPI(address string, crawlerTask *schemas.RealtorCrawlerTask, redis helper.RedisCache) {
	var proxies []util2.Proxy
	// load proxies file
	proxies, err := util2.GetProxies(viper.GetString("crawler.proxy_path"))
	if err != nil {
		rp.Logger.Fatal(fmt.Sprint("Loading proxy error:", err.Error()))
	}

	go rp.RunRealtorCrawlerAPI(address, proxies[util2.RandIntRange(0, len(proxies))], crawlerTask, redis)
}

func (rp RealtorProcessor) RunRealtorCrawlerAPI(
	address string,
	proxy util2.Proxy,
	crawlerTask *schemas.RealtorCrawlerTask,
	redis helper.RedisCache,
) {
	for {
		realtorCrawler, err := realtor.NewRealtorCrawler(
			rp.DB,
			rp.Logger,
			proxy,
		)
		if err != nil {
			rp.Logger.Error(err.Error())
			continue
		}
		rp.Logger.Info(fmt.Sprint("Start crawler on Multilogin App: ", realtorCrawler.BaseSel.Profile.UUID))
		if err := realtorCrawler.RunRealtorCrawlerAPI(address); err != nil {
			realtorCrawler.Logger.Error(err.Error())
			realtorCrawler.BaseSel.StopSessionBrowser(true)
			if realtorCrawler.CrawlerBlocked == true {
				continue
			}
			if realtorCrawler.BrowserTurnOff == true {
				continue
			}
			crawlerTask.Error = err.Error()
			crawlerTask.Status = viper.GetString("crawler.crawler_status.failed")
			redis.SetRedis(crawlerTask.TaskID, crawlerTask, time.Hour*1)
		} else {
			crawlerTask.Status = viper.GetString("crawler.crawler_status.succeeded")
			crawlerTask.RealtorDetail = realtorCrawler.CrawlerTables.RealtorData
			redis.SetRedis(crawlerTask.TaskID, crawlerTask, time.Hour*1)
		}
		realtorCrawler.BaseSel.StopSessionBrowser(true)
		break
	}
}
