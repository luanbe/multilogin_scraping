package middlewares

import (
	"github.com/alexedwards/scs/v2"
	"github.com/spf13/viper"
	"net/http"
)

type AdminMiddleware interface {
	UserAuth(next http.Handler) http.Handler
}

type AdminMiddlewareImpl struct {
	SessionManager *scs.SessionManager
}

func NewAdminMiddleware(sessionManager *scs.SessionManager) AdminMiddleware {
	return &AdminMiddlewareImpl{SessionManager: sessionManager}
}

func (ami *AdminMiddlewareImpl) UserAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ok := ami.SessionManager.Exists(r.Context(), viper.GetString("session.auth_user_key")); !ok {
			http.Redirect(w, r, "/admin/login", http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
}
