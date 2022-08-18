package main

import (
	"context"
	"fmt"
	"interview-telkom-6/handler"
	"interview-telkom-6/repository/persistence"
	"interview-telkom-6/service"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	os.Setenv("TZ", "Asia/Jakarta")
}

func main() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	rGroup := r.Group("/api")

	gin.SetMode(gin.ReleaseMode)

	if err := godotenv.Load(); err != nil {
		log.Printf(".env not found")
	}

	dbUser := os.Getenv("DATABASE_USER")
	dbPass := os.Getenv("DATABASE_PASSWORD")
	dbPort := os.Getenv("DATABASE_PORT")
	dbHost := os.Getenv("DATABASE_HOST")
	dbName := os.Getenv("DATABASE_NAME")
	ctx := context.Background()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", dbHost, dbUser,
		dbPass,
		dbName, dbPort,
	)
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("error connect to db: %v", err)
	}

	defer db.Close()

	ctx = context.WithValue(ctx, "db", db)

	productRepo := persistence.NewProductRepository(db)
	cartRepo := persistence.NewCartRepository(db)
	cartProductRepo := persistence.NewCartProductRepository(db)
	productSvc := service.NewProductService(productRepo)
	cartSvc := service.NewCartService(ctx, cartRepo, cartProductRepo, productRepo)

	handler.NewProductHandler(rGroup, productSvc)
	handler.NewCartHandler(rGroup, cartSvc)

	log.Fatal(r.Run(":" + os.Getenv("APP_PORT")))

}
