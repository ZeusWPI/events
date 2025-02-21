package repository

import (
	"context"
	"fmt"

	"github.com/ZeusWPI/events/internal/pkg/db/sqlc"
	"github.com/ZeusWPI/events/internal/pkg/model"
	"github.com/ZeusWPI/events/pkg/db"
	"github.com/ZeusWPI/events/pkg/util"
	"github.com/jackc/pgx/v5/pgtype"
)

// Member provides all model.member related database operations
type Member interface {
	GetAll() ([]*model.Member, error)
	Save(*model.Member) error
}

type memberRepo struct {
	db db.DB
}

// Interface compliance
var _ Member = (*memberRepo)(nil)

// GetAll returns all members
func (r *memberRepo) GetAll() ([]*model.Member, error) {
	members, err := r.db.Queries().MemberGetAll(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Unable to get all members %v", err)
	}

	return util.SliceMap(members, func(m sqlc.Member) *model.Member {
		username := ""
		if m.Username.Valid {
			username = m.Username.String
		}

		return &model.Member{
			ID:       int(m.ID),
			Name:     m.Name,
			Username: username,
		}
	}), nil
}

// Save creates a new member or updates an existing one
func (r *memberRepo) Save(m *model.Member) error {
	var id int32
	var err error

	if m.ID == 0 {
		// Create
		username := pgtype.Text{String: m.Username, Valid: true}
		if m.Username == "" {
			username.Valid = false
		}

		id, err = r.db.Queries().MemberCreate(context.Background(), sqlc.MemberCreateParams{
			Name:     m.Name,
			Username: username,
		})
	} else {
		// Update
		id = int32(m.ID)
		username := pgtype.Text{String: m.Username, Valid: true}
		if m.Username == "" {
			username.Valid = false
		}

		err = r.db.Queries().MemberUpdate(context.Background(), sqlc.MemberUpdateParams{
			ID:       int32(m.ID),
			Name:     m.Name,
			Username: username,
		})
	}

	if err != nil {
		return fmt.Errorf("Unable to save member %+v | %v", *m, err)
	}

	m.ID = int(id)

	return nil
}
