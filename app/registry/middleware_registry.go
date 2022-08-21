package registry

import (
	"github.com/alexedwards/scs/v2"
	"github.com/luanbe/golang-web-app-structure/app/middlewares"
)

func RegisterAdminMiddleware(sessionManager *scs.SessionManager) middlewares.AdminMiddleware {
	return middlewares.NewAdminMiddleware(sessionManager)
}
