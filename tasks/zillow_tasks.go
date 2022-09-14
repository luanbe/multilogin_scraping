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
<<<<<<< HEAD
	// Comment this for testing ID
	maindb3List, err := maindb3Service.ListMaindb3Data(
		viper.GetString("crawler.crawler_status.succeeded"),
		viper.GetInt("crawler.zillow_crawler.no_browsers"),
=======
	maindb3List, err := maindb3Service.ListMaindb3Data(
		viper.GetString("crawler.crawler_status.succeeded"),
