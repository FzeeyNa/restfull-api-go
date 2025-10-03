package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Book struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Title    string             `bson:"title" json:"title"`
	AuthorID primitive.ObjectID `bson:"author_id,omitempty" json:"author_id,omitempty"`
}
