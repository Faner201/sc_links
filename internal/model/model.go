package model

import (
	"time"
)

type Shortering struct {
	Identifier  string    `json:"identifier"`
	CreatedBy   string    `json:"created_by"`
	OriginalURL string    `json:"original_url"`
	Visits      uint64    `json:"visits"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type User struct {
	IsActive        bool      `json:"is_verified,omitempty"`
	GitHubLogin     string    `json:"git_login"`
	GitHubAccessKey string    `json:"git_access_key,omitempty"`
	CreatedAt       time.Time `json:"created_at,omitempty"`
}
