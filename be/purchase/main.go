package purchase

import (
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Product struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type Purchase struct {
	Products     []Product `json:"products"`
	CustomerInfo string    `json:"customer_info"`
}

type PurchaseEvent struct {
	EventType string `json:"event_type"`
	OrderID   string `json:"order_id"`
}

func generateOrderID() string {
	orderID, err := uuid.NewRandom()
	if err != nil {
		log.Println("Failed to generate order ID:", err)
		return ""
	}
	return orderID.String()
}

func (c *Client) sendMessageToSQS(purchaseEvent *PurchaseEvent) {
	purchaseEventJSON, err := json.Marshal(purchaseEvent)
	if err != nil {
		log.Fatal("failed to marshal purchase event: ", err)
	}

	c.sqs.SendMessage(string(purchaseEventJSON))
}

func (c *Client) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request: POST /order")
	var purchase Purchase

	err := json.NewDecoder(r.Body).Decode(&purchase)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	orderID := generateOrderID()
	if orderID == "" {
		http.Error(w, "Failed to generate order ID", http.StatusInternalServerError)
		return
	}

	for _, product := range purchase.Products {
		_, err := c.db.Exec(`
			INSERT INTO orders
			(order_id, product_id, quantity, customer_info)
			VALUES (?, ?, ?, ?)`,
			orderID,
			product.ProductID,
			product.Quantity,
			purchase.CustomerInfo,
		)

		if err != nil {
			http.Error(w, "Failed to place order", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)

	c.sendMessageToSQS(&PurchaseEvent{
		EventType: "place_order",
		OrderID:   orderID,
	})
}
