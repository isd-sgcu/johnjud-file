package image

import (
	"github.com/google/uuid"
	"github.com/isd-sgcu/johnjud-file/src/app/model"
	"github.com/isd-sgcu/johnjud-file/src/app/model/pet"
)

type Image struct {
	model.Base
	PetID    *uuid.UUID `json:"pet_id" gorm:"index:idx_name,unique"`
	Pet      *pet.Pet   `json:"pet" gorm:"foreignKey:PetID;constraint:OnUpdate:CASCADE;OnDelete:SET NULL;"`
	ImageUrl string     `json:"image_url" gorm:"mediumtext"`
}
