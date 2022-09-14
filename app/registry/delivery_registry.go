package registry

import (
	"github.com/alexedwards/scs/v2"
	"multilogin_scraping/app/delivery/delivery_admin"
	"multilogin_scraping/app/service"
)

func RegisterIndexAdminDelivery(
	userService service.UserService,
	sessionManager *scs.SessionManager,
) *delivery_admin.IndexAdminDelivery {
	return delivery_admin.NewIndexAdminDelivery(userService, sessionManager)
}

func RegisterUserAdminDelivery(
	userService service.UserService,
	sessionManager *scs.SessionManager,
) *delivery_admin.UserAdminDelivery {
	return delivery_admin.NewUserAdminDelivery(userService, sessionManager)
}
