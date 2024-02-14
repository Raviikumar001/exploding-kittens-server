package models

import (
	"github.com/google/uuid"
)

// User represents a player in your game system
type SignUp struct {
	Name     string `json:"name"`
	Username string `json:"username"`
}

type SignIn struct {
	Username string `json:"username"`
}

type User struct {
	ID             uuid.UUID `json:"ID"`
	Name           string    `json:"name"`
	Username       string    `json:"username"`
	TotalPoints    int       `json:"total_points"`
	TotalGamesLost int       `json:"total_games_lost"`
}
