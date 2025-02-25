package model

// Board represents a Zeus WPI board member
type Board struct {
	ID     int
	Role   string
	Member Member
	Year   Year
}

// Equal returns true if 2 boards are equal
func (b *Board) Equal(b2 Board) bool {
	return b.Member.Equal(b2.Member) && b.Year.Equal(b2.Year)
}
