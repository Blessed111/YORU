package entity

import "time"

type Tournament struct {
	TournamentID       int       `json:"tournament_id"`
	TournamentName     string    `json:"tournament_name"`
	Description        string    `json:"description"`
	StartDate          time.Time `json:"start_date"`
	EndDate            time.Time `json:"end_date"`
	TeamsCount         int       `json:"teams_count"`
	TotalRoundNumber   int       `json:"total_round_number"`
	StatusActive       bool      `json:"status_active"`
	WinnerTeamID       int       `json:"winner_team_id"`
	CurrentTeamsCount  int       `json:"current_teams_count"`
	CurrentRoundNumber int       `json:"current_round_number"`
	Registered         bool      `json:"registered"`
	WinnerName         string
}
