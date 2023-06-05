package repository

import (
	"Diploma/entity"

	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateTeam(team *entity.Teams) (int, error)
	GetUser(teamname, password string) (entity.Teams, error)
}
type Repository struct {
	Authorization
	TournamentRepositoryI
}

func NewConnection(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization:         NewAuthPostgres(db),
		TournamentRepositoryI: newTournamentRepository(db),
	}
}
