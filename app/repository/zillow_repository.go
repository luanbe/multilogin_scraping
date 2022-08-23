package repository

type ZillowRepository interface {
}

type ZillowRepositoryImpl struct {
	base BaseRepository
}

func NewZillowRepository(br BaseRepository) ZillowRepository {
	return &ZillowRepositoryImpl{br}
}

//func (r *ZillowRepositoryImpl) AddZillow(Zillow *entity.ZillowData) error {
//	if err := r.base.GetDB().Create(User).Error; err != nil {
//		return err
//	}
//	return nil
//}
