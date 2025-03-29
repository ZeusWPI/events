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
	GetByID(context.Context, int) (*model.Member, error)
	GetByName(context.Context, string) (*model.Member, error)
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
		return nil, fmt.Errorf("unable to get all members %w", err)
	}

	return util.SliceMap(members, func(m sqlc.Member) *model.Member {
		zauthID := 0
		if m.ZauthID.Valid {
			zauthID = int(m.ZauthID.Int32)
		}

		username := ""
		if m.Username.Valid {
			username = m.Username.String
		}

		return &model.Member{
			ID:       int(m.ID),
			ZauthID:  zauthID,
			Name:     m.Name,
			Username: username,
		}
	}), nil
}

// GetByID returns a member given an id
func (r *memberRepo) GetByID(ctx context.Context, id int) (*model.Member, error) {
	member, err := r.repo.queries(ctx).MemberGetByID(ctx, int32(id))
	if err != nil {
		return nil, fmt.Errorf("unable to get member by ID %d | %w", id, err)
	}

	zauthID := 0
	if member.ZauthID.Valid {
		zauthID = int(member.ZauthID.Int32)
	}

	username := ""
	if member.Username.Valid {
		username = member.Username.String
	}

	return &model.Member{
		ID:       int(member.ID),
		ZauthID:  zauthID,
		Name:     member.Name,
		Username: username,
	}, nil
}

// GetByName a member given a name
func (r *memberRepo) GetByName(ctx context.Context, name string) (*model.Member, error) {
	member, err := r.repo.queries(ctx).MemberGetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("unable to get member by name %s", name)
	}

	zauthID := 0
	if member.ZauthID.Valid {
		zauthID = int(member.ZauthID.Int32)
	}

	username := ""
	if member.Username.Valid {
		username = member.Username.String
	}

	return &model.Member{
		ID:       int(member.ID),
		ZauthID:  zauthID,
		Name:     member.Name,
		Username: username,
	}, nil
}

// Save creates a new member or updates an existing one
func (r *memberRepo) Save(ctx context.Context, m *model.Member) error {
	var id int32
	var err error

	zauthID := pgtype.Int4{Int32: int32(m.ZauthID), Valid: true}
	if m.ZauthID == 0 {
		zauthID.Valid = false
	}

	username := pgtype.Text{String: m.Username, Valid: true}
	if m.Username == "" {
		username.Valid = false
	}

	if m.ID == 0 {
		// Create
		id, err = r.repo.queries(ctx).MemberCreate(ctx, sqlc.MemberCreateParams{
			Name:     m.Name,
			Username: username,
			ZauthID:  zauthID,
		})
	} else {
		// Update
		id = int32(m.ID)

		err = r.repo.queries(ctx).MemberUpdate(ctx, sqlc.MemberUpdateParams{
			ID:       int32(m.ID),
			ZauthID:  zauthID,
			Name:     m.Name,
			Username: username,
		})
	}

	if err != nil {
		return fmt.Errorf("unable to save member %+v | %w", *m, err)
	}

	m.ID = int(id)

	return nil
}
