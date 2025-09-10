package poster

import (
	"context"
	"fmt"

	"github.com/ZeusWPI/events/internal/db/model"
	"github.com/ZeusWPI/events/pkg/gitmate"
	"github.com/ZeusWPI/events/pkg/utils"
)

// The visueel repo follows the following structure
// | year
//   | event name
//     | poster.png
//     | scc.png

// fetch returns the list of posters found in gitmate
func (c *Client) fetch(ctx context.Context) ([]model.Poster, error) {
	// The first step is to fetch the data for every year
	// Unfortunately we don't know if a year directory exists in the visueel repo without making an api call for every year
	years, err := c.year.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var posters []model.Poster
	for _, year := range years {
		yearPoster, err := c.fetchYear(ctx, *year)
		if err != nil {
			return nil, err
		}

		posters = append(posters, yearPoster...)
	}

	return posters, nil
}

func (c *Client) fetchYear(ctx context.Context, year model.Year) ([]model.Poster, error) {
	// For each year we fetch all files and directories
	// Remember the file structure, we're only interested in directories as they represent events
	events, err := c.event.GetByYear(ctx, year.ID)
	if err != nil {
		return nil, err
	}
	if events == nil {
		return nil, nil
	}

	path := fmt.Sprintf("%d-%d", year.Start, year.End)
	files, err := c.gitmate.Files(ctx, path)
	if err != nil {
		return nil, err
	}

	var posters []model.Poster

	for _, file := range files {
		if file.Type != gitmate.TypeDir {
			continue
		}

		// It is a directory
		// Does the directory represent an event?
		// Is it a event directory
		event, ok := utils.SliceFind(events, func(e *model.Event) bool { return e.Name == file.Name })
		if !ok {
			continue
		}

		// The directory represents an event
		// Let's get the possible poster files inside
		eventPosters, err := c.fetchEvent(ctx, *event, file)
		if err != nil {
			return nil, err
		}

		posters = append(posters, eventPosters...)
	}

	return posters, nil
}

func (c *Client) fetchEvent(ctx context.Context, event model.Event, file gitmate.File) ([]model.Poster, error) {
	// The last step is to look for the poster files inside the event directory
	// We're looking for poster.png and scc.png
	files, err := c.gitmate.Files(ctx, file.Path)
	if err != nil {
		return nil, err
	}

	var posters []model.Poster

	for _, file := range files {
		if file.Type != gitmate.TypeFile {
			continue
		}

		if file.Name == string(posterBig) {
			// Big poster found
			posters = append(posters, model.Poster{
				EventID: event.ID,
				SCC:     false,
			})
		}

		if file.Name == string(posterScc) {
			// Scc poster found
			posters = append(posters, model.Poster{
				EventID: event.ID,
				SCC:     true,
			})
		}
	}

	return posters, nil
}

// fetchPoster fetches a single poster from gitmate
func (c *Client) fetchPoster(ctx context.Context, poster model.Poster, event model.Event) ([]byte, error) {
	path := toPath(poster, event)
	bytes, err := c.gitmate.File(ctx, path)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
