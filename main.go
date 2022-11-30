package main

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"log"
	"multilogin_scraping/initialization"
	"multilogin_scraping/tasks"
	"net/http"
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

	// Int Session Manager
	sessionManager := initialization.IntSessionManager()

	// Int Router
	router := initialization.InitRouting(db, sessionManager)

	if err := PeriodicTasks(db); err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info(fmt.Sprintf("Server START on port%v", viper.GetString("server.address")))
	log.Fatal(http.ListenAndServe(
		viper.GetString("server.address"),
		sessionManager.LoadAndSave(router),
	))

}

func PeriodicTasks(db *gorm.DB) error {
	s := gocron.NewScheduler(time.UTC)
	zillowLogger := initialization.InitLogger(
		map[string]interface{}{"Logger": "Zillow"},
		viper.GetString("crawler.zillow_crawler.log_file"),
	)

	s.SetMaxConcurrentJobs(viper.GetInt("crawler.workers.concurrent"), gocron.RescheduleMode)

	zillowProcessor := tasks.ZillowProcessor{DB: db, Logger: zillowLogger}
	//_, err := s.Every(viper.GetString("crawler.zillow_crawler.periodic_run")).SingletonMode().Do(zillowProcessor.CrawlZillowData, false)
	//if err != nil {
	//	return err
	//}
	_, err := s.Every(viper.GetString("crawler.zillow_crawler.periodic_interval")).SingletonMode().Do(zillowProcessor.CrawlZillowData, true)
	if err != nil {
		return err
	}
	s.StartAsync()
	return nil
}
