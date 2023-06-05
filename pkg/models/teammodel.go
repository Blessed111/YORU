package models

import (
	"Diploma/entity"
	"Diploma/pkg/repository"
	"fmt"
	"math/rand"
	"net/smtp"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type TeamModel struct {
	conn *sqlx.DB
}

func NewTeamModel() *TeamModel {
	conn, err := repository.NewPostgresDb1()

	if err != nil {
		panic(err)
	}

	return &TeamModel{
		conn: conn,
	}
}

func (u TeamModel) CreateTeam111(team entity.Teams) (int64, error) {

	result, err := u.conn.Exec("insert into teams (team_name, player1, player2, player3, player4, player5, password) values($1,$2,$3,$4,$5,$6,$7)",
		team.Team_name, team.Player1, team.Player2, team.Player3, team.Player4, team.Player5, team.Password)

	if err != nil {
		return 0, err
	}

	lastInsertId, _ := result.LastInsertId()

	return lastInsertId, nil
}

// func (u TeamModel) CreateTeam222(team Diploma.Teamss) (int64, error) {

// 	result, err := u.conn.Exec("insert into teamss (team_email, team_name, player1_nickname, player2_nickname, player3_nickname, player4_nickname, player5_nickname, password) values($1,$2,$3,$4,$5,$6,$7,$8)",
// 		team.Team_email, team.Team_name, team.Player1_nickname, team.Player2_nickname, team.Player3_nickname, team.Player4_nickname, team.Player5_nickname, team.Password)

// 	if err != nil {
// 		return 0, err
// 	}

// 	lastInsertId, _ := result.LastInsertId()

// 	return lastInsertId, nil
// }

func (u TeamModel) Where(team *entity.Teams, fieldName, fieldValue string) error {

	row, err := u.conn.Query("select team_id, team_name, player1, player2, player3, player4, player5, password, role from teams where "+fieldName+" = $1 limit 1;", fieldValue)

	if err != nil {
		return err
	}

	defer row.Close()

	for row.Next() {
		row.Scan(&team.Team_id, &team.Team_name, &team.Player1, &team.Player2, &team.Player3, &team.Player4, &team.Player5, &team.Password, &team.Role)
	}

	return nil
}

func (p *TeamModel) FindAll() ([]entity.Teams, error) {

	rows, err := p.conn.Query("select * from teams")
	if err != nil {
		return []entity.Teams{}, err
	}
	defer rows.Close()

	var dataTeam []entity.Teams
	for rows.Next() {
		var team entity.Teams
		rows.Scan(&team.Team_id,
			&team.Team_name,
			&team.Player1,
			&team.Player2,
			&team.Player3,
			&team.Player4,
			&team.Player5,
			&team.Password,
			&team.Role,
		)

		dataTeam = append(dataTeam, team)
	}

	return dataTeam, nil

}

const teamtable = "Teams"

func (t TeamModel) CreateTeam(team entity.Teams) (int, error) {

	var id int
	query := fmt.Sprintf("INSERT INTO %s (team_name,player1,player2,player3,player4,player5,password) values ($1,$2,$3,$4,$5,$6,$7)", teamtable)
	row := t.conn.QueryRow(query, team.Team_name, team.Player1, team.Player2, team.Player3, team.Player4, team.Player5, team.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (p *TeamModel) Find(id int64, team *entity.Teams) error {

	return p.conn.QueryRow("select team_id, team_name, player1, player2, player3, player4, player5, role from teams where team_id = $1", id).Scan(
		&team.Team_id,
		&team.Team_name,
		&team.Player1,
		&team.Player2,
		&team.Player3,
		&team.Player4,
		&team.Player5,
		&team.Role)
}

func (p *TeamModel) Update(team entity.Teams) error {

	_, err := p.conn.Exec(
		"update teams set team_name = $1, player1 = $2, player2 = $3, player3 = $4, player4 = $5, player5 = $6, role = $7 where team_id = $8",
		team.Team_name, team.Player1, team.Player2, team.Player3, team.Player4, team.Player5, team.Role, team.Team_id)

	if err != nil {
		return err
	}

	return nil
}
func (p *TeamModel) UpdateProfile(team entity.Teams) error {

	_, err := p.conn.Exec(
		"update teams set player1 = $1, player2 = $2, player3 = $3, player4 = $4, player5 = $5 where team_id = $6",
		team.Player1, team.Player2, team.Player3, team.Player4, team.Player5, team.Team_name)

	if err != nil {
		return err
	}

	return nil
}

func (p *TeamModel) Delete(id int64) {
	p.conn.Exec("delete from teams where team_id = $1", id)
}

//TEST

func generatePassword() string {
	// Generate a random 8-character password
	const passwordLength = 8
	const passwordChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	var sb strings.Builder
	for i := 0; i < passwordLength; i++ {
		sb.WriteByte(passwordChars[rand.Intn(len(passwordChars))])
	}
	return sb.String()
}

func hashPassword(password string) string {
	// Hash the password using bcrypt
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func sendNewPasswordEmail(email, password string) error {
	// Set up the email message
	from := "ospan202022@gmail.com"
	to := email
	subject := "New password for your account"
	body := fmt.Sprintf("Your new password is: %s", password)
	msg := []byte("From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body)

	// Set up the SMTP client and send the email
	auth := smtp.PlainAuth("", from, "your-password", "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, msg)
	return err
}
