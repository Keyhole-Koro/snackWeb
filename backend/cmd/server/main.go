package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"snackWeb/backend/internal/db"
	"snackWeb/backend/internal/repository"

	"github.com/rs/cors"
)

func main() {
	// Initialize Database
	if err := db.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.CloseDB()

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

	// Mount plots directory
	plotsDir := os.Getenv("PLOTS_DIR")
	if plotsDir == "" {
		cwd, _ := os.Getwd()
		// Assuming running from root or cmd/server, aim for persona_data/plots
		// Let's rely on standard structure: /home/unix/snack/persona_data/plots
		plotsDir = filepath.Join(cwd, "../../../persona_data/plots")
	}
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "13579"
	}

	log.Printf("🍿 SnackWeb Go Backend running at http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
