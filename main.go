package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
	"multilogin_scraping/app/models/entity"
	"multilogin_scraping/app/registry"
	"multilogin_scraping/app/service"
	"multilogin_scraping/crawlers/zillow"
	"multilogin_scraping/initialization"
	"net/http"
	"sync"
)

func init() {
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	if viper.GetBool("debug") {
		fmt.Println("Service RUN on DEBUG mode")
	}
}

func main() {
	// Init Logger
	logger := initialization.InitLogger(
		map[string]interface{}{"Logger": "System"},
		"",
	)
	defer logger.Sync()

	// Init Db connection
	db, err := initialization.InitDb()
	if err != nil {
		logger.Fatal(fmt.Sprint("error Db connection: %v", err.Error()))
		panic(err.Error())
	}
	logger.Info("database connected")

	// Int Session Manager
	sessionManager := initialization.IntSessionManager()

	// Int Router
	router := initialization.InitRouting(db, sessionManager)

	RunCrawler(db)

	logger.Info(fmt.Sprintf("Server START on port%v", viper.GetString("server.address")))
	log.Fatal(http.ListenAndServe(
		viper.GetString("server.address"),
		sessionManager.LoadAndSave(router),
	))

}
func RunCrawler(db *gorm.DB) {
	c := colly.NewCollector()
	zillowLogger := initialization.InitLogger(
		map[string]interface{}{"Logger": "Zillow Crawler"},
		viper.GetString("crawler.zillow_crawler.log_file"),
	)
	defer zillowLogger.Sync()

	maindb3Service := registry.RegisterMaindb3Service(db)
	zillowService := registry.RegisterZillowService(db)
	maindb3List, err := maindb3Service.ListMaindb3Data(
		viper.GetString("crawler.crawler_status.succeeded"),
		viper.GetInt("crawler.zillow_crawler.concurrent"),
	)
	if err != nil {
		zillowLogger.Error(err.Error())
		return
	}
	var m sync.Mutex
	for _, maindb3 := range maindb3List {
		go RunZillowCrawler(c, maindb3, zillowService, maindb3Service, zillowLogger, &m)
	}
	if len(maindb3List) > 0 {
		defer zillowLogger.Info("Completed to crawl", zap.Int("No.Addresses", len(maindb3List)))
	}
}

func RunZillowCrawler(
	c *colly.Collector,
	maindb3 *entity.Maindb3,
	zillowService service.ZillowService,
	maindb3Service service.Maindb3Service,
	logger *zap.Logger,
	m *sync.Mutex,
) {
	cZillow := c.Clone()
	m.Lock()
	zillowCrawler, err := zillow.NewZillowCrawler(cZillow, maindb3, zillowService, maindb3Service, logger)
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
