package model

import "github.com/isd-sgcu/johnjud-file/constant"

type Pet struct {
	Base
	Type         string          `json:"type" gorm:"tinytext"`
	Species      string          `json:"species" gorm:"tinytext"`
	Name         string          `json:"name" gorm:"tinytext"`
	Birthdate    string          `json:"birthdate" gorm:"tinytext"`
	Gender       constant.Gender `json:"gender" gorm:"tinytext" example:"male"`
	Color        string          `json:"color" gorm:"tinytext"`
	Pattern      string          `json:"pattern" gorm:"tinytext"`
	Habit        string          `json:"habit" gorm:"mediumtext"`
	Caption      string          `json:"caption" gorm:"mediumtext"`
	Status       constant.Status `json:"status" gorm:"mediumtext" example:"findhome"`
	IsSterile    bool            `json:"is_sterile"`
	IsVaccinated bool            `json:"is_vaccine"`
	IsVisible    bool            `json:"is_visible"`
	IsClubPet    bool            `json:"is_club_pet"`
	Origin       string          `json:"origin" gorm:"tinytext"`
	Address      string          `json:"address" gorm:"tinytext"`
	Contact      string          `json:"contact" gorm:"tinytext"`
}
