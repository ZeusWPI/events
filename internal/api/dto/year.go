package dto

import "github.com/ZeusWPI/events/internal/db/model"

// Year is the data transferable object version of the model Organizer
type Year struct {
	ID        int `json:"id" validate:"required"`
	StartYear int `json:"start_year"`
	EndYear   int `json:"end_year"`
}

// YearDTO converts a model Year to a DTO Year
func YearDTO(y *model.Year) Year {
	return Year(*y)
}

// ToModel converts a DTO Year to a model Year
func (y *Year) ToModel() *model.Year {
	year := model.Year(*y)
	return &year
}
