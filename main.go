package main
import (
	"log"
	"os"

	"restfull-api-go/database"
	"restfull-api-go/route"
	"restfull-api-go/seed"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectMongo()
	seed.Seed()
	app := fiber.New()
	route.SetupRoutes(app)
	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}