package delivery

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/luanbe/golang-web-app-structure/app/service"
	"github.com/luanbe/golang-web-app-structure/templates"

	"github.com/luanbe/golang-web-app-structure/helper"
)

type UserDelivery struct {
	Tpl     helper.Template
	Service service.UserService
}

func NewUserDelivery(s service.UserService) *UserDelivery {
	return &UserDelivery{Service: s}
}

func (ud *UserDelivery) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(ud.UserContextBody)
	r.Post("/", ud.NewUser)
	r.Get("/signup/{username}", ud.Signup)
	r.Get("/login", ud.Login)
	r.Get("/logout", ud.Logout)

	return r
}

func (ud *UserDelivery) UserContextBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := []string{"luanbe68"}

		ctx := context.WithValue(r.Context(), "user", username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (ud *UserDelivery) Signup(w http.ResponseWriter, r *http.Request) {
	fmt.Println(chi.URLParam(r, "username"))
	w.Header().Set("Content-Type", "text/html; charset=utf8")
	ud.Tpl = helper.TplMust(ud.Tpl.TplParseFS(templates.FS, "signup.gohtml"))
	ud.Tpl.Execute(w, nil)
}

func (ud *UserDelivery) NewUser(w http.ResponseWriter, r *http.Request) {
	user := struct {
		Email    string
		UserName string
	}{
		r.FormValue("email"),
		r.FormValue("email"),
	}
	result, _ := ud.Service.AddUser(user.UserName, user.Email)
	fmt.Fprint(w, "New user:", result)

}

func (ud *UserDelivery) addSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		v, ok := ctx.Value("acl.admin").(bool)
		if !ok || !v {
			fmt.Println("err")
		}
		ctx = context.WithValue(ctx, "acl.admin", true)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (ud *UserDelivery) Login(w http.ResponseWriter, r *http.Request) {
	u := r.Context().Value("user").([]string)
	fmt.Fprint(w, u)
}

func (ud *UserDelivery) Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Logout")
}
