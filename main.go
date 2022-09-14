package main

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/hibiken/asynq"
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

	if viper.GetBool("crawler.workers.redis_task") == true {
		if err := RedisPeriodicTasks(logger); err != nil {
			logger.Fatal(fmt.Sprintf("could not create periodic task: %v", err))
		}
	} else {
		if err := PeriodicTasks(db, logger); err != nil {
			logger.Fatal(fmt.Sprintf("could not create periodic task: %v", err))
		}
	}

	logger.Info(fmt.Sprintf("Server START on port%v", viper.GetString("server.address")))
	log.Fatal(http.ListenAndServe(
		viper.GetString("server.address"),
		sessionManager.LoadAndSave(router),
	))
}

func RedisPeriodicTasks(logger *zap.Logger) error {
	// Example of using America/Los_Angeles timezone instead of the default UTC timezone.
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}
	scheduler := asynq.NewScheduler(
		asynq.RedisClientOpt{
			Addr: viper.GetString("crawler.redis.address"),
			DB:   viper.GetInt("crawler.redis.db"),
		},
		&asynq.SchedulerOpts{
			Location: loc,
		},
	)
	zillowCrawlerTask, err := tasks.NewZillowRedisTask()
	if err != nil {
		return err
	}
	entryID, err := scheduler.Register(fmt.Sprint("@every ", viper.GetString("crawler.zillow_crawler.periodic_run")), zillowCrawlerTask)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("registered an entry: %q\n", entryID))

	if err := scheduler.Run(); err != nil {
		return err
	}

	return nil
}

func PeriodicTasks(db *gorm.DB, logger *zap.Logger) error {
	s := gocron.NewScheduler(time.UTC)
	s.SetMaxConcurrentJobs(viper.GetInt("crawler.workers.concurrent"), gocron.RescheduleMode)
	_, err := s.Every(viper.GetString("crawler.zillow_crawler.periodic_run")).Do(tasks.RunCrawler, db, logger)
	if err != nil {
		return err
	}
	_, err = s.Every(viper.GetString("crawler.zillow_crawler.periodic_interval")).Do(tasks.RunCrawlerInterval, db, logger)
	if err != nil {
		return err
	}
	s.StartAsync()
	return nil
}
