package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/oauth2/clientcredentials"
)

func main() {
	// HTTP server to expose the OpenAPI endpoint
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/greeter/greet", greetHandler) // OpenAPI endpoint

	serverPort := 8080
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", serverPort),
		Handler: serverMux,
	}

	go func() {
		log.Printf("Starting HTTP Greeter on port %d\n", serverPort)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP ListenAndServe error: %v", err)
		}
		log.Println("HTTP server stopped serving new requests.")
	}()

	// Graceful shutdown on interrupt or termination signal
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)
	<-stopCh // Wait for shutdown signal

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("Shutting down the server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Server shutdown complete.")
}

// greetHandler handles the /greeter/greet API call
func greetHandler(w http.ResponseWriter, r *http.Request) {

	// OAuth2 Client credentials
	serviceURL := os.Getenv("CHOREO_MYSERVICE_SERVICE_URL")
	clientID := os.Getenv("CHOREO_MYSERVICE_CONSUMER_KEY")
	clientSecret := os.Getenv("CHOREO_MYSERVICE_CONSUMER_SECRET")
	tokenURL := os.Getenv("CHOREO_MYSERVICE_TOKEN_URL")

	if serviceURL == "" || clientID == "" || clientSecret == "" || tokenURL == "" {
		http.Error(w, "Missing required environment variables", http.StatusInternalServerError)
		return
	}

	// Set up OAuth2 configuration
	oauth2Config := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
	}

	// Create an HTTP client with OAuth2 token
	client := oauth2Config.Client(context.Background())

	// Make a request to the actual greeter service
	greeterServiceURL := fmt.Sprintf("%s/greeter/greet?name=%s", serviceURL, "name")
	resp, err := client.Get(greeterServiceURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to make a request to Greeter service: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response from the Greeter service
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read the response body", http.StatusInternalServerError)
		return
	}

	// Write the response from the Greeter service to the HTTP response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
