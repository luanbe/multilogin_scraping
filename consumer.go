package main

import (
	"multilogin_scraping/helper"
)

func main() {
	r := helper.NewRabbitMQ("amqp://root:root@127.0.0.1:5672/")
	r.ConsumeMessage("topic", "test_exchange", "test_queue", "test_exchange_key")
}
