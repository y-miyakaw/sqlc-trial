package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sqlc-trial/gen/sqlc"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
)

type requestBody struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Price  string `json:"price"`
	Color  string `json:"color"`
}

type GetProductsByIDsAndColorRequest struct {
	IDs   []string `json:"ids"`
	Color string   `json:"color"`
}

func main() {
	seedItems()
	e := echo.New()
	p := e.Group("/products")
	p.GET("/:id", getProduct())
	p.GET("/", getAllProducts())
	p.POST("/search", getProductByUserIDAndNameAndColor())
	p.POST("/searchids", getProductByIDsAndName())
	p.POST("/", createProduct())
	p.PUT("/:id", updateProduct())
	p.DELETE("/:id", deleteProduct())
	e.Logger.Fatal(e.Start("localhost:8099"))
}

func dbConn() *sql.DB {
	godotenv.Load(".env.local")
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("sql.Open: %v", err)
	}
	return db
}

func seedItems() {
	db := dbConn()
	defer db.Close()
	ctx := context.Background()
	q := sqlc.New(db)
	items, err := q.GetAllProducts(ctx)
	fmt.Println(items)
	if err != nil {
		log.Printf("q.GetAllProducts: %v", err)
	}
	if len(items) == 0 {
		var sampleProducts = []sqlc.Product{
			{
				ID:     "1",
				UserID: sql.NullString{String: "1", Valid: true},
				Name:   "sample1",
				Price:  "100",
				Color:  stringToNullString("black"),
			},
			{
				ID:     "2",
				UserID: sql.NullString{String: "1", Valid: true},
				Name:   "sample2",
				Price:  "200",
				Color:  stringToNullString("red"),
			},
			{
				ID:     "3",
				UserID: sql.NullString{String: "1", Valid: true},
				Name:   "sample3",
				Price:  "300",
				Color:  stringToNullString("red"),
			},
		}

		for _, p := range sampleProducts {
			_, err := q.CreateProduct(ctx, sqlc.CreateProductParams{
				ID:     p.ID,
				UserID: p.UserID,
				Name:   p.Name,
				Price:  p.Price,
				Color:  p.Color,
			})
			if err != nil {
				log.Printf("q.CreateProduct: %v", err)
			}
		}
	}
}

func getProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		db := dbConn()
		defer db.Close()
		ctx := context.Background()
		id := c.Param("id")
		q := sqlc.New(db)
		i, err := q.GetProduct(ctx, id)
		if err != nil {
			log.Printf("q.GetProduct: %v", err)
		}

		log.Println("getProduct")
		return c.JSON(200, i)
	}
}

func getAllProducts() echo.HandlerFunc {
	return func(c echo.Context) error {
		db := dbConn()
		defer db.Close()
		ctx := context.Background()
		q := sqlc.New(db)
		items, err := q.GetAllProducts(ctx)
		if err != nil {
			log.Printf("q.GetAllProducts: %v", err)
		}
		log.Println("getAllProducts")
		return c.JSON(200, items)
	}
}

func getProductByUserIDAndNameAndColor() echo.HandlerFunc {
	return func(c echo.Context) error {
		db := dbConn()
		defer db.Close()
		ctx := context.Background()
		var requestBody requestBody
		if err := c.Bind(&requestBody); err != nil {
			log.Printf("c.Bind: %v", err)
		}
		q := sqlc.New(db)
		args := sqlc.GetProductsByUserIDAndColorParams{
			UserID:  stringToNullString(requestBody.UserID),
			Color:   stringToNullString(requestBody.Color),
			Column3: requestBody.Name,
		}
		fmt.Println(args)
		i, err := q.GetProductsByUserIDAndColor(ctx, args)
		if err != nil {
			log.Printf("q.GetProductByUserIDAndName: %v", err)
		}
		log.Println("getProductByUserIDAndName")
		return c.JSON(200, i)
	}
}

func getProductByIDsAndName() echo.HandlerFunc {
	return func(c echo.Context) error {
		db := dbConn()
		defer db.Close()
		ctx := context.Background()
		var req GetProductsByIDsAndColorRequest
		if err := c.Bind(&req); err != nil {
			log.Printf("c.Bind: %v", err)
		}
		q := sqlc.New(db)
		idsJoined := "{" + strings.Join(req.IDs, ",") + "}"
		args := sqlc.GetProductsByIDsAndColorParams{
			Column1: idsJoined,
			Color:   stringToNullString(req.Color),
		}
		fmt.Println(args)
		i, err := q.GetProductsByIDsAndColor(ctx, args)
		if err != nil {
			log.Printf("q.GetProductsByIDsAndColor: %v", err)
		}
		log.Println("GetProductsByIDsAndColor")
		return c.JSON(200, i)
	}
}

func stringToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{String: s, Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func createProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		db := dbConn()
		defer db.Close()
		var body requestBody
		if err := c.Bind(&body); err != nil {
			log.Printf("c.Bind: %v", err)
		}
		ctx := context.Background()
		q := sqlc.New(db)
		i, err := q.CreateProduct(ctx, sqlc.CreateProductParams{
			UserID: sql.NullString{String: body.UserID, Valid: true},
			Name:   body.Name,
			Price:  body.Price,
		})
		if err != nil {
			log.Printf("q.CreateProduct: %v", err)
		}
		log.Println("createProduct")
		return c.JSON(200, i)
	}
}

func updateProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		db := dbConn()
		defer db.Close()
		id := c.Param("id")
		var body requestBody
		if err := c.Bind(&body); err != nil {
			log.Printf("c.Bind: %v", err)
		}
		ctx := context.Background()
		q := sqlc.New(db)
		i, err := q.UpdateProduct(ctx, sqlc.UpdateProductParams{
			ID:     id,
			UserID: sql.NullString{String: body.UserID, Valid: true},
			Name:   body.Name,
			Price:  body.Price,
		})
		if err != nil {
			log.Printf("q.UpdateProduct: %v", err)
		}
		log.Println("updateProduct")
		return c.JSON(200, i)
	}
}

func deleteProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		db := dbConn()
		defer db.Close()
		id := c.Param("id")
		ctx := context.Background()
		q := sqlc.New(db)
		i, err := q.DeleteProduct(ctx, id)
		if err != nil {
			log.Printf("q.DeleteProduct: %v", err)
		}
		log.Println("deleteProduct")
		return c.JSON(200, i)
	}
}
