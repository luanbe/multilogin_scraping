package delivery_admin

import (
	"github.com/alexedwards/scs/v2"
	"github.com/luanbe/golang-web-app-structure/app/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/luanbe/golang-web-app-structure/helper"
	"github.com/luanbe/golang-web-app-structure/templates"
)

type UserAdminDelivery struct {
	tpl            helper.Template
	SessionManager *scs.SessionManager
	UserService    service.UserService
}

func NewUserAdminDelivery(userService service.UserService, sessionManager *scs.SessionManager) *UserAdminDelivery {
	return &UserAdminDelivery{
		UserService:    userService,
		SessionManager: sessionManager,
	}
}

func (uad *UserAdminDelivery) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/add", uad.NewUser)

	return r
}

func (uad *UserAdminDelivery) NewUser(w http.ResponseWriter, r *http.Request) {

	uad.tpl = helper.TplMust(uad.tpl.TplParseFS(templates.FS, "admin/user/user_add.gohtml"))
	uad.tpl.Execute(w, nil)
}
