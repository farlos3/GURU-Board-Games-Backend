package search

import (
	"log"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/joho/godotenv"
	"os"
)

// Global Elasticsearch client
var es *elasticsearch.Client

// Initialize Elasticsearch connection
func InitElasticsearch() {
	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	cloudID := os.Getenv("CLOUD_ID") 
	apiKey := os.Getenv("API_KEY")

	// Ensure that both variables are set
	if cloudID == "" || apiKey == "" {
		log.Fatal("CLOUD_ID and API_KEY must be set in the .env file")
	}

	// Set up the Elasticsearch client configuration
	cfg := elasticsearch.Config{
		CloudID: cloudID,
		APIKey:  apiKey,
	}

	// Create a new Elasticsearch client
	es, err = elasticsearch.NewClient(cfg)

	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// Optional: Test the connection (useful for debugging)
	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting Elasticsearch info: %s", err)
	}
	defer res.Body.Close()

	log.Printf("✅ Connected to Elasticsearch: %s", res)
}