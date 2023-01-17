package main

import (
	"fmt"
	"github.com/spf13/viper"
	"multilogin_scraping/app/schemas"
	"multilogin_scraping/helper"
	"multilogin_scraping/initialization"
	util2 "multilogin_scraping/pkg/utils"
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

func main2() {
	// Init logger
	workerLog := initialization.InitLogger(
		map[string]interface{}{"Logger": "Crawling Address"},
		viper.GetString("crawler.workers.log_file"),
	)
	// Init Redis
	redis := helper.NewRedisCache(
		viper.GetString("crawler.redis.address"),
		"",
		viper.GetInt("crawler.redis.db"),
		workerLog,
	)

	// Init Db connection
	db, err := initialization.InitDb()
	if err != nil {
		workerLog.Fatal(fmt.Sprint("error Db connection: %v", err.Error()))
	}
	workerLog.Info("database connected")

	zillowProcessor := tasks.ZillowProcessor{DB: db, Logger: workerLog}

	r := helper.NewRabbitMQ(viper.GetString("crawler.rabbitmq.url"), workerLog)
	messages, rabbitHelper := r.ConsumeMessage(
		viper.GetString("crawler.rabbitmq.tasks.crawl_address.exchange_type"),
		viper.GetString("crawler.rabbitmq.tasks.crawl_address.exchange_name"),
		viper.GetString("crawler.rabbitmq.tasks.crawl_address.queue_name"),
		viper.GetString("crawler.rabbitmq.tasks.crawl_address.routing_key"),
	)
	defer rabbitHelper.Connect.Close()
	defer rabbitHelper.Channel.Close()

	// Make a channel to receive messages into infinite loop.
	forever := make(chan bool)

	go func() {
		for message := range messages {
			// For example, show received message in a console.
			go workerLog.Info(fmt.Sprintf(" > Received message: %s\n", message.Body))
			utils := helper.NewUtils()
			//if body, err := utils.Deserialize(message.Body); err != nil {
			//	log.Printf(" > Errors: %s\n", err.Error())
			//}
			body, _ := utils.Deserialize(message.Body)

			if body["worker"] == viper.GetString("crawler.rabbitmq.tasks.crawl_address.routing_key") {
				crawlerTask := &schemas.ZillowCrawlerTask{}
				if err := redis.GetRedis(body["task_id"].(string), crawlerTask); err != nil {
					workerLog.Error(err.Error())
				} else {
					zillowProcessor.CrawlZillowDataByAPI(body["address"].(string), crawlerTask, redis)
				}

			}

		}
	}()

	<-forever
}

func main() {
	// Init logger
	workerLog := initialization.InitLogger(
		map[string]interface{}{"Logger": "Crawling Address"},
		viper.GetString("crawler.workers.log_file"),
	)
	db, err := initialization.InitDb()
	if err != nil {
		workerLog.Fatal(err.Error())
	}
	realtorProcessor := tasks.RealtorProcessor{DB: db, Logger: workerLog}
	realtorTask := &schemas.RealtorCrawlerTask{}
	redis := helper.NewRedisCache(
		viper.GetString("crawler.redis.address"),
		"",
		viper.GetInt("crawler.redis.db"),
		workerLog,
	)

	var proxies []util2.Proxy
	// load proxies file
	proxies, err = util2.GetProxies(viper.GetString("crawler.proxy_path"))
	if err != nil {
		workerLog.Fatal(fmt.Sprint("Loading proxy error:", err.Error()))
	}

	forever := make(chan bool)

	go func() {
		realtorProcessor.NewRealtorApiTask(
			"823 Lake Grove Dr, Little Elm,\" + \"TX 75068",
			proxies[util2.RandIntRange(0, len(proxies))],
			realtorTask,
			redis,
		)
	}()
	<-forever

}
