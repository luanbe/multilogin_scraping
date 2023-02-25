package registry

import (
	"gorm.io/gorm"
	rp "multilogin_scraping/app/repository"
	"multilogin_scraping/app/service"
)

// Todo: Add logger later
func RegisterUserService(db *gorm.DB) service.UserService {
	return service.NewUserService(
		rp.NewBaseRepository(db),
		rp.NewUserRepository(rp.NewBaseRepository(db)))
}

func RegisterMaindb3Service(db *gorm.DB) service.Maindb3Service {
	return service.NewMaindb3Service(
		rp.NewBaseRepository(db),
		rp.NewMaindb3Repository(rp.NewBaseRepository(db)))
}

func RegisterZillowService(db *gorm.DB) service.ZillowService {
	return service.NewZillowService(
		rp.NewBaseRepository(db),
		rp.NewZillowRepository(rp.NewBaseRepository(db)))
}

func RegisterRealtorService(db *gorm.DB) service.RealtorService {
	return service.NewRealtorService(
		rp.NewBaseRepository(db),
		rp.NewRealtorRepository(rp.NewBaseRepository(db)))
}

func RegisterMovotoService(db *gorm.DB) service.MovotoService {
	return service.NewMovotoService(
		rp.NewBaseRepository(db),
		rp.NewMovotoRepository(rp.NewBaseRepository(db)))
}
