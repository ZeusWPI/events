package service

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"math"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/poster"
	"github.com/ZeusWPI/events/internal/server/dto"
	"github.com/ZeusWPI/events/internal/task"
	"github.com/ZeusWPI/events/pkg/image"
	"github.com/ZeusWPI/events/pkg/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	a4AspectRatio          = 1.4142
	a4AspectRatioTolerance = 0.01
)

type Poster struct {
	service Service

	event  repository.Event
	poster repository.Poster
}

func (s *Service) NewPoster() *Poster {
	return &Poster{
		service: *s,
		event:   *s.repo.NewEvent(),
		poster:  *s.repo.NewPoster(),
	}
}

func (p *Poster) GetFile(ctx context.Context, posterID int, original bool) ([]byte, error) {
	poster, err := p.poster.Get(ctx, posterID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}
	if poster == nil {
		return nil, fiber.ErrBadRequest
	}

	var file []byte

	if original {
		file, err = storage.S.Get(poster.FileID)
	} else {
		file, err = storage.S.Get(poster.WebpID)
	}

	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}
	if file == nil {
		return nil, fiber.ErrBadRequest
	}

	return file, nil
}

func (p *Poster) Save(ctx context.Context, posterSave dto.PosterSave) (dto.Poster, error) {
	poster := model.Poster{
		ID:      posterSave.ID,
		EventID: posterSave.EventID,
		SCC:     posterSave.SCC,
	}

	a4, err := isA4(posterSave.File)
	if err != nil {
		zap.S().Error(err)
		return dto.Poster{}, fiber.ErrInternalServerError
	}
	if !a4 {
		return dto.Poster{}, fiber.ErrBadRequest
	}

	// Does the event exist?
	event, err := p.event.GetByID(ctx, poster.EventID)
	if err != nil {
		zap.S().Error(err)
		return dto.Poster{}, fiber.ErrInternalServerError
	}
	if event == nil {
		return dto.Poster{}, fiber.ErrBadRequest
	}

	if poster.ID != 0 {
		// Update, delete old poster
		oldPoster, err := p.poster.Get(ctx, poster.ID)
		if err != nil {
			zap.S().Error(err)
			return dto.Poster{}, fiber.ErrInternalServerError
		}
		if oldPoster == nil {
			return dto.Poster{}, fiber.ErrBadRequest
		}

		if err = storage.S.Delete(oldPoster.FileID); err != nil {
			zap.S().Error(err) // Only log error, it's fine
		}
	}

	webp, err := image.ToWebp(posterSave.File)
	if err != nil {
		zap.S().Error(err)
		return dto.Poster{}, fiber.ErrInternalServerError
	}

	poster.WebpID = uuid.NewString()
	if err := storage.S.Set(poster.WebpID, webp, 0); err != nil {
		zap.S().Error(err)
		return dto.Poster{}, fiber.ErrInternalServerError
	}

	poster.FileID = uuid.NewString()
	if err := storage.S.Set(poster.FileID, posterSave.File, 0); err != nil {
		zap.S().Error(err)
		return dto.Poster{}, fiber.ErrInternalServerError
	}

	if err := p.service.withRollback(ctx, func(ctx context.Context) error {
		if posterSave.ID == 0 {
			err = p.poster.Create(ctx, &poster)
		} else {
			err = p.poster.Update(ctx, poster)
		}
		if err != nil {
			zap.S().Error(err)
			return fiber.ErrInternalServerError
		}

		if posterSave.ID == 0 {
			if err := p.service.poster.Create(ctx, poster); err != nil {
				zap.S().Error(err)
				return fiber.ErrInternalServerError
			}
		}

		return nil
	}); err != nil {
		return dto.Poster{}, err
	}

	return dto.PosterDTO(&poster), nil
}

func (p *Poster) Delete(ctx context.Context, posterID int) error {
	poster, err := p.poster.Get(ctx, posterID)
	if err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}
	if poster == nil {
		return fiber.ErrBadRequest
	}

	return p.service.withRollback(ctx, func(ctx context.Context) error {
		if err := p.poster.Delete(ctx, posterID); err != nil {
			zap.S().Error(err)
			return fiber.ErrInternalServerError
		}

		if err := p.service.poster.Delete(ctx, *poster); err != nil {
			zap.S().Error(err)
			return fiber.ErrInternalServerError
		}

		if err := storage.S.Delete(poster.FileID); err != nil {
			zap.S().Error(err) // Only log error, it's fine
		}

		return nil
	})
}

func (p *Poster) Sync() error {
	// The task manager runs everything in the background
	// The returned error is the status for adding it to the task manager
	// The result of the task itself is logged by the task manager
	if err := task.Manager.AddOnce(task.NewTask(poster.TaskUID, "Syncronizing posters", task.Now, p.service.poster.Sync)); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}

func isA4(data []byte) (bool, error) {
	img, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return false, fmt.Errorf("decode png image %w", err)
	}

	bounds := img.Bounds()
	aspectRatio := float64(bounds.Dy()) / float64(bounds.Dx())

	return math.Abs(aspectRatio-a4AspectRatio) <= a4AspectRatioTolerance, nil
}
