package main

import (
	"be/product"
	"be/purchase"
	"be/sqs"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var c *Client

type Client struct {
	db  *sql.DB
	sqs *sqs.SQSService
}

func Init() *Client {
	c = &Client{}

	db, err := sql.Open("mysql", "talha:skymode123@tcp(localhost:3306)/estore")
	if err != nil {
		log.Fatal("Error opening database connection:", err)
	}

	c.db = db
	c.sqs = sqs.NewSQSService("estore")
	return c
}

func main() {
	c := Init()

	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	productService := product.Init()
	http.HandleFunc("/", productService.GetAllProducts)
	http.HandleFunc("/admin", productService.GetAllProducts)
	http.HandleFunc("/admin/create", productService.CreateProduct)
	http.HandleFunc("/admin/delete", productService.DeleteProduct)
	http.HandleFunc("/admin/update", productService.UpdateProduct)

	purchaseService := purchase.Init()
	http.HandleFunc("/order", purchaseService.PlaceOrder)

	go func() {
		log.Println("Server is listening on port 8080")
		log.Fatal(http.ListenAndServe(":8080", corsMiddleware(http.DefaultServeMux)))
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	c.db.Close()
}
