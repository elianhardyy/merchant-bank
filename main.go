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

	"strings"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	routes.InitRoute()
	fmt.Println(constant.URL)
	parsedURL, err := url.Parse(constant.URL)
	if err != nil {
		panic(err)
	}

	host := parsedURL.Host
	if strings.Contains(host, ":") {
		_, port, _ := strings.Cut(host, ":")
		if port == "" {
			port = "80"
		}
		serverAddr := ":" + port
		fmt.Println("Server address:", serverAddr)
		runServer(serverAddr)
	} else {
		fmt.Println("No explicit port found, defaulting to 80")
		runServer(":80")
	}
}

func runServer(serverAddr string) {
	server := &http.Server{
		Addr:    serverAddr,
		Handler: routes.R,
	}

	go func() {
		log.Println("Server is running on", serverAddr)
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
