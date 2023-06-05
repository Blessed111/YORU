package handlers

import (
	"log"
	"net/http"

	"Diploma/entity"
	"Diploma/pkg/cfg"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handlers) Logi1(ctx *gin.Context) {
	if ctx.Request.Method == http.MethodGet {
		ctx.HTML(http.StatusUnauthorized, "login.html", nil)
		return
	}
	teamName := ctx.PostForm("teamName")
	password := ctx.PostForm("password")
	var team entity.Teams

	teamModel.Where(&team, "team_name", teamName)
	if err := bcrypt.CompareHashAndPassword([]byte(team.Password), []byte(password)); err != nil {
		ctx.HTML(http.StatusUnauthorized, "login.html", gin.H{
			"error": "Wrong teamname or password!",
		})
		log.Println("wrong password")
	}
	session, err := cfg.Store.Get(ctx.Request, cfg.SESSION_ID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	session.Values["loggedIn"] = true
	session.Values["teamname"] = team.Team_name
	session.Values["player1"] = team.Player1
	session.Values["player2"] = team.Player2
	session.Values["player3"] = team.Player3
	session.Values["player4"] = team.Player4
	session.Values["player5"] = team.Player5
	session.Values["id"] = team.Team_id
	session.Values["role"] = team.Role
	session.Save(ctx.Request, ctx.Writer)

	ctx.Redirect(http.StatusSeeOther, "/home")
}
