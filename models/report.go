package models

import "time"

type ReportRequest struct {
	Title       string  `json:"title" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Latitude    float64 `json:"latitude" validate:"required"`
	Longitude   float64 `json:"longitude" validate:"required"`
}

func (r *ReportRequest) ToReport() *Report {
	return &Report{
		Title:       r.Title,
		Description: r.Description,
		Latitude:    r.Latitude,
		Longitude:   r.Longitude,
	}
}

type Report struct {
	ID          string    `json:"id" gorm:"primary_key" param:"id"`
	AuthorID    string    `json:"author_id"`
	Author      *User     `json:"author,omitempty" gorm:"foreignKey:author_id"`
	Title       string    `json:"title"`
	Image       string    `json:"image"`
	Description string    `json:"description"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Popularity  int       `json:"popularity"`
	CreatedAt   time.Time `json:"created_at"`
}
