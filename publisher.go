package main

import (
	"multilogin_scraping/helper"
)

func main() {
	r := helper.NewRabbitMQ("amqp://root:root@127.0.0.1:5672/")
	r.PublishMessage("topic", "test_exchange", "test_exchange_key", "test 1 lan coi sao")
}
