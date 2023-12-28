package image

import (
	"github.com/isd-sgcu/johnjud-file/internal/model"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByPetId(id string, result *model.Image) error {
	return r.db.Model(&model.Image{}).Find(&result, "pet_id = ?", id).Error
}

func (r *Repository) Create(in *model.Image) error {
	return r.db.Create(&in).Error
}

func (r *Repository) Delete(id string) error {
	return r.db.First(id).Delete(&model.Image{}).Error
}
