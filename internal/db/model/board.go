package model

import "github.com/ZeusWPI/events/internal/db/sqlc"

type Board struct {
	ID          int    `json:"id"`
	MemberID    int    `json:"member_id"`
	YearID      int    `json:"year_id"`
	Role        string `json:"role"`
	IsOrganizer bool   `json:"is_organizer"`
	// Non db fields
	Member Member `json:"member"`
	Year   Year   `json:"year"`
}

func BoardModel(board sqlc.Board) *Board {
	return &Board{
		ID:          int(board.ID),
		MemberID:    int(board.MemberID),
		YearID:      int(board.YearID),
		Role:        board.Role,
		IsOrganizer: board.IsOrganizer,
	}
}

func BoardModelPopulated(board sqlc.BoardGetAllPopulatedRow) *Board {
	username := ""
	if board.Username.Valid {
		username = board.Username.String
	}

	zauthID := 0
	if board.ZauthID.Valid {
		zauthID = int(board.ZauthID.Int32)
	}

	return &Board{
		ID:       int(board.ID),
		MemberID: int(board.ID_2),
		Member: Member{
			ID:       int(board.ID_2),
			Name:     board.Name,
			Username: username,
			ZauthID:  zauthID,
		},
		YearID: int(board.ID_3),
		Year: Year{
			ID:    int(board.ID_3),
			Start: int(board.YearStart),
			End:   int(board.YearEnd),
		},
		Role:        board.Role,
		IsOrganizer: board.IsOrganizer,
	}
}

func (b *Board) Equal(b2 Board) bool {
	return b.Role == b2.Role && b.IsOrganizer == b2.IsOrganizer && b.Member.Equal(b2.Member) && b.Year.Equal(b2.Year)
}

// EqualEntry return true if the both board instances refer to the same entry on the website
func (b *Board) EqualEntry(b2 Board) bool {
	return b.Member.Equal(b2.Member) && b.Year.Equal(b2.Year)
}
