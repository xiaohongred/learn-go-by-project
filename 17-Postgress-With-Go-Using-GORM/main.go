package main

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"postgreee-with-go-using-gorm/models"
	"postgreee-with-go-using-gorm/storage"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(ctx *fiber.Ctx) error {
	book := Book{}
	err := ctx.BodyParser(&book)
	if err != nil {
		ctx.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}
	err = r.DB.Create(&book).Error
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "cloud not create book"})
		return err
	}
	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "book has been added",
	})
	return nil
}

func (r *Repository) DeleteBook(ctx *fiber.Ctx) error {
	bookModel := models.Books{}
	id := ctx.Params("id")
	if id == "" {
		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}
	err := r.DB.Delete(bookModel, id).Error
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete book",
		})
		return err
	}
	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "books delete successfully",
	})
	return nil
}

func (r *Repository) GetBookByID(ctx *fiber.Ctx) error {
	bookModel := &models.Books{}
	id := ctx.Params("id")
	if id == "" {
		ctx.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"message": "id cannot be empty",
		})
		return nil
	}
	err := r.DB.Where("id = ?", id).First(bookModel).Error
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not find book by id",
		})
		return err
	}
	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "successful",
		"data":    bookModel,
	})
	return nil
}

func (r *Repository) GetBooks(ctx *fiber.Ctx) error {
	bookModels := &[]models.Books{}
	err := r.DB.Find(bookModels).Error
	if err != nil {
		ctx.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not get books",
		})
		return err
	}
	ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "books fetched successfully",
		"data":    bookModels,
	})
	return nil
}

var r Repository

func (r *Repository) SetupRoute(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_books", r.CreateBook)
	api.Delete("delete_book/:id", r.DeleteBook)
	api.Get("/get_books/:id", r.GetBookByID)
	api.Get("/books", r.GetBooks)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	cfg := storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PWD"),
		User:     os.Getenv("DB_USR"),
		SSLModel: os.Getenv("DB_SSLModel"),
		DBName:   os.Getenv("DB_DBNAME"),
	}

	db, err := storage.NewConnection(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal("could not load database")
	}
	err = models.MigrateBooks(db)
	r := Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoute(app)
	app.Listen(":8080")
}
