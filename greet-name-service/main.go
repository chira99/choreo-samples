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
	// HTTP server to expose multiple endpoints
	serverMux := http.NewServeMux()

	// Register the /greeter/greet handler
	serverMux.HandleFunc("/greeter/greet", greetHandler)

	// Register a new endpoint for the second service
	serverMux.HandleFunc("/greeter/world", worldHandler)

	serverPort := 8080
	server := http.Server{
		Addr:    fmt.Sprintf(":%d", serverPort),
		Handler: serverMux,
	}

	go func() {
		log.Printf("Starting HTTP Server on port %d\n", serverPort)
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
	makeOAuth2Request(w, r, "SERVICE1")
}

// worldHandler handles the /another-service/action API call
func worldHandler(w http.ResponseWriter, r *http.Request) {
	makeOAuth2Request(w, r, "SERVICE2")
}

// makeOAuth2Request makes an OAuth2 authenticated request to a service
// Takes in a `serviceType` parameter to determine which environment variables to use
func makeOAuth2Request(w http.ResponseWriter, r *http.Request, serviceType string) {
	var serviceURL, clientID, clientSecret, tokenURL string

	// Choose environment variables based on the serviceType
	switch serviceType {
	case "SERVICE1":
		serviceURL = os.Getenv("CHOREO_CONNECT_MYSERVICE_TO_TESTSERVICE_SERVICEURL")
		clientID = os.Getenv("CHOREO_CONNECT_MYSERVICE_TO_TESTSERVICE_CONSUMERKEY")
		clientSecret = os.Getenv("CHOREO_CONNECT_MYSERVICE_TO_TESTSERVICE_CONSUMERSECRET")
		tokenURL = os.Getenv("CHOREO_CONNECT_MYSERVICE_TO_TESTSERVICE_TOKENURL")

	case "SERVICE2":
		serviceURL = os.Getenv("CHOREO_CONNECT_MYSERVICE_TO_TESTSERVICE2_SERVICEURL")
		clientID = os.Getenv("CHOREO_CONNECT_MYSERVICE_TO_TESTSERVICE2_CONSUMERKEY")
		clientSecret = os.Getenv("CHOREO_CONNECT_MYSERVICE_TO_TESTSERVICE2_CONSUMERSECRET")
		tokenURL = os.Getenv("CHOREO_CONNECT_MYSERVICE_TO_TESTSERVICE2_TOKENURL")

	default:
		http.Error(w, "Invalid service type", http.StatusInternalServerError)
		return
	}

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
	// Construct the service URL based on service type

	var serviceRequestURL string
	if serviceType == "SERVICE1" {
		serviceRequestURL = fmt.Sprintf("%s/greeter/greet?name=%s", serviceURL, "name")
	} else {
		// For SERVICE2, just use the base URL and path
		serviceRequestURL = fmt.Sprintf("%s/greeter/world", serviceURL)
	}

	// Make a request to the specified service API path
	resp, err := client.Get(serviceRequestURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to make a request to service: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read the response from the service
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read the response body", http.StatusInternalServerError)
		return
	}

	// Write the response from the service to the HTTP response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
