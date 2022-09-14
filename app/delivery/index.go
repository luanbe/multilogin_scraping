package delivery

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"multilogin_scraping/helper"
	"multilogin_scraping/templates"
)

type IndexDelivery struct {
	Tpl helper.Template
}

func NewIndexDelivery() *IndexDelivery {
	return &IndexDelivery{}
}

func (index *IndexDelivery) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/add", index.Home)

	return r
}

func (index *IndexDelivery) Home(w http.ResponseWriter, r *http.Request) {
	index.Tpl = helper.TplMust(index.Tpl.TplParseFS(templates.FS, "home.gohtml"))
	index.Tpl.Execute(w, nil)
}
