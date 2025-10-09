package model

import "github.com/ZeusWPI/events/internal/db/sqlc"

type Image struct {
	ID     int
	Name   string
	FileID string
}

func ImageModel(image sqlc.Image) *Image {
	return &Image{
		ID:     int(image.ID),
		Name:   image.Name,
		FileID: image.FileID,
	}
}
