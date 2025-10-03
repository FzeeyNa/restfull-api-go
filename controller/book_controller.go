package controller

import (
	"context"
	"restfull-api-go/database"
	"restfull-api-go/model"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// get next id for book
func getNextBookID() (int, error) {
	bookCollection := database.MongoDB.Collection("books")
	var lastBook model.Book
	opts := options.FindOne().SetSort(bson.D{{"_id", -1}})
	err := bookCollection.FindOne(context.Background(), bson.D{}, opts).Decode(&lastBook)
	if err != nil {
		// If no documents are found, start with ID 1
		if err.Error() == "mongo: no documents in result" {
			return 1, nil
		}
		return 0, err
	}
	return lastBook.ID + 1, nil
}

func BookControllerGetAll(c *fiber.Ctx) error {
	bookCollection := database.MongoDB.Collection("books")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := bookCollection.Find(ctx, bson.M{"deleted_at": nil})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer cursor.Close(ctx)

	var books []model.Book
	if err := cursor.All(ctx, &books); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(books)
}

func BookControllerGetById(c *fiber.Ctx) error {
	bookCollection := database.MongoDB.Collection("books")
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var book model.Book
	err = bookCollection.FindOne(ctx, bson.M{"_id": id, "deleted_at": nil}).Decode(&book)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Book not found"})
	}

	return c.JSON(book)
}

func BookControllerPost(c *fiber.Ctx) error {
	bookCollection := database.MongoDB.Collection("books")
	var book model.Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	nextID, err := getNextBookID()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate book ID"})
	}

	book.ID = nextID
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()
	book.DeletedAt = nil

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = bookCollection.InsertOne(ctx, book)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(book)
}

func BookControllerPut(c *fiber.Ctx) error {
	bookCollection := database.MongoDB.Collection("books")
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var book model.Book
	if err := c.BodyParser(&book); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	update := bson.M{
		"$set": bson.M{
			"title":      book.Title,
			"author_id":  book.AuthorID,
			"updated_at": time.Now(),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := bookCollection.UpdateOne(ctx, bson.M{"_id": id, "deleted_at": nil}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Book not found"})
	}

	// To return the updated book, we need to fetch it again
	var updatedBook model.Book
	err = bookCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&updatedBook)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Book not found after update"})
	}

	return c.JSON(updatedBook)
}

func BookControllerDelete(c *fiber.Ctx) error {
	bookCollection := database.MongoDB.Collection("books")
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// soft delete
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}
	result, err := bookCollection.UpdateOne(ctx, bson.M{"_id": id, "deleted_at": nil}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Book not found"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
