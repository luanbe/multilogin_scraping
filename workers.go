package main

import (
	"fmt"
	"github.com/hibiken/asynq"
	"github.com/spf13/viper"
	"log"
	"multilogin_scraping/initialization"
	"multilogin_scraping/tasks"
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
	logger := initialization.InitLogger(
		map[string]interface{}{"Logger": "Worker"},
		viper.GetString("crawler.workers.log_file"),
	)

	// Init Db connection
	db, err := initialization.InitDb()
	if err != nil {
		logger.Fatal(fmt.Sprintf("error Db connection: %v", err.Error()))
	}
	logger.Info("database connected")

	srv := asynq.NewServer(
		asynq.RedisClientOpt{
			Addr: viper.GetString("crawler.redis.address"),
			DB:   viper.GetInt("crawler.redis.db"),
		},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: viper.GetInt("crawler.workers.concurrent"),
			// Optionally specify multiple queues with different priority.
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
			// See the godoc for other configuration options
		},
	)

	// mux maps a type to a handler
	mux := asynq.NewServeMux()
	zillowProcessor := tasks.NewZillowProcessor(db, logger)
	mux.HandleFunc(tasks.TypeZillowCrawler, zillowProcessor.ZillowCrawlerProcessTask)

	if err := srv.Run(mux); err != nil {
		log.Fatalf("could not run server: %v", err)
	}
}
