package service

// import (
// 	"time"

// 	"Diploma"

// 	"github.com/claytonfinney/challonge-go"
// 	"github.com/sirupsen/logrus"
// )

// type TourService struct {
// 	// repos  repository.Tournament
// 	client *challonge.ChallongeClient
// }

// //	func NewTourService(repo repository.Tournament, client *challonge.ChallongeClient) *TourService {
// //		return &TourService{repo, client}
// //	}
// func (t *TourService) CreateTournament(tournament *Diploma.Tournament) error {
// 	tme, _ := time.Parse(time.RFC3339, "2020-05-13T19:45:07.000Z")
// 	// add to tme tournament.Time
// 	tme = tme.Add(tournament.Time)
// 	trn := challonge.Tournament{
// 		challonge.TournamentKey{
// 			Name:                tournament.TournamentName,
// 			TournamentType:      "single elimination",
// 			Description:         tournament.Description,
// 			OpenSignup:          true,
// 			HoldThirdPlaceMatch: true,
// 			Private:             false,
// 			StartAt:             tme,
// 			AcceptAttachments:   true,
// 		},
// 	}
// 	_, err := t.client.CreateTournament(&trn)
// 	if err != nil {
// 		logrus.Fatal(err)
// 	}
// 	return nil
// }

// func (t *TourService) UpdateTournament(tournament *Diploma.Tournament) (*challonge.Tournament, error) {
// 	trn, err := t.client.UpdateTournament(tournament.Url, &challonge.Tournament{
// 		challonge.TournamentKey{
// 			Name:                tournament.TournamentName,
// 			TournamentType:      "single elimination",
// 			Description:         tournament.Description,
// 			OpenSignup:          true,
// 			HoldThirdPlaceMatch: true,
// 			Private:             false,
// 			StartAt:             time.Now(),
// 			AcceptAttachments:   true,
// 		},
// 	})
// 	if err != nil {
// 		logrus.Fatal(err)
// 	}
// 	return trn, nil
// }

// func (t *TourService) DeleteTournament(tournament *Diploma.Tournament) error {
// 	if _, err := t.client.DeleteTournament(tournament.Url); err != nil {
// 		logrus.Fatal(err)
// 	}
// 	return nil
// }

// func (t *TourService) GetTournament(tournament *Diploma.Tournament) (int, error) {
// 	return 0, nil
// }

// func (t *TourService) AppendTour(teams Diploma.Teams) (int, error) {
// 	return 0, nil
// }
