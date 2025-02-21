package model

// Board represents a Zeus WPI board member
type Board struct {
	ID           int
	Member       Member
	AcademicYear AcademicYear
	Role         string
}

// Equal returns true if 2 boards are equal
func (b *Board) Equal(b2 Board) bool {
	return b.Member.Equal(b2.Member) && b.AcademicYear.Equal(b2.AcademicYear)
}
