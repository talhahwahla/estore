package product

import (
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Product struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Category    string `json:"category"`
}

type ProductEvent struct {
	EventType string `json:"event_type"`
	ProductID int    `json:"product_id"`
}

func (c *Client) getLastInsertID() (int, error) {
	var lastInsertID int
	err := c.db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&lastInsertID)
	if err != nil {
		return 0, err
	}
	return lastInsertID, nil
}

func (c *Client) sendMessageToSQS(productEvent *ProductEvent) {
	productEventJSON, err := json.Marshal(productEvent)
	if err != nil {
		log.Fatal("failed to marshal product event: ", err)
	}

	c.sqs.SendMessage(string(productEventJSON))
}

func (c *Client) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request: GET /products")
	rows, err := c.db.Query("SELECT * FROM products")
	if err != nil {
		log.Fatal("Error pinging database:", err)
	}
	defer rows.Close()

	var products []Product

	for rows.Next() {
		var product Product
		if err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.Category,
		); err != nil {
			log.Fatal("Error scanning row:", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	productsJSON, err := json.Marshal(products)
	if err != nil {
		log.Fatal("Error marshalling products to JSON:", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(productsJSON)
}

func (c *Client) CreateProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request: POST /products/create")
	var product Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	_, err = c.db.Exec(`
		INSERT INTO products
		(name, description, price, category)
		VALUES (?, ?, ?, ?)`,
		product.Name,
		product.Description,
		product.Price,
		product.Category,
	)
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	pid, _ := c.getLastInsertID()
	c.sendMessageToSQS(&ProductEvent{
		EventType: "create",
		ProductID: pid,
	})
}

func (c *Client) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request: DELETE /products/delete")
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Product ID is required", http.StatusBadRequest)
		return
	}

	_, err := c.db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

	pid, _ := c.getLastInsertID()
	c.sendMessageToSQS(&ProductEvent{
		EventType: "delete",
		ProductID: pid,
	})
}

func (c *Client) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request: PUT /products/update")
	var product Product
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	_, err = c.db.Exec(`
		UPDATE products
		SET name = ?, description = ?, price = ?, category = ?
		WHERE id = ?`,
		product.Name,
		product.Description,
		product.Price,
		product.Category,
		product.Id,
	)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	pid, _ := c.getLastInsertID()
	c.sendMessageToSQS(&ProductEvent{
		EventType: "update",
		ProductID: pid,
	})
}

// func Main() {
// 	c = Init()

// corsMiddleware := func(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
// 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")

// 		if r.Method == "OPTIONS" {
// 			w.WriteHeader(http.StatusOK)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }

// http.HandleFunc("/products", c.GetAllProducts)
// http.HandleFunc("/products/create", c.CreateProduct)
// http.HandleFunc("/products/delete", c.DeleteProduct)
// http.HandleFunc("/products/update", c.UpdateProduct)

// go func() {
// 	log.Println("Server is listening on port 8080")
// 	log.Fatal(http.ListenAndServe(":8080", corsMiddleware(http.DefaultServeMux)))
// }()

// stop := make(chan os.Signal, 1)
// signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

// <-stop

// c.db.Close()
// }
