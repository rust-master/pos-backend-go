package main

import (
	"fmt"
	"log"
	"net/http"
	"pos-backend-go/router"

	"github.com/rs/cors"
)

func main() {
	r := router.Router()
	fmt.Println("Starting server on the port 8080...")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://127.0.0.1:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
		Debug:            true,
	})

	handler := c.Handler(r)

	log.Fatal(http.ListenAndServe(":8080", handler))
}
