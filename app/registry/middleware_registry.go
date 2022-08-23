package registry

import (
	"github.com/alexedwards/scs/v2"
	"multilogin_scraping/app/middlewares"
)

func RegisterAdminMiddleware(sessionManager *scs.SessionManager) middlewares.AdminMiddleware {
	return middlewares.NewAdminMiddleware(sessionManager)
}
