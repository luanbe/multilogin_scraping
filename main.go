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
	logger := initialization.InitLogger()
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

	//RunCrawler(db)

	logger.Info(fmt.Sprintf("Server START on port%v", viper.GetString("server.address")))
	log.Fatal(http.ListenAndServe(
		viper.GetString("server.address"),
		sessionManager.LoadAndSave(router),
	))

}
func RunCrawler(db *gorm.DB, logger *zap.Logger) {
	c := colly.NewCollector()

	maindb3Service := registry.RegisterMaindb3Service(db)
	zillowService := registry.RegisterZillowService(db)
	maindb3List, err := maindb3Service.ListMaindb3Data(
		viper.GetString("crawler.crawler_status.succeeded"),
		viper.GetInt("crawler.zillow_crawler.concurrent"),
	)
	if err != nil {
		logger.Error(err.Error(), zap.String("Crawler", "Zillow"))
		return
	}
	if len(maindb3List) > 0 {
		defer logger.Info("ZillowCrawler: Crawled Done", zap.String("Crawler", "Zillow"))
	}
	for _, maindb3 := range maindb3List {
		if err := RunZillowCrawler(c, maindb3, zillowService, maindb3Service, logger); err != nil {
			logger.Error(err.Error(), zap.String("Crawler", "Zillow"))
		}
	}

}

func RunZillowCrawler(
	c *colly.Collector,
	maindb3 *entity.Maindb3,
	zillowService service.ZillowService,
	maindb3Service service.Maindb3Service,
	logger *zap.Logger,
) error {
	cZillow := c.Clone()
	zillowCrawler, err := zillow.NewZillowCrawler(cZillow, maindb3, zillowService, maindb3Service, logger)
	if err != nil {
		return err
	}
	if err := zillowCrawler.RunZillowCrawler(true); err != nil {
		return err
	}
	return nil
}
