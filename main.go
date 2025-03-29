package main

import (
	"fmt"
	"go-json/constant"
	"go-json/internal/routes"
	"log"
	"net/http"
	"net/url"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading.env file")
	}
}
func main() {

	routes.InitRoute()
	parsedURL, err := url.Parse(constant.URL)
	if err != nil {
		panic(err)
	}
	port := parsedURL.Port()
	if port == "" {
		port = "80"
	}
	serverAddr := ":" + port
	fmt.Println(serverAddr)
	http.ListenAndServe(serverAddr, routes.R)
}
