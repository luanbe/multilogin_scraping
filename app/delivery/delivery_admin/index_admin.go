package delivery_admin

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/luanbe/golang-web-app-structure/app/middlewares"
	"github.com/spf13/viper"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/luanbe/golang-web-app-structure/app/service"
	"github.com/luanbe/golang-web-app-structure/helper"
	"github.com/luanbe/golang-web-app-structure/templates"
	"golang.org/x/crypto/bcrypt"
)

type IndexAdminDelivery struct {
	Tpl            helper.Template
	SessionManager *scs.SessionManager
	UserService    service.UserService
}

func NewIndexAdminDelivery(userService service.UserService, sessionManager *scs.SessionManager) *IndexAdminDelivery {
	return &IndexAdminDelivery{UserService: userService, SessionManager: sessionManager}
}

func (iad *IndexAdminDelivery) Routes(adminMiddleware middlewares.AdminMiddleware) chi.Router {
	r := chi.NewRouter()
	r.With(adminMiddleware.UserAuth).Get("/", iad.HomeAdmin)
	r.Get("/signup", iad.Signup)
	r.Post("/signup", iad.AddUser)
	r.Get("/login", iad.Login)
	r.Post("/login", iad.VerifyUser)
	return r
}

func (iad *IndexAdminDelivery) Signup(w http.ResponseWriter, r *http.Request) {
	iad.Tpl = helper.TplMust(iad.Tpl.TplParseFS(templates.FS, "admin/signup.gohtml"))
	iad.Tpl.Execute(w, nil)
}

func (iad *IndexAdminDelivery) Login(w http.ResponseWriter, r *http.Request) {
	iad.Tpl = helper.TplMust(iad.Tpl.TplParseFS(templates.FS, "admin/login.gohtml"))
	iad.Tpl.Execute(w, nil)
}

func (iad *IndexAdminDelivery) AddUser(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("psw")

	// Salt and hash the password using the bcrypt algorithm
	// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 8)
	if err != nil {
		http.Error(w, "Password generated error.", http.StatusInternalServerError)
	}
	user, err := iad.UserService.AddUser(email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Added User got error.", http.StatusInternalServerError)
	}
	fmt.Fprint(w, user)

}

func (iad *IndexAdminDelivery) VerifyUser(w http.ResponseWriter, r *http.Request) {
	email := r.FormValue("email")
	password := r.FormValue("psw")

	user, err := iad.UserService.GetUser(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// prevent session fixation attacks
	if err = iad.SessionManager.RenewToken(r.Context()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	iad.SessionManager.Put(r.Context(), viper.GetString("session.auth_user_key"), user)
	http.Redirect(w, r, "/admin", http.StatusMovedPermanently)

}

func (iad *IndexAdminDelivery) HomeAdmin(w http.ResponseWriter, r *http.Request) {
	iad.Tpl = helper.TplMust(iad.Tpl.TplParseFS(templates.FS, "admin/home.gohtml"))
	iad.Tpl.Execute(w, nil)
}
