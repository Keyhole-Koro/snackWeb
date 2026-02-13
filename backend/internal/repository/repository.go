package repository

import (
	"database/sql"
	"log"
	"snackWeb/backend/internal/db"
	"snackWeb/backend/internal/models"
)

// GetPosts returns a list of feed items with nested replies
func GetPosts(limit int) ([]models.FeedItem, error) {
	query := `
		SELECT 
			p.id, 
			p.created_at, 
			p.event_type, 
			COALESCE(per.name, '') as agent_name, 
			p.content, 
			p.topic
		FROM post p
		LEFT JOIN persona per ON p.persona_id = per.id
		ORDER BY p.created_at DESC
		LIMIT ?
	`
	rows, err := db.DB.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feed []models.FeedItem
	for rows.Next() {
		var item models.FeedItem
		// Actually, let's scan into time.Time and format later
		var createdAt sql.NullTime

		if err := rows.Scan(&item.ID, &createdAt, &item.EventType, &item.AgentName, &item.Content, &item.Topic); err != nil {
			log.Println("Error scanning post:", err)
			continue
		}
		if createdAt.Valid {
			item.Timestamp = createdAt.Time.Format("2006-01-02T15:04:05.000000")
		}

		replies, err := GetRepliesForPost(item.ID)
		if err != nil {
			log.Println("Error fetching replies:", err)
			item.Replies = []models.APIReply{}
		} else {
			item.Replies = replies
		}

		feed = append(feed, item)
	}

	return feed, nil
}

func GetRepliesForPost(postID string) ([]models.APIReply, error) {
	query := `
		SELECT 
			r.id, 
			r.created_at, 
			r.event_type, 
			COALESCE(per.name, '') as agent_name, 
			r.content
		FROM reply r
		LEFT JOIN persona per ON r.persona_id = per.id
		WHERE r.post_id = ?
		ORDER BY r.created_at ASC
	`
	rows, err := db.DB.Query(query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var replies []models.APIReply
	for rows.Next() {
		var r models.APIReply
		var createdAt sql.NullTime
		if err := rows.Scan(&r.ID, &createdAt, &r.EventType, &r.AgentName, &r.Content); err != nil {
			continue
		}
		if createdAt.Valid {
			r.Timestamp = createdAt.Time.Format("2006-01-02T15:04:05.000000")
		}
		replies = append(replies, r)
	}
	// Return empty slice instead of nil for JSON consistency
	if replies == nil {
		replies = []models.APIReply{}
	}
	return replies, nil
}

func GetPersonas() ([]models.Persona, error) {
	query := `SELECT id, name, bio, is_active FROM persona WHERE is_active = 1`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var personas []models.Persona
	for rows.Next() {
		var p models.Persona
		// is_active is boolean in Go, but Integer 0/1 in SQLite. driver handles it? 
		// mattn/go-sqlite3 handles scans to types well.
		if err := rows.Scan(&p.ID, &p.Name, &p.Bio, &p.IsActive); err != nil {
			log.Println("Scan persona error:", err)
			continue
		}
		
		stats, err := GetLatestFitness(p.ID)
		if err == nil {
			p.Stats = stats
		}
		personas = append(personas, p)
	}
	return personas, nil
}

func GetLatestFitness(personaID string) (*models.PersonaStats, error) {
	query := `
		SELECT 
			post_quality, 
			incisiveness, 
			judiciousness, 
			raw_fitness 
		FROM fitness_snapshot 
		WHERE persona_id = ? 
		ORDER BY created_at DESC 
		LIMIT 1
	`
	var s models.PersonaStats
	err := db.DB.QueryRow(query, personaID).Scan(&s.PostQuality, &s.Incisiveness, &s.Judiciousness, &s.RawFitness)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func GetGenerationStats() ([]models.GenerationStats, error) {
	query := `
		SELECT 
			id, 
			population_diversity, 
			fitness_mean 
		FROM generation 
		ORDER BY id ASC
	`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []models.GenerationStats
	for rows.Next() {
		var gs models.GenerationStats
		// Scan NULLs? Schema says default 0.0 but let's be safe if needed. 
		// Assuming non-null for simplified logic
		if err := rows.Scan(&gs.Generation, &gs.PopulationDiversity, &gs.FitnessMean); err != nil {
			continue
		}
		stats = append(stats, gs)
	}
	return stats, nil
}
