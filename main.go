package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"multilogin_scraping/helper"
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

	// Int RabbitMQ
	rabbitMQ := helper.NewRabbitMQ(viper.GetString("crawler.rabbitmq.url"), logger)
	defer rabbitMQ.CloseRabbitMQ()

	//Int Redis
	redis := helper.NewRedisCache(
		viper.GetString("crawler.redis.address"),
		"",
		viper.GetInt("crawler.redis.db"),
		logger,
	)

	// Int Router
	router := initialization.InitRouting(db, sessionManager, rabbitMQ, redis)

	logger.Info(fmt.Sprintf("Server START on port%v", viper.GetString("server.address")))

	log.Fatal(http.ListenAndServe(
		viper.GetString("server.address"),
		sessionManager.LoadAndSave(router),
	))

}
