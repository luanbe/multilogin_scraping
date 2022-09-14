package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"multilogin_scraping/app/models/entity"
	"multilogin_scraping/app/registry"
	"multilogin_scraping/app/service"
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
func (processor *ZillowProcessor) ZillowRedisProcessTask(ctx context.Context, t *asynq.Task) error {

	var p ZillowCrawlerPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}
	RunCrawler(processor.DB, processor.Logger)
	return nil
}

func RunCrawler(db *gorm.DB, zillowLogger *zap.Logger) {
	c := colly.NewCollector()

	maindb3Service := registry.RegisterMaindb3Service(db)
	zillowService := registry.RegisterZillowService(db)
	// Comment this for testing ID
	maindb3List, err := maindb3Service.ListMaindb3Data(
		viper.GetString("crawler.crawler_status.succeeded"),
		viper.GetInt("crawler.zillow_crawler.no_browsers"),
	)
	if err != nil {
		zillowLogger.Error(err.Error())
		return
	}
	var m sync.Mutex
	for _, maindb3 := range maindb3List {
		go RunZillowCrawler(c, maindb3, zillowService, maindb3Service, zillowLogger, &m, false)
	}
	if len(maindb3List) > 0 {
		defer zillowLogger.Info("Completed to crawl", zap.Int("No.Addresses", len(maindb3List)))
	}

	//NOTE: This is for testing data by id
	//maindb3, err := maindb3Service.GetMaindb3(1)
	//if err != nil {
	//	zillowLogger.Error(err.Error())
	//	return
	//}
	//var m sync.Mutex
	//go RunZillowCrawler(c, maindb3, zillowService, maindb3Service, zillowLogger, &m)

}

func RunCrawlerInterval(db *gorm.DB, zillowLogger *zap.Logger) {
	c := colly.NewCollector()

	maindb3Service := registry.RegisterMaindb3Service(db)
	zillowService := registry.RegisterZillowService(db)
	// Comment this for testing ID
	maindb3List, err := maindb3Service.ListMaindb3IntervalData(
		viper.GetInt("crawler.zillow_crawler.days_interval"),
		viper.GetString("crawler.crawler_status.succeeded"),
		viper.GetInt("crawler.zillow_crawler.no_browsers_interval"),
	)
	if err != nil {
		zillowLogger.Error(err.Error())
		return
	}
	var m sync.Mutex
	for _, maindb3 := range maindb3List {
		go RunZillowCrawler(c, maindb3, zillowService, maindb3Service, zillowLogger, &m, true)
	}
	if len(maindb3List) > 0 {
		defer zillowLogger.Info("Completed to crawl & updated history table", zap.Int("No.Addresses", len(maindb3List)))
	}

	//NOTE: This is for testing data by id
	//maindb3, err := maindb3Service.GetMaindb3(1)
	//if err != nil {
	//	zillowLogger.Error(err.Error())
	//	return
	//}
	//var m sync.Mutex
	//go RunZillowCrawler(c, maindb3, zillowService, maindb3Service, zillowLogger, &m)

}

func RunZillowCrawler(
	c *colly.Collector,
	maindb3 *entity.Maindb3,
	zillowService service.ZillowService,
	maindb3Service service.Maindb3Service,
	logger *zap.Logger,
	m *sync.Mutex,
	onlyHistoryTable bool,
) {
	cZillow := c.Clone()
	m.Lock()
	zillowCrawler, err := zillow.NewZillowCrawler(cZillow, maindb3, zillowService, maindb3Service, logger, onlyHistoryTable)
	if err != nil {
		logger.Error(err.Error(), zap.Uint64("mainDBID", maindb3.ID))
		m.Unlock()
		return
	}
	m.Unlock()
	zillowCrawler.ShowLogInfo("Zillow Data is crawling...")
	if err := zillowCrawler.RunZillowCrawler(true); err != nil {
		zillowCrawler.ShowLogError(err.Error())
	}
	zillowCrawler.ShowLogInfo("Zillow Data Crawled!")
}
