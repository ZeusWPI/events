package repository

import (
	"context"
	"fmt"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/sqlc"
	"github.com/ZeusWPI/events/pkg/util"
	"github.com/jackc/pgx/v5/pgtype"
)

// Member provides all model.member related database operations
type Member interface {
	GetAll(context.Context) ([]*model.Member, error)
	Save(context.Context, *model.Member) error
}

type memberRepo struct {
	repo Repository
}

// Interface compliance
var _ Member = (*memberRepo)(nil)

// GetAll returns all members
func (r *memberRepo) GetAll(ctx context.Context) ([]*model.Member, error) {
	members, err := r.repo.queries(ctx).MemberGetAll(ctx)
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
func (r *memberRepo) Save(ctx context.Context, m *model.Member) error {
	var id int32
	var err error

	if m.ID == 0 {
		// Create
		username := pgtype.Text{String: m.Username, Valid: true}
		if m.Username == "" {
			username.Valid = false
		}

		id, err = r.repo.queries(ctx).MemberCreate(ctx, sqlc.MemberCreateParams{
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

		err = r.repo.queries(ctx).MemberUpdate(ctx, sqlc.MemberUpdateParams{
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
