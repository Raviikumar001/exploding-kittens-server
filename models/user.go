package models

import (
	"github.com/google/uuid"
)


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


type UpdateData struct {
    ID         string `json:"ID"`
    GameResult bool   `json:"gameResult"`
}


type IDData struct { 
	ID string `json:"ID"`
} 