package registry

import (
	rp "github.com/luanbe/golang-web-app-structure/app/repository"
	"github.com/luanbe/golang-web-app-structure/app/service"
	"gorm.io/gorm"
)

// Todo: Add logger later
func RegisterUserService(db *gorm.DB) service.UserService {
	return service.NewUserService(
		rp.NewBaseRepository(db),
		rp.NewUserRepository(rp.NewBaseRepository(db)))
}
