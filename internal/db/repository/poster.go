package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/sqlc"
)

type Poster struct {
	repo Repository
}

func (r *Repository) NewPoster() *Poster {
	return &Poster{
		repo: *r,
	}
}

func (p *Poster) Get(ctx context.Context, posterID int) (*model.Poster, error) {
	poster, err := p.repo.queries(ctx).PosterGet(ctx, int32(posterID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get poster %d | %w", posterID, err)
	}

	return model.PosterModel(poster), nil
}

func (p *Poster) Create(ctx context.Context, poster *model.Poster) error {
	id, err := p.repo.queries(ctx).PosterCreate(ctx, sqlc.PosterCreateParams{
		EventID: int32(poster.EventID),
		FileID:  poster.FileID,
		Scc:     poster.SCC,
	})
	if err != nil {
		return fmt.Errorf("create poster %+v | %w", *poster, err)
	}

	poster.ID = int(id)

	return nil
}

func (p *Poster) Update(ctx context.Context, poster model.Poster) error {
	if err := p.repo.queries(ctx).PosterUpdate(ctx, sqlc.PosterUpdateParams{
		ID:      int32(poster.ID),
		EventID: int32(poster.EventID),
		FileID:  poster.FileID,
		Scc:     poster.SCC,
	}); err != nil {
		return fmt.Errorf("update poster %+v | %w", poster, err)
	}

	return nil
}

func (p *Poster) Delete(ctx context.Context, posterID int) error {
	if err := p.repo.queries(ctx).PosterDelete(ctx, int32(posterID)); err != nil {
		return fmt.Errorf("delete poster %d | %w", posterID, err)
	}

	return nil
}
