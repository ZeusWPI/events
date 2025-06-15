package model

import (
	"fmt"

	"github.com/ZeusWPI/events/internal/db/sqlc"
)

type Year struct {
	ID    int `json:"id"`
	Start int `json:"year_start"`
	End   int `json:"year_end"`
}

func YearModel(year sqlc.Year) *Year {
	return &Year{
		ID:    int(year.ID),
		Start: int(year.YearStart),
		End:   int(year.YearEnd),
	}
}

func (y *Year) String() string {
	return fmt.Sprintf("%02d-%02d", y.Start%100, y.End%100)
}

func (y *Year) Equal(y2 Year) bool {
	return y.Start == y2.Start && y.End == y2.End
}
