package entity

import (
	"github.com/google/uuid"
	"time"
)

type MovieStatus string

const (
	MovieStatusSetting    MovieStatus = "setting"
	MovieStatusProcessing MovieStatus = "processing"
	MovieStatusReady      MovieStatus = "ready"
	MovieStatusFailed     MovieStatus = "failed"
)

type Movie struct {
	ID          string
	Title       string
	Description string
	LaunchYear  int
	Genre       string
	Duration    int
	Status      MovieStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewMovie(title, description, genre string, launchYear, duration int) *Movie {
	return &Movie{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Genre:       genre,
		LaunchYear:  launchYear,
		Duration:    duration,
		Status:      MovieStatusSetting,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
