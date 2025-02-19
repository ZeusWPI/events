package models

import (
	"fmt"
)

// AcademicYear represents an academic year
type AcademicYear struct {
	ID        int
	StartYear int
	EndYear   int
}

// String returns the shortened string version
func (a *AcademicYear) String() string {
	return fmt.Sprintf("%02d-%02d", a.StartYear%100, a.EndYear%100)
}

// Equal returns true if 2 academic years are equal
func (a *AcademicYear) Equal(a2 AcademicYear) bool {
	return a.StartYear == a2.StartYear && a.EndYear == a2.EndYear
}
