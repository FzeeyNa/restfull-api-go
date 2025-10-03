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

// get next id for author
func getNextAuthorID() (int, error) {
	authorCollection := database.MongoDB.Collection("authors")
	var lastAuthor model.Author
	opts := options.FindOne().SetSort(bson.D{{"_id", -1}})
	err := authorCollection.FindOne(context.Background(), bson.D{}, opts).Decode(&lastAuthor)
	if err != nil {
		// If no documents are found, start with ID 1
		if err.Error() == "mongo: no documents in result" {
			return 1, nil
		}
		return 0, err
	}
	return lastAuthor.ID + 1, nil
}

func AuthorControllerGetAll(c *fiber.Ctx) error {
	authorCollection := database.MongoDB.Collection("authors")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// find all authors that are not soft deleted
	cursor, err := authorCollection.Find(ctx, bson.M{"deleted_at": nil})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	defer cursor.Close(ctx)

	var authors []model.Author
	if err := cursor.All(ctx, &authors); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(authors)
}

func AuthorControllerGetById(c *fiber.Ctx) error {
	authorCollection := database.MongoDB.Collection("authors")
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var author model.Author
	// find author by id and not soft deleted
	err = authorCollection.FindOne(ctx, bson.M{"_id": id, "deleted_at": nil}).Decode(&author)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Author not found"})
	}

	return c.JSON(author)
}

func AuthorControllerPost(c *fiber.Ctx) error {
	authorCollection := database.MongoDB.Collection("authors")
	var author model.Author
	if err := c.BodyParser(&author); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	nextID, err := getNextAuthorID()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate author ID"})
	}

	author.ID = nextID
	author.CreatedAt = time.Now()
	author.UpdatedAt = time.Now()
	author.DeletedAt = nil

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = authorCollection.InsertOne(ctx, author)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(author)
}

func AuthorControllerPut(c *fiber.Ctx) error {
	authorCollection := database.MongoDB.Collection("authors")
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	var author model.Author
	if err := c.BodyParser(&author); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	update := bson.M{
		"$set": bson.M{
			"name":       author.Name,
			"updated_at": time.Now(),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := authorCollection.UpdateOne(ctx, bson.M{"_id": id, "deleted_at": nil}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Author not found"})
	}

	// To return the updated author, we need to fetch it again
	var updatedAuthor model.Author
	err = authorCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&updatedAuthor)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Author not found after update"})
	}

	return c.JSON(updatedAuthor)
}

func AuthorControllerDelete(c *fiber.Ctx) error {
	authorCollection := database.MongoDB.Collection("authors")
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid ID"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// soft delete
	update := bson.M{"$set": bson.M{"deleted_at": time.Now()}}
	result, err := authorCollection.UpdateOne(ctx, bson.M{"_id": id, "deleted_at": nil}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if result.MatchedCount == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Author not found"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
