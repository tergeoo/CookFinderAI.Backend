package main

import (
	_ "CookFinder.Backend/docs"
	"CookFinder.Backend/internal/handler"
	repository "CookFinder.Backend/internal/repo"
	"CookFinder.Backend/internal/service"
	"CookFinder.Backend/migrations"
	"CookFinder.Backend/pkg/db"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	slog.Info("Starting CookFinder Backend")
	dbURL := os.Getenv("DATABASE_PUBLIC_URL")

	slog.Info("Connecting to database", "url", dbURL)

	//DB, err := sqlx.Connect("postgres", "postgres://cook_finder:cook_finder@localhost:6464/cook_finder?sslmode=disable")
	DB, err := sqlx.Connect("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	err = db.Migrate(DB.DB, migrations.FS, "cook_finder", "")
	if err != nil {
		slog.Error("error migration", "err", err)
	}

	ingRepo := repository.NewIngredientRepository(DB)
	catRepo := repository.NewCategoryRepository(DB)
	recipeRepo := repository.NewRecipeRepository(DB)
	fileRepo := repository.NewFileRepository(DB)
	recipeIngredientRepo := repository.NewRecipeIngredientRepository(DB)

	ingService := service.NewIngredientService(ingRepo)
	catService := service.NewCategoryService(catRepo)
	recipeService := service.NewRecipeService(recipeRepo, recipeIngredientRepo)
	fileService := service.NewFileService(fileRepo)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Swagger docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	handler.NewIngredientHandler(r, ingService)
	handler.NewCategoryHandler(r, catService)
	handler.NewRecipeHandler(r, recipeService)
	handler.NewFileHandler(r, fileService)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Println("Server running on :8080")

	r.Run(":8080")
}
