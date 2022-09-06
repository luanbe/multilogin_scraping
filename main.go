package main

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
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

	if err := PeriodicTasks(logger); err != nil {
		logger.Fatal(fmt.Sprintf("could not create periodic task: %v", err))
	}

	logger.Info(fmt.Sprintf("Server START on port%v", viper.GetString("server.address")))
	log.Fatal(http.ListenAndServe(
		viper.GetString("server.address"),
		sessionManager.LoadAndSave(router),
	))
}

func ClientHandle(logger *zap.Logger) error {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: viper.GetString("crawler.redis.address"),
		DB:   viper.GetInt("crawler.redis.db"),
	})
	defer client.Close()
	task, err := tasks.NewZillowCrawlerTask()
	if err != nil {
		return err
	}
	info, err := client.Enqueue(task, asynq.ProcessIn(viper.GetDuration("crawler.zillow_crawler.periodic_run")*time.Minute))
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("enqueued task: id=%s queue=%s", info.ID, info.Queue))
	return nil
}

func PeriodicTasks(logger *zap.Logger) error {
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
	zillowCrawlerTask, err := tasks.NewZillowCrawlerTask()
	if err != nil {
		return err
	}
	entryID, err := scheduler.Register(viper.GetString("crawler.zillow_crawler.periodic_run"), zillowCrawlerTask)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("registered an entry: %q\n", entryID))

	if err := scheduler.Run(); err != nil {
		return err
	}

	return nil
}
