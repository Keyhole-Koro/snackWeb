package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"snackWeb/backend/internal/adapter"
	"snackWeb/backend/internal/db"
	"snackWeb/backend/internal/repository"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/rs/cors"
)

func main() {
	// Initialize Database (DynamoDB)
	if err := db.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	// No defer db.CloseDB() needed for DynamoDB client v2

	mux := http.NewServeMux()

	// ── Endpoints ────────────────────────────────────────────────

	// GET /api/feed
	mux.HandleFunc("/api/feed", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		limitStr := r.URL.Query().Get("limit")
		limit := 50
		if limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
				limit = l
			}
		}

		feed, err := repository.GetPosts(limit)
		if err != nil {
			log.Println("Error getting feed:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(feed)
	})

	// GET /api/personas
	mux.HandleFunc("/api/personas", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		personas, err := repository.GetPersonas()
		if err != nil {
			log.Println("Error getting personas:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(personas)
	})

	// GET /api/stats
	mux.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		stats, err := repository.GetGenerationStats()
		if err != nil {
			log.Println("Error getting stats:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(stats)
	})

	// Mount plots directory (Only relevant for local/EC2, not Lambda usually, unless EFS)
	// In Lambda, we might serve via S3 presigned URLs or CloudFront, but if mounted EFS/S3-Local...
	// Just keep logic for local dev/legacy.
	plotsDir := os.Getenv("PLOTS_DIR")
	if plotsDir == "" {
		cwd, _ := os.Getwd()
		// Assuming running from root or cmd/server
		plotsDir = filepath.Join(cwd, "../../../persona_data/plots")
	}
	// Check existence before serving to avoid panic in FileServer? No, FileServer is safe.
	fs := http.FileServer(http.Dir(plotsDir))
	mux.Handle("/plots/", http.StripPrefix("/plots/", fs))

	// Middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	handler := c.Handler(mux)

	// Check execution mode
	mode := os.Getenv("EXECUTION_MODE")
	if mode == "lambda" {
		log.Println("Starting Lambda Handler")
		lambdaAdapter := adapter.New(handler)
		lambda.Start(lambdaAdapter.Proxy)
	} else {
		port := os.Getenv("PORT")
		if port == "" {
			port = "13579"
		}
		log.Printf("🍿 SnackWeb Go Backend running at http://localhost:%s", port)
		log.Fatal(http.ListenAndServe(":"+port, handler))
	}
}
