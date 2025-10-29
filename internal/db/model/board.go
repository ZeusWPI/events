package model

import "github.com/ZeusWPI/events/internal/db/sqlc"

type Board struct {
	ID          int
	MemberID    int
	YearID      int
	Role        string
	IsOrganizer bool
	Mattermost  string
	// Non db fields
	Member Member
	Year   Year
}

func BoardModel(board sqlc.Board, member sqlc.Member, year sqlc.Year) *Board {
	mattermost := ""
	if board.Mattermost.Valid {
		mattermost = board.Mattermost.String
	}

	return &Board{
		ID:          int(board.ID),
		MemberID:    int(board.MemberID),
		YearID:      int(board.YearID),
		Role:        board.Role,
		IsOrganizer: board.IsOrganizer,
		Mattermost:  mattermost,
		Member:      *MemberModel(member),
		Year:        *YearModel(year),
	}
}

func (b *Board) Equal(b2 Board) bool {
	return b.Role == b2.Role && b.IsOrganizer == b2.IsOrganizer && b.Mattermost == b2.Mattermost && b.Member.Equal(b2.Member) && b.Year.Equal(b2.Year)
}

// EqualEntry return true if the both board instances refer to the same entry on the website
func (b *Board) EqualEntry(b2 Board) bool {
	return b.Member.Equal(b2.Member) && b.Year.Equal(b2.Year)
}
