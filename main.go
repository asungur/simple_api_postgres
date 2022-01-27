package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"simple_api_postgres/router"

	"github.com/joho/godotenv"
)

func getPort(portKey string) (port string) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	port = os.Getenv(portKey)
	return
}

func main() {
	r := router.Router()
	port := ":" + getPort("PORT")
	fmt.Println("Starting server on the port:", port)

	log.Fatal(http.ListenAndServe(port, r))
}
