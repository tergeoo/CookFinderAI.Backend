package main

import (
	_ "CookFinder.Backend/docs"
	"CookFinder.Backend/internal/handler"
	repository "CookFinder.Backend/internal/repo"
	"CookFinder.Backend/internal/service"
	"CookFinder.Backend/migrations"
	"CookFinder.Backend/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	dbURL := os.Getenv("DATABASE_PUBLIC_URL")
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

	ingService := service.NewIngredientService(ingRepo)
	catService := service.NewCategoryService(catRepo)
	recipeService := service.NewRecipeService(recipeRepo, ingRepo, catRepo)

	r := gin.Default()

	// Swagger docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	handler.NewIngredientHandler(r, ingService)
	handler.NewCategoryHandler(r, catService)
	handler.NewRecipeHandler(r, recipeService)
	handler.NewUploadHandler(r)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	log.Println("Server running on :8080")
	r.Run(":8080")
}
