package model

import (
	"fmt"
)

// Year represents an academic year
type Year struct {
	ID        int
	StartYear int
	EndYear   int
}

// String returns the shortened string version
func (a *Year) String() string {
	return fmt.Sprintf("%02d-%02d", a.StartYear%100, a.EndYear%100)
}

// Equal returns true if 2 years are equal
func (a *Year) Equal(a2 Year) bool {
	return a.StartYear == a2.StartYear && a.EndYear == a2.EndYear
}
