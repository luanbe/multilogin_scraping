package main

import (
	"fmt"
	"github.com/hibiken/asynq"
	"log"
	"multilogin_scraping/initialization"
	"multilogin_scraping/tasks"
)

const redisAddr = "127.0.0.1:6379"

func main2() {
	logger := initialization.InitLogger(
		map[string]interface{}{"Logger": "Worker"},
		"./logs/workers.log",
	)

	// Init Db connection
	db, err := initialization.InitDb()
	if err != nil {
		logger.Fatal(fmt.Sprintf("error Db connection: %v", err.Error()))
		panic(err.Error())
	}
	logger.Info("database connected")

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			// Specify how many concurrent workers to use
			Concurrency: 10,
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
