package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"Diploma/entity"

	"github.com/jmoiron/sqlx"
)

type TournamentRepositoryI interface {
	CreateTournament(tour entity.Tournament) error
	GetAllTournament(active *bool) ([]entity.Tournament, error)
	RegisterTeam(tournamentID, teamID int) error
	GetAllParticipiant(teamID *int) ([]entity.Participant, error)
	GetAllParticipiantByTourID(tournamentID int) ([]int, error)
	UnregisterTeam(tournamentID, teamID int) error
	ActiveTournament(touramentID int) error
	InactiveTournament(touramentID int) error
	GetTournamentByID(tournamentID int) (entity.Tournament, error)
	CreateMatch(match entity.Match) error
	GetAllMatchesByTournamentID(tournamentID int) ([]entity.Match, error)

	UpdateWinnerMatch(mathcID, winnerID, loserID int) error
	GetMatchesByTournamentID(tournamentID int) ([]entity.Match, error)
	GetAllMatchesByTournamentIDWithName(tournamentID int) ([]entity.Match, error)
	GetWinnersByRoundTournament(roundNumber, tournamentID int) ([]int, error)

	CheckWinnerLoserIDsNotEmpty(roundNumber, tournamentID int) (bool, error)
	UpdateCurrentRoundNumber(tournamentID int) error
	TournamentWinner(teamID, tournamentID int) error
	GetMatchByID(matchID int) (entity.Match, error)
	GetTeamByID(teamID int) (entity.Teams, error)

	GetMatchForNextLevel(tournamentID, roundNumber, participant1, participant2 int) (entity.Match, error)
}

type tournamentRepository struct {
	db *sqlx.DB
}

func newTournamentRepository(db *sqlx.DB) TournamentRepositoryI {
	return &tournamentRepository{
		db: db,
	}
}

func (t *tournamentRepository) GetTeamByID(teamID int) (entity.Teams, error) {
	query := `SELECT team_id, team_name, player1, player2, player3, player4, player5, role, password FROM teams WHERE team_id = $1`
	row := t.db.QueryRow(query, teamID)

	var team entity.Teams
	err := row.Scan(&team.Team_id, &team.Team_name, &team.Player1, &team.Player2, &team.Player3, &team.Player4, &team.Player5, &team.Role, &team.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Teams{}, errors.New("team not found")
		}
		return entity.Teams{}, err
	}

	return team, nil
}

func (t *tournamentRepository) CreateTournament(tour entity.Tournament) error {
	query := "INSERT INTO tournaments (tournament_name, description, start_date, end_date, teams_count, total_round_number, active, winner_team_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"
	_, err := t.db.Exec(query, tour.TournamentName, tour.Description, tour.StartDate, tour.EndDate, tour.TeamsCount, tour.TotalRoundNumber, tour.StatusActive, tour.WinnerTeamID)
	if err != nil {
		return err
	}

	return nil
}

func (t *tournamentRepository) GetAllTournament(active *bool) ([]entity.Tournament, error) {
	var query string
	if active != nil {
		if *active {
			query = "SELECT * FROM tournaments WHERE active = true ORDER BY tournament_id ASC"
		} else if !*active {
			query = "SELECT * FROM tournaments WHERE active = false ORDER BY tournament_id ASC"
		}
	} else {
		query = "SELECT * FROM tournaments ORDER BY tournament_id ASC"
	}

	var tournaments []entity.Tournament

	rows, err := t.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tournament entity.Tournament
		err := rows.Scan(
			&tournament.TournamentID,
			&tournament.TournamentName,
			&tournament.Description,
			&tournament.StartDate,
			&tournament.EndDate,
			&tournament.TeamsCount,
			&tournament.TotalRoundNumber,
			&tournament.StatusActive,
			&tournament.WinnerTeamID,
			&tournament.CurrentTeamsCount,
			&tournament.CurrentRoundNumber,
		)
		if err != nil {
			return nil, err
		}
		tournaments = append(tournaments, tournament)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tournaments, nil
}

func (t *tournamentRepository) RegisterTeam(tournamentID, teamID int) error {
	tx, err := t.db.Beginx()
	if err != nil {
		return err
	}

	var currentTeamsCount, maxTeamsCount int
	query := `
		SELECT current_teams_count, teams_count
		FROM tournaments
		WHERE tournament_id = $1
		FOR UPDATE
	`
	row := t.db.QueryRow(query, tournamentID)
	err = row.Scan(&currentTeamsCount, &maxTeamsCount)
	if err != nil {
		tx.Rollback()
		return err
	}

	if currentTeamsCount >= maxTeamsCount {
		tx.Rollback()
		return fmt.Errorf("tournament registration is full")
	}

	_, err = tx.Exec("INSERT INTO participants_team (team_id, tournament_id) VALUES ($1, $2)", teamID, tournamentID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("You're already registered")
	}

	_, err = tx.Exec("UPDATE tournaments SET current_teams_count = current_teams_count + 1 WHERE tournament_id = $1", tournamentID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (t *tournamentRepository) UnregisterTeam(tournamentID, teamID int) error {
	tx, err := t.db.Beginx()
	if err != nil {
		return err
	}

	tour, err := t.GetTournamentByID(tournamentID)
	if err != nil {
		return err
	}

	if tour.CurrentRoundNumber != 0 {
		return fmt.Errorf("You can't unregister cause tournament was started")
	}

	query := "DELETE FROM participants_team WHERE team_id = $1 AND tournament_id = $2"
	result, err := tx.Exec(query, teamID, tournamentID)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("UPDATE tournaments SET current_teams_count = current_teams_count - 1 WHERE tournament_id = $1", tournamentID)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (t *tournamentRepository) GetAllParticipiant(teamID *int) ([]entity.Participant, error) {
	var participiants []entity.Participant
	var haveID bool
	if teamID != nil {
		haveID = true
	}

	query := "SELECT * FROM participants_team"
	if haveID {
		query = "SELECT * FROM participants_team WHERE team_id=$1"
	}
	rows, err := t.db.Query(query, *teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var participant entity.Participant
		err := rows.Scan(
			&participant.TeamID,
			&participant.TournamentID,
		)
		if err != nil {
			return nil, err
		}
		participiants = append(participiants, participant)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return participiants, nil
}

func (t *tournamentRepository) GetAllParticipiantByTourID(tournamentID int) ([]int, error) {
	var allParticipantID []int

	query := "SELECT team_id FROM participants_team WHERE tournament_id = $1"

	rows, err := t.db.Query(query, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var participantID int
		if err := rows.Scan(&participantID); err != nil {
			return nil, err
		}
		allParticipantID = append(allParticipantID, participantID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return allParticipantID, nil
}

func (t *tournamentRepository) ActiveTournament(touramentID int) error {
	query := "UPDATE tournaments SET active = true WHERE tournament_id = $1"
	_, err := t.db.Exec(query, touramentID)
	return err
}

func (t *tournamentRepository) InactiveTournament(touramentID int) error {
	query := "UPDATE tournaments SET active = false WHERE tournament_id = $1"
	_, err := t.db.Exec(query, touramentID)
	return err
}

func (t *tournamentRepository) GetTournamentByID(tournamentID int) (entity.Tournament, error) {
	var tournament entity.Tournament
	query := "SELECT * FROM tournaments WHERE tournament_id = $1"
	row := t.db.QueryRow(query, tournamentID)
	if err := row.Scan(
		&tournament.TournamentID,
		&tournament.TournamentName,
		&tournament.Description,
		&tournament.StartDate,
		&tournament.EndDate,
		&tournament.TeamsCount,
		&tournament.TotalRoundNumber,
		&tournament.StatusActive,
		&tournament.WinnerTeamID,
		&tournament.CurrentTeamsCount,
		&tournament.CurrentRoundNumber,
	); err != nil {
		return entity.Tournament{}, err
	}

	return tournament, nil
}

// Функция для создания нового матча
func (r *tournamentRepository) CreateMatch(match entity.Match) error {
	query := `
		INSERT INTO matches (tournament_id, round_number, participant1_id, participant2_id, winner_id, loser_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(query, match.TournamentID, match.RoundNumber, match.FirstParticipantID, match.SecondParticipantID, match.WinnerID, match.LoserID)
	if err != nil {
		return err
	}
	return nil
}

func (r *tournamentRepository) GetAllMatchesByTournamentID(tournamentID int) ([]entity.Match, error) {
	var matches []entity.Match
	query := `
		SELECT match_id, tournament_id, round_number, participant1_id, participant2_id, winner_id, loser_id
		FROM matches
		WHERE tournament_id = $1
		ORDER BY match_id ASC
	`
	rows, err := r.db.Query(query, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var match entity.Match
		err := rows.Scan(&match.MatchID, &match.TournamentID, &match.RoundNumber, &match.FirstParticipantID, &match.SecondParticipantID, &match.WinnerID, &match.LoserID)
		if err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

func (t *tournamentRepository) UpdateWinnerMatch(matchID, winnerID, loserID int) error {
	query := "UPDATE matches SET winner_id = $1, loser_id = $2 WHERE match_id = $3"

	_, err := t.db.Exec(query, winnerID, loserID, matchID)
	return err
}

func (t *tournamentRepository) GetAllMatchesByTournamentIDWithName(tournamentID int) ([]entity.Match, error) {
	query := "SELECT m.match_id, m.tournament_id, m.round_number, m.participant1_id, m.participant2_id, m.winner_id, m.loser_id, t1.team_id, t1.team_name, t2.team_id, t2.team_name, tr.tournament_id, tr.tournament_name, tr.description, tr.start_date, tr.end_date, tr.teams_count, tr.total_round_number, tr.active, tr.winner_team_id, tr.current_teams_count, tr.current_round_number FROM matches AS m INNER JOIN teams AS t1 ON m.participant1_id = t1.team_id INNER JOIN teams AS t2 ON m.participant2_id = t2.team_id INNER JOIN tournaments AS tr ON m.tournament_id = tr.tournament_id WHERE m.tournament_id = $1 ORDER BY match_id ASC"
	rows, err := t.db.Query(query, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := []entity.Match{}
	for rows.Next() {
		match := entity.Match{}
		firstParticipant := entity.Teams{}
		secondParticipant := entity.Teams{}
		tournament := entity.Tournament{}

		err := rows.Scan(
			&match.MatchID, &match.TournamentID, &match.RoundNumber,
			&match.FirstParticipantID, &match.SecondParticipantID, &match.WinnerID, &match.LoserID,
			&firstParticipant.Team_id, &firstParticipant.Team_name,
			&secondParticipant.Team_id, &secondParticipant.Team_name,
			&tournament.TournamentID, &tournament.TournamentName, &tournament.Description,
			&tournament.StartDate, &tournament.EndDate, &tournament.TeamsCount,
			&tournament.TotalRoundNumber, &tournament.StatusActive, &tournament.WinnerTeamID,
			&tournament.CurrentTeamsCount, &tournament.CurrentRoundNumber,
		)
		if err != nil {
			return nil, err
		}

		match.FirstParticipant = firstParticipant
		match.SecondParticipant = secondParticipant
		match.Tournament = tournament

		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

func (r *tournamentRepository) GetMatchesByTournamentID(tournamentID int) ([]entity.Match, error) {
	query := "SELECT match_id, tournament_id, round_number, participant1_id, participant2_id, winner_id, loser_id FROM matches WHERE tournament_id = $1 ORDER BY match_id ASC"

	rows, err := r.db.Query(query, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := []entity.Match{}
	for rows.Next() {
		var match entity.Match
		err := rows.Scan(
			&match.MatchID,
			&match.TournamentID,
			&match.RoundNumber,
			&match.FirstParticipantID,
			&match.SecondParticipantID,
			&match.WinnerID,
			&match.LoserID,
		)
		if err != nil {
			return nil, err
		}

		matches = append(matches, match)
	}

	return matches, nil
}

func (t *tournamentRepository) GetWinnersByRoundTournament(roundNumber, tournamentID int) ([]int, error) {
	query := "SELECT winner_id FROM matches WHERE round_number = $1 AND tournament_id = $2"
	rows, err := t.db.Query(query, roundNumber, tournamentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Извлечение идентификаторов победителей
	var winners []int
	for rows.Next() {
		var winnerID int
		err := rows.Scan(&winnerID)
		if err != nil {
			return nil, err
		}
		winners = append(winners, winnerID)
	}

	return winners, nil
}

func (t *tournamentRepository) CheckWinnerLoserIDsNotEmpty(roundNumber, tournamentID int) (bool, error) {
	query := "SELECT winner_id, loser_id FROM matches WHERE round_number = $1 AND tournament_id = $2"
	rows, err := t.db.Query(query, roundNumber, tournamentID)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		var winnerID, loserID sql.NullInt64
		err := rows.Scan(&winnerID, &loserID)
		if err != nil {
			return false, err
		}

		if (!winnerID.Valid || winnerID.Int64 == 0) && (!loserID.Valid || loserID.Int64 == 0) {
			return false, nil
		}

		found = true
	}

	if err = rows.Err(); err != nil {
		return false, err
	}

	if !found {
		return false, nil
	}

	return true, nil
}

func (t *tournamentRepository) UpdateCurrentRoundNumber(tournamentID int) error {
	query := "UPDATE tournaments SET current_round_number = current_round_number + 1 WHERE tournament_id = $1"
	_, err := t.db.Exec(query, tournamentID)
	if err != nil {
		return err
	}

	return nil
}

// func (t *tournamentRepository) GetTournamentByMatchID(matchID int) (entity.Tournament, error) {
// 	return entity.Tournament{}, nil
// }

func (t *tournamentRepository) GetMatchByID(matchID int) (entity.Match, error) {
	var match entity.Match

	// Подготовка SQL-запроса
	query := "SELECT match_id, tournament_id, round_number, participant1_id, participant2_id, winner_id, loser_id FROM matches WHERE match_id = $1"
	row := t.db.QueryRow(query, matchID)

	// Извлечение данных из результирующей строки
	err := row.Scan(&match.MatchID, &match.TournamentID, &match.RoundNumber, &match.FirstParticipantID, &match.SecondParticipantID, &match.WinnerID, &match.LoserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return match, fmt.Errorf("матч с ID %d не найден", matchID)
		}
		return match, err
	}

	return match, nil
}

func (t *tournamentRepository) TournamentWinner(teamID, tournamentID int) error {
	query := "UPDATE tournaments SET winner_team_id = $1 WHERE tournament_id = $2"
	_, err := t.db.Exec(query, teamID, tournamentID)
	if err != nil {
		return err
	}
	return nil
}

func (t *tournamentRepository) GetMatchForNextLevel(tournamentID, roundNumber, participant1, participant2 int) (entity.Match, error) {
	query := `
		SELECT match_id, tournament_id, round_number, participant1_id, participant2_id, winner_id, loser_id
		FROM matches
		WHERE tournament_id = $1 AND round_number = $2 AND (participant1_id = $3 OR participant2_id = $4)
		LIMIT 1`

	var match entity.Match
	err := t.db.QueryRow(query, tournamentID, roundNumber, participant1, participant2).
		Scan(&match.MatchID, &match.TournamentID, &match.RoundNumber, &match.FirstParticipantID, &match.SecondParticipantID, &match.WinnerID, &match.LoserID)
	if err != nil {
		return entity.Match{}, err
	}

	return match, nil
}
