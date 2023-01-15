package api

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gofrs/uuid"
	"github.com/spf13/viper"
	"multilogin_scraping/app/schemas"
	"multilogin_scraping/helper"
	"net/http"
	"time"
)

type CrawlerDelivery struct {
	RabbitMQ helper.RabbitMQBroker
	Redis    helper.RedisCache
}

func NewCrawlerDelivery() *CrawlerDelivery {
	return &CrawlerDelivery{}
}

func (cd *CrawlerDelivery) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/crawl", cd.Crawl)
	r.Get("/status/{taskID}", cd.CrawlerStatus)
	return r
}

func (cd *CrawlerDelivery) Crawl(w http.ResponseWriter, r *http.Request) {
	data := &schemas.CrawlerRequest{}

	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, schemas.ErrInvalidRequest(err))
		return
	}
	redisId, err := uuid.NewV4()
	if err != nil {
		render.Render(w, r, schemas.ErrServer(fmt.Errorf("Redis: %v", err.Error())))
		return
	}
	msg := helper.Message{
		"task_id": redisId.String(),
		"address": data.Address,
		"worker":  viper.GetString("crawler.rabbitmq.tasks.crawl_address.routing_key"),
	}

	utils := helper.NewUtils()

	object := schemas.ZillowCrawlerTask{
		Status:  viper.GetString("crawler.crawler_status.start"),
		TaskID:  redisId.String(),
		Address: data.Address,
		Error:   "",
	}
	if err := cd.Redis.SetRedis(redisId.String(), object, time.Hour*1); err != nil {
		render.Render(w, r, schemas.ErrServer(err))
		return
	}

	dataByte, err := utils.Serialize(msg)
	if err := cd.RabbitMQ.PublishMessage(
		viper.GetString("crawler.rabbitmq.tasks.crawl_address.exchange_type"),
		viper.GetString("crawler.rabbitmq.tasks.crawl_address.exchange_name"),
		viper.GetString("crawler.rabbitmq.tasks.crawl_address.routing_key"),
		dataByte,
	); err != nil {
		render.Render(w, r, schemas.ErrServer(fmt.Errorf("RabbitMQ: %v", err.Error())))
		return
	}

	render.Status(r, http.StatusOK)
	render.Render(w, r, schemas.NewCrawlerResponse(redisId.String(), data.Address))
}
func (cd *CrawlerDelivery) CrawlerStatus(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	crawlerTask := &schemas.ZillowCrawlerTask{}
	if err := cd.Redis.GetRedis(taskID, crawlerTask); err != nil {
		render.Render(w, r, schemas.ErrNotFound(err))
		return
	}
	render.Status(r, http.StatusOK)
	render.Render(w, r, crawlerTask)

}
