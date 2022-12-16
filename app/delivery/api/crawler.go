package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"multilogin_scraping/app/schemas"
	"multilogin_scraping/helper"
	"net/http"
)

type CrawlerDelivery struct {
}

func NewCrawlerDelivery() *CrawlerDelivery {
	return &CrawlerDelivery{}
}

func (cd *CrawlerDelivery) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/crawl", cd.Crawl)
	return r
}

func (cd *CrawlerDelivery) Crawl(w http.ResponseWriter, r *http.Request) {
	data := &schemas.CrawlerRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, schemas.ErrInvalidRequest(err))
		return
	}
	req := helper.NewRabbitMQ("amqp://root:root@127.0.0.1:5672/")
	req.PublishMessage("topic", "test_exchange", "test_exchange_key", data.Address)

	render.Status(r, http.StatusOK)
	render.Render(w, r, schemas.NewCrawlerResponse(data.Address))
}
