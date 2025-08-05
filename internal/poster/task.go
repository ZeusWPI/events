package poster

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"slices"
	"time"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/pkg/gitmate"
	"github.com/ZeusWPI/events/pkg/storage"
	"github.com/ZeusWPI/events/pkg/utils"
	"github.com/google/uuid"
)

func (c *Client) Sync(ctx context.Context) error {
	dbPosters, err := c.poster.GetAll(ctx)
	if err != nil {
		return err
	}

	gitmatePosters, err := c.fetch(ctx)
	if err != nil {
		return err
	}

	pullRequests, err := c.gitmate.Pulls(ctx)
	if err != nil {
		return err
	}

	for _, poster := range dbPosters {
		if exists := slices.ContainsFunc(gitmatePosters, func(p model.Poster) bool { return p.EqualEntry(*poster) }); exists {
			// Both the db and gitmate have a poster entry
			// Let's compare the contents of the posters and use gitmate as the source of thruth
			if err := c.checkContent(ctx, *poster); err != nil {
				return err
			}

			continue
		}

		// Poster not yet in the visueel repo
		// Add it if the event is finished
		event, err := c.event.GetByIDPopulated(ctx, poster.EventID)
		if err != nil {
			return err
		}
		if event.EndTime.After(time.Now()) {
			continue
		}

		// Is there already a pull request open for it?
		pr := toPull(*poster, *event)
		if exists := slices.ContainsFunc(pullRequests, func(p gitmate.Pull) bool { return p.Title == pr.Title && p.Body == pr.Body }); exists {
			continue
		}

		// Not in the repo yet and no open pull request to add it
		// Let's create a pr for it
		if err := c.toGitmate(ctx, *poster, *event); err != nil {
			return err
		}
	}

	for _, poster := range gitmatePosters {
		if exists := slices.ContainsFunc(utils.SliceDereference(dbPosters), func(p model.Poster) bool { return p.EqualEntry(poster) }); exists {
			continue
		}

		// Poster is not yet in our db
		// Do we have the associated event?
		event, err := c.event.GetByIDPopulated(ctx, poster.EventID)
		if err != nil {
			return err
		}
		if event == nil {
			continue
		}

		// We have the event but the poster is not in our database yet
		// Let's add it
		if err := c.toDB(ctx, poster, *event); err != nil {
			return err
		}
	}

	return nil
}

// checkContent compares the content of the db poster with the gitmate poster
func (c *Client) checkContent(ctx context.Context, poster model.Poster) error {
	dbPoster, err := storage.S.Get(poster.FileID)
	if err != nil {
		return fmt.Errorf("retrieve file for poster %+v | %w", poster, err)
	}

	event, err := c.event.GetByIDPopulated(ctx, poster.EventID)
	if err != nil {
		return err
	}

	gitmatePoster, err := c.fetchPoster(ctx, poster, *event)
	if err != nil {
		return err
	}

	if bytes.Equal(dbPoster, gitmatePoster) {
		// Same posters
		return nil
	}

	// Different posters, gitmate is the source of thruth
	// Delete the database poster
	_ = storage.S.Delete(poster.FileID) // Not a big deal if it fails
	if err := c.poster.Delete(ctx, poster.ID); err != nil {
		return err
	}

	// Get the poster from gitmate
	if err := c.toDB(ctx, poster, *event); err != nil {
		return err
	}

	return nil
}

// toGitmate will create a pull request to add the poster to the visueel repository
func (c *Client) toGitmate(ctx context.Context, poster model.Poster, event model.Event) error {
	branchName := toBranch(poster, event)

	content, err := storage.S.Get(poster.FileID)
	if err != nil {
		return fmt.Errorf("retrieve file for poster %+v | %+v | %w", poster, event, err)
	}

	if err := c.gitmate.BranchCreate(ctx, gitmate.BranchCreate{
		Name: branchName,
		Base: branchMain,
	}); err != nil {
		return err
	}

	posterMsg := "poster"
	if poster.SCC {
		posterMsg = "scc poster"
	}

	if err := c.gitmate.FileCreate(ctx, toPath(poster, event), gitmate.FileCreate{
		Branch:  branchName,
		Content: base64.StdEncoding.EncodeToString(content),
		Message: fmt.Sprintf("feat: add %s for %s", posterMsg, event.Name),
	}); err != nil {
		return err
	}

	pull := toPull(poster, event)

	if err := c.gitmate.PullCreate(ctx, gitmate.PullCreate{
		Title: pull.Title,
		Body:  pull.Body,
		Base:  branchMain,
		Head:  branchName,
	}); err != nil {
		return err
	}

	return nil
}

// toDB will fetch a poster from gitmate and insert it in the db
func (c *Client) toDB(ctx context.Context, poster model.Poster, event model.Event) error {
	bytes, err := c.fetchPoster(ctx, poster, event)
	if err != nil {
		return err
	}

	poster.FileID = uuid.NewString()
	if err := storage.S.Set(poster.FileID, bytes, 0); err != nil {
		return fmt.Errorf("unable to store new poster %+v | %w", poster, err)
	}

	if err := c.poster.Create(ctx, &poster); err != nil {
		_ = storage.S.Delete(poster.FileID) // Not a big deal if this fails
		return err
	}

	return nil
}
