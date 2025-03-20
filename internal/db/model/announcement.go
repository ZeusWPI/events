package model

import (
	"errors"
	"slices"
	"time"
)

type Target int

const (
	Discord Target = iota
	Mattermost
)

var TARGET_STRINGS = []string{"discord", "mattermost"} // Ohno dont forget!!!

func (t Target) String() string {
	return TARGET_STRINGS[t]
}

func FromTargetString(s string) (Target, error) {
	idx := slices.Index(TARGET_STRINGS, s)
	if idx == -1 {
		return -1, errors.New("invalid target")
	}

	return Target(idx), nil
}

// Announcement represents an announcement
type Announcement struct {
	ID        int
	Content   string
	Time      time.Time
	Target    Target
	Sent      bool
	Event     Event
	Member    Member
	CreatedAt time.Time
	UpdatedAt time.Time
}
