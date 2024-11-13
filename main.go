package main

import (
	"fmt"
	home "k-space-go/handlers"
	utils "k-space-go/utils"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	utils.Db()

	http.HandleFunc("/", home.HelloHandler)
	http.HandleFunc("/my-first-endpoint", home.MyFirstEndpoint)

	log.Println("Listening on port", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
