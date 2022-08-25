package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/spf13/viper"
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
	// Init Db connection
	db, err := initialization.InitDb()
	if err != nil {
		fmt.Errorf("error Db connection: %v", err.Error())
		panic(err.Error())
	}

	// Init Logger
	logger := initialization.InitLogger()
	defer logger.Sync()

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
func RunCrawler(db *gorm.DB) {
	c := colly.NewCollector()

	maindb3Service := registry.RegisterMaindb3Service(db)
	zillowService := registry.RegisterZillowService(db)
	maindb3List := maindb3Service.ListMaindb3Data(
		viper.GetString("crawler.crawler_status.succeeded"),
		viper.GetInt("crawler.zillow_crawler.concurrent"),
	)
	if len(maindb3List) > 0 {
		defer fmt.Println("ZillowCrawler: Crawled Done")
	}
	for _, maindb3 := range maindb3List {
		RunZillowCrawler(c, maindb3, zillowService, maindb3Service)
	}

}

func RunZillowCrawler(c *colly.Collector, maindb3 *entity.Maindb3, zillowService service.ZillowService, maindb3Service service.Maindb3Service) {
	cZillow := c.Clone()
	zillowCrawler := zillow.NewZillowCrawler(cZillow, maindb3, zillowService, maindb3Service)
	if zillowCrawler == nil {
		return
	}
	zillowCrawler.RunZillowCrawler(true)
}
