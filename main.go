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
		panic(err.Error())
	}
	logger.Info("database connected")

	// Int Session Manager
	sessionManager := initialization.IntSessionManager()

	// Int Router
	router := initialization.InitRouting(db, sessionManager)

	if err := ClientHandle(logger); err != nil {
		logger.Fatal(fmt.Sprintf("could not create task: %v", err))
	}

	logger.Info(fmt.Sprintf("Server START on port%v", viper.GetString("server.address")))
	log.Fatal(http.ListenAndServe(
		viper.GetString("server.address"),
		sessionManager.LoadAndSave(router),
	))
}

func ClientHandle(logger *zap.Logger) error {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: "127.0.0.1:6379"})
	defer client.Close()
	task, err := tasks.NewZillowCrawlerTask()
	if err != nil {
		return err
	}
	info, err := client.Enqueue(task)
	if err != nil {
		return err
	}
	logger.Info(fmt.Sprintf("enqueued task: id=%s queue=%s", info.ID, info.Queue))
	return nil
}
