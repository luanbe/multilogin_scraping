package delivery

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/luanbe/golang-web-app-structure/helper"
	"github.com/luanbe/golang-web-app-structure/templates"
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
