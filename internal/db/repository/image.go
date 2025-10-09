package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/sqlc"
)

type Image struct {
	repo Repository
}

func (r *Repository) NewImage() *Image {
	return &Image{
		repo: *r,
	}
}

func (i *Image) Get(ctx context.Context, imageID int) (*model.Image, error) {
	image, err := i.repo.queries(ctx).ImageGet(ctx, int32(imageID))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get image by id %d | %w", imageID, err)
	}

	return model.ImageModel(image), nil
}

func (i *Image) Create(ctx context.Context, image *model.Image) error {
	id, err := i.repo.queries(ctx).ImageCreate(ctx, sqlc.ImageCreateParams{
		Name:   image.Name,
		FileID: image.FileID,
	})
	if err != nil {
		return fmt.Errorf("create image %+v | %w", *image, err)
	}

	image.ID = int(id)

	return nil
}
