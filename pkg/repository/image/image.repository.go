package image

import (
	"github.com/isd-sgcu/johnjud-file/internal/model"
)

type Repository interface {
	FindByPetId(id string, result *[]*model.Image) error
	Create(in *model.Image) error
	Update(id string, in *model.Image) error
	Delete(id string) error
}
