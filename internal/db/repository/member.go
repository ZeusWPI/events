package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/sqlc"
	"github.com/ZeusWPI/events/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
)

type Member struct {
	repo Repository
}

func (r *Repository) NewMember() *Member {
	return &Member{
		repo: *r,
	}
}

func (m *Member) GetAll(ctx context.Context) ([]*model.Member, error) {
	members, err := m.repo.queries(ctx).MemberGetAll(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get all members %w", err)
	}

	return utils.SliceMap(members, model.MemberModel), nil
}

func (m *Member) GetByID(ctx context.Context, id int) (*model.Member, error) {
	member, err := m.repo.queries(ctx).MemberGetByID(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get member by ID %d | %w", id, err)
	}

	return model.MemberModel(member), nil
}

func (m *Member) GetByName(ctx context.Context, name string) (*model.Member, error) {
	member, err := m.repo.queries(ctx).MemberGetByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get member by name %s | %w", name, err)
	}

	return model.MemberModel(member), nil
}

func (m *Member) Create(ctx context.Context, member *model.Member) error {
	zauthID := pgtype.Int4{Int32: int32(member.ZauthID), Valid: member.ZauthID != 0}
	username := pgtype.Text{String: member.Username, Valid: member.Username != ""}

	id, err := m.repo.queries(ctx).MemberCreate(ctx, sqlc.MemberCreateParams{
		Name:     member.Name,
		Username: username,
		ZauthID:  zauthID,
	})
	if err != nil {
		return fmt.Errorf("create member %+v | %w", *member, err)
	}

	member.ID = int(id)

	return nil
}

func (m *Member) Update(ctx context.Context, member model.Member) error {
	zauthID := pgtype.Int4{Int32: int32(member.ZauthID), Valid: member.ZauthID != 0}
	username := pgtype.Text{String: member.Username, Valid: member.Username != ""}

	if err := m.repo.queries(ctx).MemberUpdate(ctx, sqlc.MemberUpdateParams{
		ID:       int32(member.ID),
		ZauthID:  zauthID,
		Name:     member.Name,
		Username: username,
	}); err != nil {
		return fmt.Errorf("update member %+v | %w", member, err)
	}

	return nil
}

func (m *Member) Delete(ctx context.Context, memberID int) error {
	if err := m.repo.queries(ctx).MemberDelete(ctx, int32(memberID)); err != nil {
		return fmt.Errorf("delete member %d | %w", memberID, err)
	}

	return nil
}
