package main

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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
		"",
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

	if err := PeriodicTasks(db, logger); err != nil {
		logger.Fatal(err.Error())
	}

	logger.Info(fmt.Sprintf("Server START on port%v", viper.GetString("server.address")))
	log.Fatal(http.ListenAndServe(
		viper.GetString("server.address"),
		sessionManager.LoadAndSave(router),
	))
}

func PeriodicTasks(db *gorm.DB, logger *zap.Logger) error {
	s := gocron.NewScheduler(time.UTC)
	s.SetMaxConcurrentJobs(viper.GetInt("crawler.workers.concurrent"), gocron.RescheduleMode)
	_, err := s.Every(viper.GetString("crawler.zillow_crawler.periodic_run")).SingletonMode().Do(tasks.RunCrawler, db, logger, false)
	if err != nil {
		return err
	}
	_, err = s.Every(viper.GetString("crawler.zillow_crawler.periodic_interval")).SingletonMode().Do(tasks.RunCrawler, db, logger, true)
	if err != nil {
		return err
	}
	s.StartAsync()
	return nil
}
