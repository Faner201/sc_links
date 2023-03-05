package db

import "time"

type MgoShortening struct {
	Identifier  string    `bson:"_id"`
	CreatedBy   string    `bson:"created_by"`
	OriginalURL string    `bson:"original_url"`
	Visits      uint64    `bson:"visits"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}
