// Package repository provides all repositories
package repository

import "github.com/ZeusWPI/events/pkg/db"

// Repository is used to create specific repositories
type Repository struct {
	db db.DB
}

// New creates a new repository
func New(db db.DB) *Repository {
	return &Repository{db: db}
}

// NewEvent creates a new EventRepository
func (r *Repository) NewEvent() Event {
	return &eventRepo{db: r.db}
}
