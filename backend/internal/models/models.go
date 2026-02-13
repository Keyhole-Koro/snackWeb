package models

import (
	"time"
)

// Persona corresponds to the persona table
type Persona struct {
	ID        string    `json:"id" dynamodbav:"id"`
	Name      string    `json:"name" dynamodbav:"name"`
	Bio       string    `json:"bio" dynamodbav:"bio"`
	Metadata  *string   `json:"-" dynamodbav:"metadata,omitempty"` 
	CreatedAt string    `json:"created_at" dynamodbav:"created_at"`
	UpdatedAt time.Time `json:"updated_at" dynamodbav:"updated_at"`
	IsActive  bool      `json:"is_active" dynamodbav:"is_active"`
	Stats     *PersonaStats `json:"stats,omitempty" dynamodbav:"stats,omitempty"`
}

// PersonaStats corresponds to a subset of fields from fitness_snapshot
type PersonaStats struct {
	PostQuality   float64 `json:"post_quality" dynamodbav:"post_quality"`
	Incisiveness  float64 `json:"incisiveness" dynamodbav:"incisiveness"`
	Judiciousness float64 `json:"judiciousness" dynamodbav:"judiciousness"`
	RawFitness    float64 `json:"raw_fitness" dynamodbav:"raw_fitness"`
}

// Post corresponds to the post table
type Post struct {
	ID              string    `json:"id" dynamodbav:"id"`
	PersonaID       string    `json:"persona_id" dynamodbav:"persona_id"`
	GenerationID    *int64    `json:"generation_id" dynamodbav:"generation_id,omitempty"`
	Content         string    `json:"content" dynamodbav:"content"`
	Topic           string    `json:"topic" dynamodbav:"topic"`
	EventType       string    `json:"event_type" dynamodbav:"event_type"`
	ResearchContext *string   `json:"-" dynamodbav:"research_context,omitempty"`
	CreatedAt       string    `json:"created_at" dynamodbav:"created_at"`
	AgentName       string    `json:"agent_name" dynamodbav:"agent_name"` 
}

// Reply corresponds to the reply table
type Reply struct {
	ID        string    `json:"id" dynamodbav:"id"`
	PostID    string    `json:"post_id" dynamodbav:"post_id"`
	PersonaID string    `json:"persona_id" dynamodbav:"persona_id"`
	Content   string    `json:"content" dynamodbav:"content"`
	EventType string    `json:"event_type" dynamodbav:"event_type"`
	CreatedAt string    `json:"created_at" dynamodbav:"created_at"`
	AgentName string    `json:"agent_name" dynamodbav:"agent_name"`
}

// FeedItem is the response structure for the feed API
type FeedItem struct {
	ID        string     `json:"id"`
	Timestamp string     `json:"timestamp"`
	EventType string     `json:"event_type"`
	AgentName string     `json:"agent_name"`
	Content   string     `json:"content"`
	Topic     string     `json:"topic"`
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
	Generation          int     `json:"generation" dynamodbav:"generation"`
	PopulationDiversity float64 `json:"population_diversity" dynamodbav:"population_diversity"`
	FitnessMean         float64 `json:"fitness_mean" dynamodbav:"fitness_mean"`
}
