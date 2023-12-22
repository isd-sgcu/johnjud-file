package image

import (
	"github.com/isd-sgcu/johnjud-file/src/app/model/image"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByPetId(id string, result *image.Image) error {
	return r.db.Model(&image.Image{}).Find(&result, "pet_id = ?", id).Error
}

func (r *Repository) Create(in *image.Image) error {
	return r.db.Create(&in).Error
}

func (r *Repository) Delete(id string) error {
	return r.db.First(id).Delete(&image.Image{}).Error
}
