package main

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"multilogin_scraping/initialization"
	"multilogin_scraping/tasks"
	"time"
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
		"system.log",
	)
	// Init Db connection
	db, err := initialization.InitDb()
	if err != nil {
		logger.Fatal(fmt.Sprint("error Db connection: %v", err.Error()))
	}
	logger.Info("database connected")

	if err := PeriodicTasks(db); err != nil {
		logger.Fatal(err.Error())
	}

}

func PeriodicTasks(db *gorm.DB) error {
	s := gocron.NewScheduler(time.UTC)
	//zillowLogger := initialization.InitLogger(
	//	map[string]interface{}{"Logger": "Zillow"},
	//	viper.GetString("crawler.zillow_crawler.log_file"),
	//)

	s.SetMaxConcurrentJobs(viper.GetInt("crawler.workers.concurrent"), gocron.RescheduleMode)

	//zillowProcessor := tasks.ZillowProcessor{DB: db, Logger: zillowLogger}
	//_, err := s.Every(viper.GetString("crawler.zillow_crawler.periodic_run")).SingletonMode().Do(zillowProcessor.CrawlZillowData, false)
	//if err != nil {
	//	return err
	//}
	//_, err := s.Every(viper.GetString("crawler.zillow_crawler.periodic_interval")).SingletonMode().Do(zillowProcessor.CrawlZillowData, true)

	multiLoginLogger := initialization.InitLogger(
		map[string]interface{}{"Logger": "Zillow"},
		viper.GetString("crawler.zillow_crawler.log_file"),
	)
	multiLoginProcessor := tasks.MultiLoginProcessor{Logger: multiLoginLogger}
	_, err := s.Every("5m").SingletonMode().Do(multiLoginProcessor.DeleteProfiles, []string{"movoto-Crawler", "zillow-Crawler", "realtor-Crawler"})
	if err != nil {
		return err
	}
	s.StartAsync()
	return nil
}
