package route
import (
	"restfull-api-go/controller"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	//v1 := api.Group("/v1")

	books := api.Group("/books")
	books.Get("/", controller.BookControllerGetAll)
	books.Get("/:id", controller.BookControllerGetById)
	books.Post("/", controller.BookControllerPost)
	books.Put("/:id", controller.BookControllerPut)
	books.Delete("/:id", controller.BookControllerDelete)

	authors := api.Group("/authors")
	authors.Get("/", controller.AuthorControllerGetAll)
	authors.Get("/:id", controller.AuthorControllerGetById)
	authors.Post("/", controller.AuthorControllerPost)
	authors.Put("/:id", controller.AuthorControllerPut)
	authors.Delete("/:id", controller.AuthorControllerDelete)
}