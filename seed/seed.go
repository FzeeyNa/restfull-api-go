package seed

import (
	"context"
	"fmt"
	"log"
	"restfull-api-go/database"
	"restfull-api-go/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var authorCollection = database.MongoDB.Collection("authors")
var bookCollection = database.MongoDB.Collection("buku")

func Seed() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Clear existing data
	authorCollection.DeleteMany(ctx, bson.M{})
	bookCollection.DeleteMany(ctx, bson.M{})

	// Seed Authors
	author1ID := primitive.NewObjectID()
	author2ID := primitive.NewObjectID()
	author3ID := primitive.NewObjectID()

	authors := []interface{}{
		model.Author{ID: author1ID, Name: "J.K. Rowling"},
		model.Author{ID: author2ID, Name: "George R.R. Martin"},
		model.Author{ID: author3ID, Name: "J.R.R. Tolkien"},
	}

	_, err := authorCollection.InsertMany(ctx, authors)
	if err != nil {
		log.Fatal("Error seeding authors: ", err)
	}

	// Seed Books
	books := []interface{}{
		model.Book{ID: primitive.NewObjectID(), Title: "Harry Potter and the Philosopher's Stone", AuthorID: author1ID},
		model.Book{ID: primitive.NewObjectID(), Title: "A Game of Thrones", AuthorID: author2ID},
		model.Book{ID: primitive.NewObjectID(), Title: "The Lord of the Rings", AuthorID: author3ID},
	}

	_, err = bookCollection.InsertMany(ctx, books)
	if err != nil {
		log.Fatal("Error seeding books: ", err)
	}

	fmt.Println("Seeder executed successfully!")
}
