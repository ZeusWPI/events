package dto

// Organizer is the data transferable object of the model organizer
type Organizer struct {
	ID   int    `json:"id"`
	Role string `json:"role"`
	Name string `json:"name"`
}
