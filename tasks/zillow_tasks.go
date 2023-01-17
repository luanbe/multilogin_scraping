package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"multilogin_scraping/app/models/entity"
	"multilogin_scraping/app/registry"
	"multilogin_scraping/app/schemas"
	"multilogin_scraping/crawlers/zillow"
	"multilogin_scraping/helper"
	util2 "multilogin_scraping/pkg/utils"
	"sync"
	"time"
)

const (
	TypeZillowCrawler = "zillow:crawler"
)

type ZillowCrawlerPayload struct {
}

type ZillowProcessor struct {
	DB     *gorm.DB
	Logger *zap.Logger
}

func NewZillowProcessor(db *gorm.DB, logger *zap.Logger) *ZillowProcessor {
	return &ZillowProcessor{
		DB:     db,
		Logger: logger,
	}
}

func NewZillowRedisTask() (*asynq.Task, error) {
	payload, err := json.Marshal(ZillowCrawlerPayload{})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeZillowCrawler, payload), nil
}

func (Zillow ZillowProcessor) CrawlZillowData(onlyHistoryTable bool) {

	var wg sync.WaitGroup
	//noBrowser := viper.GetInt("crawler.zillow_crawler.no_browsers")
	// TODO: We will add more browsers later when finding the method query multiple records at per concurrency
	noBrowser := viper.GetInt("crawler.zillow_crawler.periodic_browser")
	recordSize := viper.GetInt("crawler.zillow_crawler.periodic_record_size")
	if onlyHistoryTable == true {
		noBrowser = viper.GetInt("crawler.zillow_crawler.periodic_browser_interval")
		recordSize = viper.GetInt("crawler.zillow_crawler.periodic_record_size_interval")
	}
	var proxies []util2.Proxy
	// load proxies file
	proxies, err := util2.GetProxies(viper.GetString("crawler.zillow_crawler.proxy_path"))
	if err != nil {
		Zillow.Logger.Fatal(fmt.Sprint("Loading proxy error:", err.Error()))
	}

	maindb3Service := registry.RegisterMaindb3Service(Zillow.DB)

	wg.Add(noBrowser)
	page := 0
	pageSize := recordSize / noBrowser
	for i := 0; i < noBrowser; i++ {
		page += 1
		maindb3DataList := []*entity.ZillowMaindb3Address{}
		err := error(nil)
		if onlyHistoryTable == true {
			maindb3DataList, err = maindb3Service.ListMaindb3IntervalData(
				viper.GetInt("crawler.zillow_crawler.days_interval"),
				viper.GetString("crawler.crawler_status.succeeded"),
				page,
				pageSize,
			)
		} else {
			maindb3DataList, err = maindb3Service.ListMaindb3Data(
				//viper.GetString("crawler.crawler_status.failed"),
				"",
				page,
				pageSize,
			)
		}

		if err != nil {
			Zillow.Logger.Error(err.Error())
		}
		if maindb3DataList == nil || len(maindb3DataList) < 1 {
			Zillow.Logger.Info("Not found maindb3 data")
			return
		}
		// Begin to crawl data
		go Zillow.RunCrawler(onlyHistoryTable, maindb3DataList, &wg, proxies[i])
	}
	wg.Wait()
	return
}
func (Zillow ZillowProcessor) RunCrawler(
	onlyHistoryTable bool,
	maindb3DataList []*entity.ZillowMaindb3Address,
	wg *sync.WaitGroup,
	proxy util2.Proxy,
) {
	defer wg.Done()

	for {
		zillowCrawler, err := zillow.NewZillowCrawler(
			Zillow.DB,
			maindb3DataList,
			Zillow.Logger,
			onlyHistoryTable,
			proxy,
		)
		if err != nil {
			Zillow.Logger.Error(err.Error())
			continue
		}
		Zillow.Logger.Info(fmt.Sprint("Start crawler on Multilogin App: ", zillowCrawler.BaseSel.Profile.UUID))
		if err := zillowCrawler.RunZillowCrawler(); err != nil {
			zillowCrawler.ShowLogError(err.Error())
			zillowCrawler.BaseSel.StopSessionBrowser(true)
			if zillowCrawler.CrawlerBlocked == true {
				continue
			}
			if zillowCrawler.BrowserTurnOff == true {
				continue
			}
			break
		}
		zillowCrawler.BaseSel.StopSessionBrowser(true)
		break
	}
}

func (zp ZillowProcessor) CrawlZillowDataByAPI(address string, crawlerTask *schemas.ZillowCrawlerTask, redis helper.RedisCache) {
	var proxies []util2.Proxy
	// load proxies file
	proxies, err := util2.GetProxies(viper.GetString("crawler.proxy_path"))
	if err != nil {
		zp.Logger.Fatal(fmt.Sprint("Loading proxy error:", err.Error()))
	}

	go zp.RunZillowCrawlerAPI(address, proxies[util2.RandIntRange(0, len(proxies))], crawlerTask, redis)
}

func (zp ZillowProcessor) RunZillowCrawlerAPI(
	address string,
	proxy util2.Proxy,
	crawlerTask *schemas.ZillowCrawlerTask,
	redis helper.RedisCache,
) {
	// Delare a empty slice but we will not use it
	var maindb3DataList []*entity.ZillowMaindb3Address

	for {
		zillowCrawler, err := zillow.NewZillowCrawler(
			zp.DB,
			maindb3DataList,
			zp.Logger,
			false,
			proxy,
		)
		if err != nil {
			zp.Logger.Error(err.Error())
			continue
		}
		zp.Logger.Info(fmt.Sprint("Start crawler on Multilogin App: ", zillowCrawler.BaseSel.Profile.UUID))
		if err := zillowCrawler.RunZillowCrawlerAPI(address); err != nil {
			zillowCrawler.ShowLogError(err.Error())
			zillowCrawler.BaseSel.StopSessionBrowser(true)
			if zillowCrawler.CrawlerBlocked == true {
				continue
			}
			if zillowCrawler.BrowserTurnOff == true {
				continue
			}
			crawlerTask.Error = err.Error()
			crawlerTask.Status = viper.GetString("crawler.crawler_status.failed")
			redis.SetRedis(crawlerTask.TaskID, crawlerTask, time.Hour*1)
		} else {
			crawlerTask.Status = viper.GetString("crawler.crawler_status.succeeded")
			crawlerTask.ZillowDetail = zillowCrawler.CrawlerTables.ZillowData
			redis.SetRedis(crawlerTask.TaskID, crawlerTask, time.Hour*1)
		}
		zillowCrawler.BaseSel.StopSessionBrowser(true)
		break
	}
}
