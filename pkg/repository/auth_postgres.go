package repository

import (
	"Diploma/entity"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateTeam(team *entity.Teams) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (team_name,player1,player2,player3,player4,player5,password) values ($1,$2,$3,$4,$5,$6,$7) Returning team_id", teamtable)
	row := r.db.QueryRow(query, team.Team_name, team.Player1, team.Player2, team.Player3, team.Player4, team.Player5, team.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// func (r *AuthPostgres) CreateTeam111(team Diploma.Teams) (int64, error) {

// 	result, err := r.db.Exec("insert into teams (team_name, player1, player2, player3, player4, player5, password) values($1,$2,$3,$4,$5,$6,$7)",
// 		team.Team_name, team.Player1, team.Player2, team.Player3, team.Player4, team.Player5, team.Password)

// 	if err != nil {
// 		return 0, err
// 	}

// 	lastInsertId, _ := result.LastInsertId()

//		return lastInsertId, nil
//	}
func (r *AuthPostgres) GetUser(teamname, password string) (entity.Teams, error) {
	var Team entity.Teams
	query := fmt.Sprintf("SELECT team_id from %s where team_name=$1 and password=$2", teamtable)
	err := r.db.Get(&Team, query, teamname, password)
	return Team, err
}
