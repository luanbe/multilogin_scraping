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
	crawlerRequest := &schemas.CrawlerRequest{}

	if err := render.Bind(r, crawlerRequest); err != nil {
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
		"data":    crawlerRequest,
		"worker":  viper.GetString("crawler.rabbitmq.tasks.crawl_address.routing_key"),
	}

	utils := helper.NewUtils()
	zillowStatus := "disabled"
	if crawlerRequest.CrawlersActive.Zillow == true {
		zillowStatus = viper.GetString("crawler.crawler_status.start")
	}
	realtorStatus := "disabled"
	if crawlerRequest.CrawlersActive.Realtor == true {
		realtorStatus = viper.GetString("crawler.crawler_status.start")
	}
	movotoStatus := "disabled"
	if crawlerRequest.CrawlersActive.Movoto == true {
		movotoStatus = viper.GetString("crawler.crawler_status.start")
	}
	crawlerSearchRes := schemas.CrawlerSearchRes{
		TaskID:         redisId.String(),
		CrawlerRequest: crawlerRequest,
		Zillow: &schemas.ZillowDetail{
			Status: zillowStatus,
			Error:  "",
		},
		Realtor: &schemas.RealtorDetail{
			Status: realtorStatus,
			Error:  "",
		},
		Movoto: &schemas.MovotoDetail{
			Status: movotoStatus,
			Error:  "",
		},
	}
	if err := cd.Redis.SetRedis(redisId.String(), crawlerSearchRes, time.Hour*1); err != nil {
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
	render.Render(w, r, schemas.NewCrawlerResponse(redisId.String(), crawlerSearchRes))
}
func (cd *CrawlerDelivery) CrawlerStatus(w http.ResponseWriter, r *http.Request) {
	taskID := chi.URLParam(r, "taskID")
	crawlerTask := &schemas.CrawlerSearchRes{}
	if err := cd.Redis.GetRedis(taskID, crawlerTask); err != nil {
		render.Render(w, r, schemas.ErrNotFound(err))
		return
	}
	render.Status(r, http.StatusOK)
	render.Render(w, r, crawlerTask)
}
