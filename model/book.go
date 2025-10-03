package model

import "time"

type Book struct {
	ID        int        `bson:"_id" json:"id"`
	Title     string     `bson:"title" json:"title"`
	AuthorID  int        `bson:"author_id" json:"author_id"`
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}
