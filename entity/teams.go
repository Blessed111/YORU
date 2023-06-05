package entity

type Teams struct {
	Team_id   int    `json:"-" db:"team_id"`
	Team_name string `json:"teamName" validate:"required,isunique=teams-team_name"`
	Player1   string `json:"player1" validate:"required,gte=3"`
	Player2   string `json:"player2" validate:"required,gte=3"`
	Player3   string `json:"player3" validate:"required,gte=3"`
	Player4   string `json:"player4" validate:"required,gte=3"`
	Player5   string `json:"player5" validate:"required,gte=3"`
	Password  string `json:"password" validate:"required"`
	Cpassword string `json:"cpassword" validate:"required,eqfield=Password" label:"Confirm Password"`
	Role      string `json:"role"`
}

type Participant struct {
	TeamID       int `json:"team_id"`
	TournamentID int `json:"tournament_id"`
}

// LEGACY
// type Teams struct {
// 	Team_id   int    `json:"-" db:"team_id"`
// 	Team_name string `json:"teamName" validate:"required"`
// 	Player1   string `json:"player1" validate:"required,gte=3"`
// 	Player2   string `json:"player2" validate:"required,gte=3"`
// 	Player3   string `json:"player3" validate:"required,gte=3"`
// 	Player4   string `json:"player4" validate:"required,gte=3"`
// 	Player5   string `json:"player5" validate:"required,gte=3"`
// 	Password  string `json:"password" validate:"required"`
// 	Cpassword string `json:"cpassword" validate:"required,eqfield=Password" label:"Confirm Password"`
// 	Role      string `json:"role"`
// }
