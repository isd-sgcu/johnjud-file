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
	return r.db.First(id).Delete(&model.Image{}).Error
}
