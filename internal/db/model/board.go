package model

import "github.com/ZeusWPI/events/internal/db/sqlc"

type Board struct {
	ID       int    `json:"id"`
	MemberID int    `json:"member_id"`
	YearID   int    `json:"year_id"`
	Role     string `json:"role"`
	// Non db fields
	Member Member `json:"member"`
	Year   Year   `json:"year"`
}

func BoardModel(board sqlc.Board) *Board {
	return &Board{
		ID:       int(board.ID),
		MemberID: int(board.MemberID),
		YearID:   int(board.YearID),
		Role:     board.Role,
	}
}

func (b *Board) Equal(b2 Board) bool {
	return b.Member.Equal(b2.Member) && b.Year.Equal(b2.Year)
}
