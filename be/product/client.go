package product

import (
	"be/sqs"
	"database/sql"
	"log"
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
