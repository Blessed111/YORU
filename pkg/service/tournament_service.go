package service

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"Diploma/entity"
	"Diploma/pkg/repository"

	_ "github.com/lib/pq"
)

type TournamentServiceI interface {
	CreateTournament(tour entity.Tournament) error
	GetAllTournament(active *bool, teamID *int) ([]entity.Tournament, error)
	RegisterTeam(tournamentID, teamID int) error
	UnregisterTeam(tournamentID, teamID int) error
	StartTournament(roundNumber, tournamentID int) error
	ActiveTournament(touramentID int) error
	InactiveTournament(touramentID int) error
	UpdateWinnerMatch(matchID, winnerID, loserID int) error
	GetProcessTour(tournamentID int) ([]entity.Match, error)
	FinishTournament(tournamentID int) error
	GetTournamentByID(tournamentID int) (entity.Tournament, error)
}

type tournamentService struct {
	tournamentRepository repository.TournamentRepositoryI
}

func newTournamentService(repo repository.TournamentRepositoryI) TournamentServiceI {
	return &tournamentService{
		tournamentRepository: repo,
	}
}

func (t *tournamentService) GetTournamentByID(tournamentID int) (entity.Tournament, error) {
	return t.tournamentRepository.GetTournamentByID(tournamentID)
}

func (t *tournamentService) CreateTournament(tour entity.Tournament) error {
	// NEED input data validation logic
	count := calculateRounds(tour.TeamsCount)
	tour.TotalRoundNumber = count
	return t.tournamentRepository.CreateTournament(tour)
}

func (t *tournamentService) GetAllTournament(active *bool, teamID *int) ([]entity.Tournament, error) {
	// NEED input data validation logic
	tournaments, err := t.tournamentRepository.GetAllTournament(active)
	if err != nil {
		return nil, err
	}

	if teamID != nil {
		participants, err := t.tournamentRepository.GetAllParticipiant(teamID)
		if err != nil {
			return nil, err
		}
		for _, v := range participants {
			for i, k := range tournaments {
				if v.TournamentID == k.TournamentID {
					tournaments[i].Registered = true
				}
			}
		}
	}

	for i, v := range tournaments {
		if v.WinnerTeamID != 0 {
			team, err := t.tournamentRepository.GetTeamByID(v.WinnerTeamID)
			if err != nil {
				return nil, err
			}
			tournaments[i].WinnerName = team.Team_name
		}
	}

	return tournaments, nil
}

func (t *tournamentService) RegisterTeam(tournamentID, teamID int) error {
	return t.tournamentRepository.RegisterTeam(tournamentID, teamID)
}

func (t *tournamentService) UnregisterTeam(tournamentID, teamID int) error {
	return t.tournamentRepository.UnregisterTeam(tournamentID, teamID)
}

func (t *tournamentService) ActiveTournament(touramentID int) error {
	return t.tournamentRepository.ActiveTournament(touramentID)
}

func (t *tournamentService) InactiveTournament(tournamentID int) error {
	return t.tournamentRepository.InactiveTournament(tournamentID)
}

func splitArrayIntoPairs(arr []int) [][]int {
	var pairs [][]int

	for i := 0; i < len(arr); i += 2 {
		end := i + 2
		if end > len(arr) {
			end = len(arr)
		}

		pair := arr[i:end]
		pairs = append(pairs, pair)
	}

	return pairs
}

func createMatch(roundNumber, tournamentID, participant1ID, participant2ID int) entity.Match {
	match := entity.Match{
		TournamentID:        tournamentID,
		RoundNumber:         roundNumber,
		FirstParticipantID:  participant1ID,
		SecondParticipantID: participant2ID,
	}

	return match
}

// Функция для создания матча и записи его в базу данных
func (t *tournamentService) createAndSaveMatch(roundNumber, tournamentID, participant1ID, participant2ID int) error {
	match := createMatch(roundNumber, tournamentID, participant1ID, participant2ID)
	err := t.tournamentRepository.CreateMatch(match)
	if err != nil {
		return err
	}
	return nil
}

var (
	net     = make(map[int][]int)
	reverse = true
	mu      sync.Mutex
)

func (t *tournamentService) StartTournament(roundNumber, tournamentID int) error {
	if roundNumber == 0 {
		tour, err := t.tournamentRepository.GetTournamentByID(tournamentID)
		if err != nil {
			return err
		}
		if tour.CurrentTeamsCount < 4 {
			return fmt.Errorf("need minimum 4 teams to start")
		}
	}
	var err error

	if roundNumber == 0 {
		err := t.tournamentRepository.UpdateCurrentRoundNumber(tournamentID)
		if err != nil {
			fmt.Println(err)
			return err
		}

	} else {

		status, err := t.tournamentRepository.CheckWinnerLoserIDsNotEmpty(roundNumber, tournamentID)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !status {
			return fmt.Errorf("prev round not finished1")
		}
	}

	var allParticipants []int
	if roundNumber != 0 {
		allParticipants, err = t.tournamentRepository.GetWinnersByRoundTournament(roundNumber, tournamentID)
		if err != nil {
			fmt.Println(err, "1")
			return err
		}
	} else {
		allParticipants, err = t.tournamentRepository.GetAllParticipiantByTourID(tournamentID)
		if err != nil {
			fmt.Println(err, "2")
			return err
		}
	}
	// fmt.Println("ALL PARTICIPANTS", allParticipants)
	// fmt.Println("ALL PARTICIPANTS LEN", len(allParticipants))

	if len(allParticipants) == 1 {

		err := t.tournamentRepository.TournamentWinner(allParticipants[0], tournamentID)
		if err != nil {
			fmt.Println(err, "3")
			return err
		}
		return fmt.Errorf(fmt.Sprintf("tournament was finished winner team id: %d", allParticipants[0]))
	}

	if len(allParticipants) == 0 {
		return fmt.Errorf("prev round not finished")
	}
	if roundNumber != 0 {
		err := t.tournamentRepository.UpdateCurrentRoundNumber(tournamentID)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}

	if roundNumber == 0 {
		shuffleIntSlice(allParticipants)
		net[tournamentID] = allParticipants
	} else {
		// fmt.Println(allParticipants, "ALL")
		// fmt.Println(net[tournamentID], "NETT")
		allParticipantsMap := make(map[int]bool)
		for _, participant := range allParticipants {
			allParticipantsMap[participant] = true
		}

		for tournamentID, participantIDs := range net {
			updatedIDs := make([]int, 0, len(participantIDs))
			for _, participant := range participantIDs {
				if allParticipantsMap[participant] {
					updatedIDs = append(updatedIDs, participant)
				}
			}
			net[tournamentID] = updatedIDs
		}

		if reverse {
			for _, participantIDs := range net {
				for i, j := 0, len(participantIDs)-1; i < j; i, j = i+1, j-1 {
					participantIDs[i], participantIDs[j] = participantIDs[j], participantIDs[i]
				}
			}
			reverse = false
		} else {
			mu.Lock()
			reverse = true
			mu.Unlock()
		}
		// fmt.Println(allParticipants, "ALL")
		// fmt.Println(net[tournamentID], "NETT")

	}

	pairs := splitArrayIntoPairs(net[tournamentID])
	fmt.Println("pairs", pairs)
	for _, v := range pairs {
		if len(v) == 1 {

			// need logic to go to the next round
			err := t.createAndSaveMatch(roundNumber+1, tournamentID, v[0], v[0])
			if err != nil {
				fmt.Println(err, "1")
				return err
			}
			match, err := t.tournamentRepository.GetMatchForNextLevel(tournamentID, roundNumber+1, v[0], v[0])
			if err != nil {
				fmt.Println("207")
				return err
			}
			err = t.UpdateWinnerMatch(match.MatchID, v[0], v[0])
			if err != nil {
				return err
			}
			continue
		} else if len(v) < 2 {
			continue
		} else {

			err := t.createAndSaveMatch(roundNumber+1, tournamentID, v[0], v[1])
			if err != nil {
				fmt.Println(err, "1")
				return err
			}
		}
	}

	return nil
}

func (t *tournamentService) UpdateWinnerMatch(matchID, winnerID, loserID int) error {
	if winnerID == loserID {
		loserID = -1
	}
	match, err := t.tournamentRepository.GetMatchByID(matchID)
	if err != nil {
		return err
	}
	tour, err := t.tournamentRepository.GetTournamentByID(match.TournamentID)
	if err != nil {
		return err
	}

	if tour.CurrentRoundNumber > match.RoundNumber {
		return fmt.Errorf("can't update match because this match was ended")
	}

	return t.tournamentRepository.UpdateWinnerMatch(matchID, winnerID, loserID)
}

func calculateRounds(participants int) int {
	if participants < 2 {
		return 0
	}

	if participants&(participants-1) == 0 {
		rounds := int(math.Log2(float64(participants)))
		return rounds
	}
	return int(math.Ceil(math.Log2(float64(participants))))
}

func (t *tournamentService) GetProcessTour(tournamentID int) ([]entity.Match, error) {
	matches, err := t.tournamentRepository.GetAllMatchesByTournamentIDWithName(tournamentID)
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func (t *tournamentService) FinishTournament(tournamentID int) error {
	tour, err := t.tournamentRepository.GetTournamentByID(tournamentID)
	if err != nil {
		return err
	}
	if tour.CurrentRoundNumber > tour.TotalRoundNumber {
		return nil
	}
	if tour.WinnerTeamID != 0 {
		return t.tournamentRepository.UpdateCurrentRoundNumber(tournamentID)
	}

	return fmt.Errorf("tournament hasn't winner. use inactive ")
}

func shuffleIntSlice(slice []int) {
	rand.Seed(time.Now().UnixNano())

	for i := len(slice) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		(slice)[i], (slice)[j] = (slice)[j], (slice)[i]
	}
}
