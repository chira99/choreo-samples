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
	serverMux.HandleFunc("/greeter/greetGo", greetGo)
	serverMux.HandleFunc("/greeter/greetGo2", greetGo2)
	serverMux.HandleFunc("/greeter/greetGoProj", greetGoProj)
	serverMux.HandleFunc("/greeter/greetGoProj2", greetGoProj2)
	serverMux.HandleFunc("/greeter/greetGoOrg", greetGoOrg)
	serverMux.HandleFunc("/greeter/greetGoOrg2", greetGoOrg2)
	serverMux.HandleFunc("/greeter/greetGoPublic", greetGoPublic)
	serverMux.HandleFunc("/greeter/greetGoPublic2", greetGoPublic2)

	// serverMux.HandleFunc("/greeter/greetOrg", greetHandlerOrg)
	// serverMux.HandleFunc("/greeter/greetDb", greetDb)
	// serverMux.HandleFunc("/greeter/greetDbProd", greetDbProd)

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

func greetGo(w http.ResponseWriter, r *http.Request) {
	serviceURL := os.Getenv("CHOREO_CONNECT_MULTI_1_SERVICEURL")
	clientID := os.Getenv("CHOREO_CONNECT_MULTI_1_CONSUMERKEY")
	clientSecret := os.Getenv("CHOREO_CONNECT_MULTI_1_CONSUMERSECRET")
	tokenURL := os.Getenv("CHOREO_CONNECT_MULTI_1_TOKENURL")

	if serviceURL == "" || clientID == "" || clientSecret == "" || tokenURL == "" {
		missingVars := []string{}
		if serviceURL == "" {
			missingVars = append(missingVars, "CHOREO_CONNECT_MULTI_1_SERVICEURL")
		}
		if clientID == "" {
			missingVars = append(missingVars, "CHOREO_CONNECT_MULTI_1_CONSUMERKEY")
		}
		if clientSecret == "" {
			missingVars = append(missingVars, "CHOREO_CONNECT_MULTI_1_CONSUMERSECRET")
		}
		if tokenURL == "" {
			missingVars = append(missingVars, "CHOREO_CONNECT_MULTI_1_TOKENURL")
		}
		http.Error(w, fmt.Sprintf("Missing required environment variables: %v", missingVars), http.StatusInternalServerError)
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
	serviceRequestURL := fmt.Sprintf("%s/greeter/greet?name=%s", serviceURL, "chira")

	// Make a request to the specified service API path
	resp, err := client.Get(serviceRequestURL)
	if err != nil {
		log.Printf("Failed to make a request to GO-1 service: %v", err)
		http.Error(w, fmt.Sprintf("Failed to make a request to Go-1 service: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read the response body for Go-1 service")
		http.Error(w, "Failed to read the response body of Go-1", http.StatusInternalServerError)
		return
	}

	// Write the response from the service to the HTTP response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

// greetHandler handles the /greeter/greet API call for PROJECT
func greetGo2(w http.ResponseWriter, r *http.Request) {
	serviceURL := os.Getenv("CHOREO_CONNECT_MULTI_2_SERVICEURL")

	if serviceURL == "" {
		http.Error(w, "Missing required environment variable: CHOREO_CONNECT_MULTI_2_SERVICEURL", http.StatusInternalServerError)
		return
	}

	serviceRequestURL := fmt.Sprintf("%s/greeter/greet?name=%s", serviceURL, "chira")

	// Make a request to the specified service API path
	resp, err := http.Get(serviceRequestURL)
	if err != nil {
		log.Printf("Failed to make a request to GO-2 service: %v", err)
		http.Error(w, fmt.Sprintf("Failed to make a request to Go-2 service: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read the response body for Go-2 service")
		http.Error(w, "Failed to read the response body of Go-2", http.StatusInternalServerError)
		return
	}

	// Write the response from the service to the HTTP response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
func greetGoProj(w http.ResponseWriter, r *http.Request) {
	// sample go code snippet
	serviceURL := os.Getenv("CHOREO_CONNECT_PROJ_1_SERVICEURL")

	if serviceURL == "" {
		http.Error(w, "Missing required environment variable: CHOREO_CONNECT_PROJ_1_SERVICEURL", http.StatusInternalServerError)
		return
	}

	serviceRequestURL := fmt.Sprintf("%s/greeter/greet?name=%s", serviceURL, "chira")

	// Make a request to the specified service API path
	resp, err := http.Get(serviceRequestURL)
	if err != nil {
		log.Printf("Failed to make a request to GO-Proj1 service: %v", err)
		http.Error(w, fmt.Sprintf("Failed to make a request to Go-Proj1 service: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read the response body for Go-Proj1 service")
		http.Error(w, "Failed to read the response body of Go-Proj1", http.StatusInternalServerError)
		return
	}

	// Write the response from the service to the HTTP response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func greetGoProj2(w http.ResponseWriter, r *http.Request) {
	serviceURL := os.Getenv("CHOREO_CONNECT_PROJ_2_SERVICEURL")

	if serviceURL == "" {
		http.Error(w, "Missing required environment variable: CHOREO_CONNECT_PROJ_2_SERVICEURL", http.StatusInternalServerError)
		return
	}

	serviceRequestURL := fmt.Sprintf("%s/greeter/greet?name=%s", serviceURL, "chira")

	// Make a request to the specified service API path
	resp, err := http.Get(serviceRequestURL)
	if err != nil {
		log.Printf("Failed to make a request to GO-proj2 service: %v", err)
		http.Error(w, fmt.Sprintf("Failed to make a request to Go-Proj2 service: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read the response body for Go-Proj2 service")
		http.Error(w, "Failed to read the response body of Go-Proj2", http.StatusInternalServerError)
		return
	}

	// Write the response from the service to the HTTP response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func greetGoOrg(w http.ResponseWriter, r *http.Request) {
	serviceURL := os.Getenv("CHOREO_CONNECT_ORG_1_SERVICEURL")
	clientID := os.Getenv("CHOREO_CONNECT_ORG_1_CONSUMERKEY")
	clientSecret := os.Getenv("CHOREO_CONNECT_ORG_1_CONSUMERSECRET")
	tokenURL := os.Getenv("CHOREO_CONNECT_ORG_1_TOKENURL")

	if serviceURL == "" || clientID == "" || clientSecret == "" || tokenURL == "" {
		missingVars := []string{}
		if serviceURL == "" {
			missingVars = append(missingVars, "CHOREO_CONNECT_ORG_1_SERVICEURL")
		}
		if clientID == "" {
			missingVars = append(missingVars, "CHOREO_CONNECT_ORG_1_CONSUMERKEY")
		}
		if clientSecret == "" {
			missingVars = append(missingVars, "CHOREO_CONNECT_ORG_1_CONSUMERSECRET")
		}
		if tokenURL == "" {
			missingVars = append(missingVars, "CHOREO_CONNECT_ORG_1_TOKENURL")
		}
		http.Error(w, fmt.Sprintf("Missing required environment variables: %v", missingVars), http.StatusInternalServerError)
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
	serviceRequestURL := fmt.Sprintf("%s/greeter/greet?name=%s", serviceURL, "chira")

	// Make a request to the specified service API path
	resp, err := client.Get(serviceRequestURL)
	if err != nil {
		log.Printf("Failed to make a request to GO-org-1 service: %v", err)
		http.Error(w, fmt.Sprintf("Failed to make a request to GO-org-1 service: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read the response body for GO-org-1 service")
		http.Error(w, "Failed to read the response body of GO-org-1", http.StatusInternalServerError)
		return
	}

	// Write the response from the service to the HTTP response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
func greetGoOrg2(w http.ResponseWriter, r *http.Request) {
	serviceURL := os.Getenv("CHOREO_CONNECT_ORG_2_SERVICEURL")
	clientID := os.Getenv("CHOREO_CONNECT_ORG_2_CONSUMERKEY")
	clientSecret := os.Getenv("CHOREO_CONNECT_ORG_2_CONSUMERSECRET")
	tokenURL := os.Getenv("CHOREO_CONNECT_ORG_2_TOKENURL")

	if serviceURL == "" || clientID == "" || clientSecret == "" || tokenURL == "" {
		missingVars := []string{}
		if serviceURL == "" {
			missingVars = append(missingVars, "CHOREO_CONNECT_ORG_2_SERVICEURL")
		}
		if clientID == "" {
			missingVars = append(missingVars, "CHOREO_CONNECT_ORG_2_CONSUMERKEY")
		}
		if clientSecret == "" {
			missingVars = append(missingVars, "CHOREO_CONNECT_ORG_2_CONSUMERSECRET")
		}
		if tokenURL == "" {
			missingVars = append(missingVars, "CHOREO_CONNECT_ORG_2_TOKENURL")
		}
		http.Error(w, fmt.Sprintf("Missing required environment variables: %v", missingVars), http.StatusInternalServerError)
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
	serviceRequestURL := fmt.Sprintf("%s/greeter/greet?name=%s", serviceURL, "chira")

	// Make a request to the specified service API path
	resp, err := client.Get(serviceRequestURL)
	if err != nil {
		log.Printf("Failed to make a request to GO-org-2 service: %v", err)
		http.Error(w, fmt.Sprintf("Failed to make a request to GO-org-2 service: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read the response body for GO-org-2 service")
		http.Error(w, "Failed to read the response body of GO-org-2", http.StatusInternalServerError)
		return
	}

	// Write the response from the service to the HTTP response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func greetGoPublic(w http.ResponseWriter, r *http.Request) {
	serviceURL := os.Getenv("CHOREO_CONNECT_PUBLIC_1_PROJ_SERVICEURL")
	clientID := os.Getenv("CHOREO_CONNECT_PUBLIC_1_PROJ_CONSUMERKEY")
	clientSecret := os.Getenv("CHOREO_CONNECT_PUBLIC_1_PROJ_CONSUMERSECRET")
	tokenURL := os.Getenv("CHOREO_CONNECT_PUBLIC_1_PROJ_TOKENURL")

	if serviceURL == "" || clientID == "" || clientSecret == "" || tokenURL == "" {
		missingVars := []string{}
		if serviceURL == "" {
			missingVars = append(missingVars, "CHOREO_CONNECT_PUBLIC_1_PROJ_SERVICEURL")
		}
		if clientID == "" {
			missingVars = append(missingVars, "CHOREO_CONNECT_PUBLIC_1_PROJ_CONSUMERKEY")
		}
		if clientSecret == "" {
			missingVars = append(missingVars, "CHOREO_CONNECT_PUBLIC_1_PROJ_CONSUMERSECRET")
		}
		if tokenURL == "" {
			missingVars = append(missingVars, "CHOREO_CONNECT_PUBLIC_1_PROJ_TOKENURL")
		}
		http.Error(w, fmt.Sprintf("Missing required environment variables: %v", missingVars), http.StatusInternalServerError)
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
	serviceRequestURL := fmt.Sprintf("%s/greeter/greet?name=%s", serviceURL, "chira")

	// Make a request to the specified service API path
	resp, err := client.Get(serviceRequestURL)
	if err != nil {
		log.Printf("Failed to make a request to GO-public-1 service: %v", err)
		http.Error(w, fmt.Sprintf("Failed to make a request to GO-public-1 service: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read the response body for GO-public-1 service")
		http.Error(w, "Failed to read the response body of GO-public-1", http.StatusInternalServerError)
		return
	}

	// Write the response from the service to the HTTP response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
func greetGoPublic2(w http.ResponseWriter, r *http.Request) {
	serviceURL := os.Getenv("CHOREO_CONNECT_PUBLIC_2_PROJ_SERVICEURL")
	clientID := os.Getenv("CHOREO_CONNECT_PUBLIC_2_PROJ_CONSUMERKEY")
	clientSecret := os.Getenv("CHOREO_CONNECT_PUBLIC_2_PROJ_CONSUMERSECRET")
	tokenURL := os.Getenv("CHOREO_CONNECT_PUBLIC_2_PROJ_TOKENURL")

	if serviceURL == "" || clientID == "" || clientSecret == "" || tokenURL == "" {
		missingVars := []string{}
		if serviceURL == "" {
			missingVars = append(missingVars, "CHOREO_CONNECT_PUBLIC_2_PROJ_SERVICEURL")
		}
		if clientID == "" {
			missingVars = append(missingVars, "CHOREO_CONNECT_PUBLIC_2_PROJ_CONSUMERKEY");
		}
		if clientSecret == "" {
			missingVars = append(missingVars, "CHOREO_CONNECT_PUBLIC_2_PROJ_CONSUMERSECRET");
		}
		if tokenURL == "" {
			missingVars = append(missingVars, "CHOREO_CONNECT_PUBLIC_2_PROJ_TOKENURL");
		}
		http.Error(w, fmt.Sprintf("Missing required environment variables: %v", missingVars), http.StatusInternalServerError)
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
	serviceRequestURL := fmt.Sprintf("%s/greeter/greet?name=%s", serviceURL, "chira")

	// Make a request to the specified service API path
	resp, err := client.Get(serviceRequestURL)
	if err != nil {
		log.Printf("Failed to make a request to GO-public-2 service: %v", err)
		http.Error(w, fmt.Sprintf("Failed to make a request to GO-public-2 service: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read the response body for GO-public-2 service")
		http.Error(w, "Failed to read the response body of GO-public-2", http.StatusInternalServerError)
		return
	}

	// Write the response from the service to the HTTP response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

// // greetHandlerOrg handles the /greeter/greetOrg API call for ORG
// func greetHandlerOrg(w http.ResponseWriter, r *http.Request) {
// 	makeOrgRequest(w, r)
// }

// func greetDbProd(w http.ResponseWriter, r *http.Request) {
// 	hostName := os.Getenv("CHOREO_CONNECTDBPROD_HOSTNAME")
// 	port := os.Getenv("CHOREO_CONNECTDBPROD_PORT")
// 	username := os.Getenv("CHOREO_CONNECTDBPROD_USERNAME")
// 	password := os.Getenv("CHOREO_CONNECTDBPROD_PASSWORD")
// 	dbName := os.Getenv("CHOREO_CONNECTDBPROD_DATABASENAME")

// 	if hostName == "" || port == "" || username == "" || password == "" || dbName == "" {
// 		missingVars := []string{}
// 		if hostName == "" {
// 			missingVars = append(missingVars, "HNAME")
// 		}
// 		if port == "" {
// 			missingVars = append(missingVars, "PORT");
// 		}
// 		if username == "" {
// 			missingVars = append(missingVars, "UNAME");
// 		}
// 		if password == "" {
// 			missingVars = append(missingVars, "PWD");
// 		}
// 		if dbName == "" {
// 			missingVars = append(missingVars, "DBNAME");
// 		}
// 		http.Error(w, fmt.Sprintf("Missing required environment variables: %v", missingVars), http.StatusInternalServerError)
// 		return
// 	}

// 	// Display the environment variables as a simple response
// 	w.Header().Set("Content-Type", "text/plain")
// 	fmt.Fprintf(w, "DB Environment Variables:\n1. Hostname: %s\n2. Port: %s\n3. Username: %s\n4. Password: %s\n5. Database Name: %s\n", hostName, port, username, password, dbName)
// }

// func greetDb(w http.ResponseWriter, r *http.Request) {
// 	makeDbReq(w, r)
// }

// makeProjectRequest makes a request to the PROJECT service without OAuth2 authentication
// func makeProjectRequest(w http.ResponseWriter, r *http.Request) {
// 	serviceURL := os.Getenv("CHOREO_CONNECT_TO_GREETER_PROJECTACCESS_SERVICEURL")
// 	if serviceURL == "" {
// 		http.Error(w, "Missing required environment variable: CHOREO_CONNECT_TO_GREETER_PROJECTACCESS_SERVICEURL", http.StatusInternalServerError)
// 		return
// 	}

// 	serviceRequestURL := fmt.Sprintf("%s/greeter/greet?name=%s", serviceURL, "person-project")
// 	resp, err := http.Get(serviceRequestURL)
// 	if err != nil {
// 		log.Printf("Failed to make a request to PROJECT service: %v", err)
// 		http.Error(w, fmt.Sprintf("Failed to make a request to service: %v", err), http.StatusInternalServerError)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Println("Failed to read the response body for PROJECT service")
// 		http.Error(w, "Failed to read the response body", http.StatusInternalServerError)
// 		return
// 	}

// 	// Write the response from the service to the HTTP response
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(resp.StatusCode)
// 	w.Write(body)
// }

// func makeDbReq(w http.ResponseWriter, r *http.Request) {
// 	hostName := os.Getenv("CHOREO_CONNECTDBDEV_HOSTNAME")
// 	port := os.Getenv("CHOREO_CONNECTDBDEV_PORT")
// 	username := os.Getenv("CHOREO_CONNECTDBDEV_USERNAME")
// 	password := os.Getenv("CHOREO_CONNECTDBDEV_PASSWORD")
// 	dbName := os.Getenv("CHOREO_CONNECTDBDEV_DATABASENAME")

// 	if hostName == "" || port == "" || username == "" || password == "" || dbName == "" {
// 		missingVars := []string{}
// 		if hostName == "" {
// 			missingVars = append(missingVars, "HNAME")
// 		}
// 		if port == "" {
// 			missingVars = append(missingVars, "PORT");
// 		}
// 		if username == "" {
// 			missingVars = append(missingVars, "UNAME");
// 		}
// 		if password == "" {
// 			missingVars = append(missingVars, "PWD");
// 		}
// 		if dbName == "" {
// 			missingVars = append(missingVars, "DBNAME");
// 		}
// 		http.Error(w, fmt.Sprintf("Missing required environment variables: %v", missingVars), http.StatusInternalServerError)
// 		return
// 	}

// 	// Display the environment variables as a simple response
// 	w.Header().Set("Content-Type", "text/plain")
// 	fmt.Fprintf(w, "DB dev Environment Variables:\n1. Hostname: %s\n2. Port: %s\n3. Username: %s\n4. Password: %s\n5. Database Name: %s\n", hostName, port, username, password, dbName)
// }
