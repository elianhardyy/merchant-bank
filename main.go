package main

import (
	"fmt"
	"go-json/constant"
	"go-json/internal/routes"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

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
	server := &http.Server{
		Addr:    serverAddr,
		Handler: routes.R,
	}
	go func() {
		log.Println("Server is running on port", port)
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Server shutting down")
}
