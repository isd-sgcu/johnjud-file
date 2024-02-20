package image

import (
	"github.com/isd-sgcu/johnjud-file/internal/model"
	"github.com/isd-sgcu/johnjud-file/pkg/repository/image"
	"gorm.io/gorm"
)

type repositoryImpl struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) image.Repository {
	return &repositoryImpl{db: db}
}

func (r *repositoryImpl) FindAll(result *[]*model.Image) error {
	return r.db.Model(&model.Image{}).Find(result).Error
}

func (r *repositoryImpl) FindOne(id string, result *model.Image) error {
	return r.db.Model(&model.Image{}).First(result, "id = ?", id).Error
}

func (r *repositoryImpl) FindByPetId(id string, result *[]*model.Image) error {
	return r.db.Model(&model.Image{}).Find(&result, "pet_id = ?", id).Error
}

func (r *repositoryImpl) Create(in *model.Image) error {
	return r.db.Create(&in).Error
}

func (r *repositoryImpl) Update(id string, in *model.Image) error {
	return r.db.Where(id, "id = ?", id).Updates(&in).First(&in, "id = ?", id).Error
}

func (r *repositoryImpl) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.Image{}).Error
}
func (r *repositoryImpl) DeleteMany(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return r.db.Delete(&model.Image{}, ids).Error
}
