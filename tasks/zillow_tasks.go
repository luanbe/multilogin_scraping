package tasks

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"multilogin_scraping/app/models/entity"
	"multilogin_scraping/app/registry"
	"multilogin_scraping/crawlers/zillow"
	"sync"
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
	return &ZillowProcessor{DB: db, Logger: logger}
}

func NewZillowRedisTask() (*asynq.Task, error) {
	payload, err := json.Marshal(ZillowCrawlerPayload{})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeZillowCrawler, payload), nil
}

func RunCrawler(db *gorm.DB, zillowLogger *zap.Logger, onlyHistoryTable bool) {
	c := colly.NewCollector()
	maindb3Service := registry.RegisterMaindb3Service(db)
	zillowService := registry.RegisterZillowService(db)
	var wg sync.WaitGroup
	//noBrowser := viper.GetInt("crawler.zillow_crawler.no_browsers")
	// TODO: We will add more browsers later when finding the method query multiple records at per concurrency
	noBrowser := viper.GetInt("crawler.zillow_crawler.periodic_browser")
	recordSize := viper.GetInt("crawler.zillow_crawler.periodic_record_size")
	if onlyHistoryTable == true {
		noBrowser = viper.GetInt("crawler.zillow_crawler.periodic_browser_interval")
		recordSize = viper.GetInt("crawler.zillow_crawler.periodic_record_size_interval")
	}

	wg.Add(noBrowser)
	page := 0
	pageSize := recordSize / noBrowser
	for i := 0; i < noBrowser; i++ {
		page += 1
		maindb3DataList := []*entity.Maindb3{}
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
			zillowLogger.Error(err.Error())
		}
		if maindb3DataList == nil || len(maindb3DataList) < 1 {
			zillowLogger.Info("Not found maindb3 data")
			return
		}
		go func(maindb3DataList []*entity.Maindb3, wg *sync.WaitGroup) {
			defer wg.Done()
			for {
				zillowCrawler, err := zillow.NewZillowCrawler(c, maindb3DataList, zillowService, maindb3Service, zillowLogger, onlyHistoryTable)
				if err != nil {
					zillowLogger.Error(err.Error())
					continue
				}
				zillowLogger.Info(fmt.Sprint("Start crawler on Multilogin App: ", zillowCrawler.BaseSel.Profile.UUID))
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
		}(maindb3DataList, &wg)
	}
	wg.Wait()
	return
}
