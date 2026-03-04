package main

import (
	"blog/internal/router"
	"log"
	"net/http"
)

func main() {
	router := router.NewRouter()
	log.Println("Server is running")
	log.Fatal(http.ListenAndServe(":8080", router.SetRouter()))
}
