package model

// Member represents a Zeus WPI member
type Member struct {
	ID       int
	ZauthID  int
	Name     string
	Username string
}

// Equal returns true if 2 members are equal
func (m *Member) Equal(m2 Member) bool {
	return m.Name == m2.Name
}
