package model

import "time"

type Author struct {
	ID        int        `bson:"_id" json:"id"`
	Name      string     `bson:"name" json:"name"`
	CreatedAt time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time  `bson:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `bson:"deleted_at,omitempty" json:"deleted_at,omitempty"`
}
