package image

import (
	"github.com/isd-sgcu/johnjud-file/internal/model"
)

type Repository interface {
	FindAll(*[]*model.Image) error
	FindOne(string, *model.Image) error
	FindByPetId(string, *[]*model.Image) error
	Create(*model.Image) error
	Update(string, *model.Image) error
	Delete(string) error
	DeleteMany([]string) error
}
