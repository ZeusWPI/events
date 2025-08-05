package service

import (
	"context"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/internal/db/repository"
	"github.com/ZeusWPI/events/internal/poster"
	"github.com/ZeusWPI/events/internal/server/dto"
	"github.com/ZeusWPI/events/internal/task"
	"github.com/ZeusWPI/events/pkg/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Poster struct {
	service Service

	poster repository.Poster
}

func (s *Service) NewPoster() *Poster {
	return &Poster{
		service: *s,
		poster:  *s.repo.NewPoster(),
	}
}

func (p *Poster) GetFile(ctx context.Context, posterID int) ([]byte, error) {
	poster, err := p.poster.Get(ctx, posterID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}
	if poster == nil {
		return nil, fiber.ErrBadRequest
	}

	file, err := storage.S.Get(poster.FileID)
	if err != nil {
		zap.S().Error(err)
		return nil, fiber.ErrInternalServerError
	}

	return file, nil
}

func (p *Poster) Save(ctx context.Context, posterSave dto.PosterSave) (dto.Poster, error) {
	poster := model.Poster{
		ID:      posterSave.ID,
		EventID: posterSave.EventID,
		SCC:     posterSave.SCC,
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

	poster.FileID = uuid.NewString()
	if err := storage.S.Set(poster.FileID, posterSave.File, 0); err != nil {
		zap.S().Error(err)
		return dto.Poster{}, fiber.ErrInternalServerError
	}

	var err error
	if poster.ID == 0 {
		err = p.poster.Create(ctx, &poster)
	} else {
		err = p.poster.Update(ctx, poster)
	}
	if err != nil {
		zap.S().Error(err)
		return dto.Poster{}, fiber.ErrInternalServerError
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

	if err := p.poster.Delete(ctx, posterID); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	if err := storage.S.Delete(poster.FileID); err != nil {
		zap.S().Error(err) // Only log error, it's fine
	}

	return nil
}

func (p *Poster) Sync() error {
	// The task manager runs everything in the background
	// The returned error is the status for adding it to the task manager
	// The result of the task itself if logged by the task manager
	if err := p.service.task.AddOnce(task.NewTask(poster.SyncTask, task.Now, p.service.poster.Sync)); err != nil {
		zap.S().Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}
