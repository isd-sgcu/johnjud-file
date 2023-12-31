package image

import (
	"github.com/isd-sgcu/johnjud-file/internal/model"
	"github.com/isd-sgcu/johnjud-file/internal/repository/image"
	"gorm.io/gorm"
)

type Repository interface {
	FindByPetId(id string, result *[]*model.Image) error
	Create(in *model.Image) error
	Update(id string, in *model.Image) error
	Delete(id string) error
}

func NewRepository(db *gorm.DB) Repository {
	return image.NewRepository(db)
}
