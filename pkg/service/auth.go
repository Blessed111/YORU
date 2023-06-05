package service

import (
	"Diploma/pkg/repository"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

const (
	salt       = "s1a2l3t4"
	timeTl     = time.Hour * 24
	signingkey = "qwertyd15341543"
)

type TokenClaims struct {
	jwt.StandardClaims
	Team_id int `json:"team_Id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo}
}

// func (s *AuthService) CreateTeam(team *Diploma.Teams) (int, error) {
// 	team.Password = generatePassword(team.Password)
// 	return s.repo.CreateTeam(team)
// }

// func (s *AuthService) CreateTeam111(team *Diploma.Teams) (int, error) {
// 	team.Password = generatePassword(team.Password)
// 	return s.repo.CreateTeam111(team)
// }

// func (s *AuthService) CreateTeamReg(team *Diploma.Teams) (int, error) {
// 	team.Password = generatePassword(team.Password)
// 	return s.repo.CreateTeam(team)
// }

func generatePassword(password string) string {
	// hash := sha1.New()
	// hash.Write([]byte(password))
	// return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	return string(hash)

}
// func (s *AuthService) GenerateToken(teamname, password string) (string, error) {
// 	team, err := s.repo.GetUser(teamname, generatePassword(password))
// 	if err != nil {
// 		return "", err
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
// 		jwt.StandardClaims{
// 			IssuedAt:  time.Now().Unix(),
// 			ExpiresAt: time.Now().Add(timeTl).Unix(),
// 		},
// 		team.Team_id,
// 	})

// 	return token.SignedString([]byte(signingkey))
// }
func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signing key value ")
		}
		return []byte(signingkey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return 0, errors.New("token claims are not correct (type)")
	}
	return claims.Team_id, nil
}
