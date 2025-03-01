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

// NewAcademicYear creates a new AcademicYear repository
func (r *Repository) NewAcademicYear() AcademicYear {
	return &academicYearRepo{db: r.db}
}

// NewEvent creates a new Event repository
func (r *Repository) NewEvent() Event {
	return &eventRepo{db: r.db, year: r.NewAcademicYear()}
}

// NewBoard creates a new Board repository
func (r *Repository) NewBoard() Board {
	return &boardRepo{db: r.db, member: r.NewMember(), year: r.NewAcademicYear()}
}

// NewMember creates a new Member repository
func (r *Repository) NewMember() Member {
	return &memberRepo{db: r.db}
}
