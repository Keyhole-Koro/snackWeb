package models

import (
	"database/sql"
	"time"
)

// Persona corresponds to the persona table
type Persona struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Bio       string         `json:"bio"`
	Metadata  sql.NullString `json:"-"` // Internal use, handled as raw JSON string if needed
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	IsActive  bool           `json:"is_active"`
	Stats     *PersonaStats  `json:"stats,omitempty"` // Enriched field
}

// PersonaStats corresponds to a subset of fields from fitness_snapshot
type PersonaStats struct {
	PostQuality   float64 `json:"post_quality"`
	Incisiveness  float64 `json:"incisiveness"`
	Judiciousness float64 `json:"judiciousness"`
	RawFitness    float64 `json:"raw_fitness"`
}

// Post corresponds to the post table
type Post struct {
	ID              string         `json:"id"`
	PersonaID       string         `json:"persona_id"`
	GenerationID    sql.NullInt64  `json:"generation_id"`
	Content         string         `json:"content"`
	Topic           string         `json:"topic"`
	EventType       string         `json:"event_type"`
	ResearchContext sql.NullString `json:"-"`
	CreatedAt       time.Time      `json:"created_at"`
	AgentName       string         `json:"agent_name"` // Joined field
}

// Reply corresponds to the reply table
type Reply struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	PersonaID string    `json:"persona_id"`
	Content   string    `json:"content"`
	EventType string    `json:"event_type"`
	CreatedAt time.Time `json:"created_at"`
	AgentName string    `json:"agent_name"` // Joined field
}

// FeedItem is the response structure for the feed API
type FeedItem struct {
	ID        string    `json:"id"`
	Timestamp string    `json:"timestamp"`
	EventType string    `json:"event_type"`
	AgentName string    `json:"agent_name"`
	Content   string    `json:"content"`
	Topic     string    `json:"topic"`
	Replies   []APIReply `json:"replies"`
}

type APIReply struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	EventType string `json:"event_type"`
	AgentName string `json:"agent_name"`
	Content   string `json:"content"`
}

// GenerationStats corresponds to the generation table
type GenerationStats struct {
	Generation          int     `json:"generation"`
	PopulationDiversity float64 `json:"population_diversity"`
	FitnessMean         float64 `json:"fitness_mean"`
}
