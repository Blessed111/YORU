package entity

type Match struct {
	MatchID             int        `json:"match_id"`
	TournamentID        int        `json:"tournament_id"`
	RoundNumber         int        `json:"round_number"`
	FirstParticipantID  int        `json:"first_participant_id"`
	SecondParticipantID int        `json:"second_participant_id"`
	WinnerID            int        `json:"winner_id"`
	LoserID             int        `json:"loser_id"`
	FirstParticipant    Teams      `json:"first_participant"`
	SecondParticipant   Teams      `json:"second_participant"`
	Tournament          Tournament `json:"tournament"`
}
