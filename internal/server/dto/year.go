package dto

import "github.com/ZeusWPI/events/internal/db/model"

type Year struct {
	ID    int `json:"id"`
	Start int `json:"start"`
	End   int `json:"end"`
}

func YearDTO(y *model.Year) Year {
	return Year(*y)
}

func (y *Year) ToModel() *model.Year {
	year := model.Year(*y)
	return &year
}
