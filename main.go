package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sqlc-trial/gen/sqlc"
	"strconv"

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

type getProductsByNameOrPriceOrCompanyIDRequest struct {
	Name      string `json:"name"`
	Price     string `json:"price"`
	CompanyID string `json:"company_id"`
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
	p.POST("/search", getProductsByNameOrPriceOrCompanyID())
	// p.GET("/company/:id", getProductsAndCompanyByCompanyID())
	// p.POST("/returning", createProductWithReturning())
	// p.POST("/", createProductWithoutReturning())
	// p.PUT("/:id", updateProduct())
	// p.DELETE("/:id", deleteProduct())
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
	products, err := q.GetAllProducts(ctx)
	if err != nil {
		log.Fatalf("q.GetAllProducts: %v", err)
	}
	if len(products) == 0 {
		sampleProducts := []sqlc.CreateProductWithoutReturningParams{
			{
				Name:      sql.NullString{String: "product A", Valid: true},
				Price:     sql.NullInt32{Int32: 1000, Valid: true},
				CompanyID: sql.NullInt32{Int32: 1, Valid: true},
			},
			{
				Name:      sql.NullString{String: "product B", Valid: true},
				Price:     sql.NullInt32{Int32: 2000, Valid: true},
				CompanyID: sql.NullInt32{Int32: 1, Valid: true},
			},
			{
				Name:      sql.NullString{String: "product C", Valid: true},
				Price:     sql.NullInt32{Int32: 1500, Valid: true},
				CompanyID: sql.NullInt32{Int32: 2, Valid: true},
			},
		}

		for _, param := range sampleProducts {
			if err := q.CreateProductWithoutReturning(ctx, param); err != nil {
				log.Fatalf("q.CreateProduct: %v", err)
			}
		}
	}
	companies, err := q.GetAllCompanies(ctx)
	if err != nil {
		log.Fatalf("q.GetAllProducts: %v", err)
	}
	if len(companies) == 0 {
		sampleCompanies := []sqlc.CreateCompanyWithoutReturningParams{
			{
				Name:    "company A",
				Address: sql.NullString{String: "address A", Valid: true},
				Person:  sql.NullString{String: "person A", Valid: true},
			},
			{
				Name:    "company B",
				Address: sql.NullString{String: "address B", Valid: true},
				Person:  sql.NullString{String: "person B", Valid: true},
			},
		}
		for _, param := range sampleCompanies {
			if err := q.CreateCompanyWithoutReturning(ctx, param); err != nil {
				log.Fatalf("q.CreateCompany: %v", err)
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
		param, err := strconv.Atoi(id)
		if err != nil {
			log.Printf("strconv.Atoi: %v", err)
			return c.JSON(http.StatusBadRequest, "invalid id")
		}
		q := sqlc.New(db)
		i, err := q.GetProduct(ctx, int32(param))
		if err != nil {
			log.Printf("q.GetProduct: %v", err)
			return c.JSON(http.StatusInternalServerError, "error")
		}

		log.Println("getProduct")
		return c.JSON(http.StatusOK, i)
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
			return c.JSON(http.StatusInternalServerError, "error")
		}
		log.Println("getAllProducts")
		return c.JSON(200, items)
	}
}

func getProductsByNameOrPriceOrCompanyID() echo.HandlerFunc {
	return func(c echo.Context) error {
		db := dbConn()
		defer db.Close()
		ctx := context.Background()
		q := sqlc.New(db)
		var req getProductsByNameOrPriceOrCompanyIDRequest
		if err := c.Bind(&req); err != nil {
			log.Printf("c.Bind: %v", err)
			return c.JSON(http.StatusBadRequest, "error")
		}
		params := sqlc.GetProductsByNameOrPriceOrCompanyIDParams{
			Name:      stringToNullString(req.Name),
			Price:     stringToInt32(req.Price),
			CompanyID: stringToInt32(req.CompanyID),
		}
		fmt.Println(params)
		items, err := q.GetProductsByNameOrPriceOrCompanyID(ctx, params)
		if err != nil {
			log.Printf("q.GetAllProducts: %v", err)
			return c.JSON(http.StatusInternalServerError, "error")
		}
		log.Println("getAllProducts")
		return c.JSON(http.StatusOK, items)
	}
}

func stringToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{String: s, Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func stringToInt32(s string) sql.NullInt32 {
	if s == "" {
		return sql.NullInt32{Int32: 0, Valid: false}
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("strconv.Atoi: %v", err)
		return sql.NullInt32{Int32: 0, Valid: false}
	}
	return sql.NullInt32{Int32: int32(i), Valid: true}
}
