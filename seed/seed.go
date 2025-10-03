package seed

import (
	"context"
	"fmt"
	"log"
	"restfull-api-go/database"
	"restfull-api-go/model"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func Seed() {
	authorCollection := database.MongoDB.Collection("authors")
	bookCollection := database.MongoDB.Collection("books")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Clear existing data
	authorCollection.DeleteMany(ctx, bson.M{})
	bookCollection.DeleteMany(ctx, bson.M{})

	now := time.Now()

	// Seed Authors
	authors := []interface{}{
		model.Author{ID: 1, Name: "J.K. Rowling", CreatedAt: now, UpdatedAt: now},
		model.Author{ID: 2, Name: "George R.R. Martin", CreatedAt: now, UpdatedAt: now},
		model.Author{ID: 3, Name: "J.R.R. Tolkien", CreatedAt: now, UpdatedAt: now},
	}

	_, err := authorCollection.InsertMany(ctx, authors)
	if err != nil {
		log.Fatal("Error seeding authors: ", err)
	}

	// Seed Books
	books := []interface{}{
		model.Book{ID: 1, Title: "Harry Potter and the Philosopher's Stone", AuthorID: 1, CreatedAt: now, UpdatedAt: now},
		model.Book{ID: 2, Title: "A Game of Thrones", AuthorID: 2, CreatedAt: now, UpdatedAt: now},
		model.Book{ID: 3, Title: "The Lord of the Rings", AuthorID: 3, CreatedAt: now, UpdatedAt: now},
	}

	_, err = bookCollection.InsertMany(ctx, books)
	if err != nil {
		log.Fatal("Error seeding books: ", err)
	}

	fmt.Println("Seeder executed successfully!")
}
