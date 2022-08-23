package main

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"log"
	"multilogin_scraping/app/models/entity"
	"multilogin_scraping/app/registry"
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

	// Int Session Manager
	sessionManager := initialization.IntSessionManager()

	// Int Router
	router := initialization.InitRouting(db, sessionManager)

	RunCrawler(db)

	fmt.Printf("Server START on port%v\n", viper.GetString("server.address"))
	log.Fatal(http.ListenAndServe(
		viper.GetString("server.address"),
		sessionManager.LoadAndSave(router),
	))

}
func RunCrawler(db *gorm.DB) {
	c := colly.NewCollector()

	zillowService := registry.RegisterMaindb3Service(db)
	maindb3List := zillowService.ListMaindb3Data(
		viper.GetString("crawler.zillow_crawler.crawling_succeed_status"),
		viper.GetInt("crawler.zillow_crawler.concurrent"),
	)
	for _, maindb3 := range maindb3List {
		go RunZillowCrawler(c, maindb3)
	}

}

func RunZillowCrawler(c *colly.Collector, maindb3 *entity.Maindb3) {
	cZillow := c.Clone()
	zillowCrawler := zillow.NewZillowCrawler(cZillow, maindb3)
	zillowCrawler.RunZillowCrawler(true)
}
