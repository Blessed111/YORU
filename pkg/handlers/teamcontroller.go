package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"Diploma/entity"
	"Diploma/pkg/cfg"

	"github.com/gin-gonic/gin"
)

func (H *Handlers) TournamentInfo(c *gin.Context) {
	session, _ := cfg.Store.Get(c.Request, cfg.SESSION_ID)
	//var team models.TeamModel
	if c.Request.Method == http.MethodGet {
		data := gin.H{
			"isnotAuthorized": session.Values["loggedIn"] != true,
			"teamname":        session.Values["teamname"],
			"Team_id":         session.Values["id"],
			"role":            session.Values["role"],
		}
		c.HTML(http.StatusOK, "teampage.html", data)

	}
}

func (H *Handlers) aboutPage(c *gin.Context) {
	session, _ := cfg.Store.Get(c.Request, cfg.SESSION_ID)
	//var team models.TeamModel
	if c.Request.Method == http.MethodGet {
		data := gin.H{
			"isnotAuthorized": session.Values["loggedIn"] != true,
			"teamname":        session.Values["teamname"],
			"Team_id":         session.Values["id"],
			"role":            session.Values["role"],
		}
		c.HTML(http.StatusOK, "about.html", data)
		
	}
}

func (h *Handlers) teamProfile(c *gin.Context) {
	session, _ := cfg.Store.Get(c.Request, cfg.SESSION_ID)
	if c.Request.Method == http.MethodGet {
		id, _ := strconv.ParseInt(c.Query("id"), 10, 64)

		var team entity.Teams
		teamModel.Find(id, &team)

		data := gin.H{

			"isnotAuthorized": session.Values["loggedIn"] != true,
			//"role":            session.Values["role"],

			"team": team,
		}

		c.HTML(http.StatusOK, "profile.html", data)
	}
}

func (h *Handlers) otherTeamProfile(c *gin.Context) {
	session, _ := cfg.Store.Get(c.Request, cfg.SESSION_ID)
	if len(session.Values) == 0 {
		c.Redirect(http.StatusSeeOther, "/auth/logi")
	} else {
		if c.Request.Method == http.MethodGet {
			id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
			session_id := session.Values["id"]
			session_int := session_id.(int)
			var team entity.Teams
			teamModel.Find(id, &team)
			if id != int64(session_int) {
				c.Redirect(http.StatusOK, "/teams/team")
				fmt.Println(session.Values["id"])
				fmt.Println(id)
			}

			data := gin.H{
				"team": team,
			}

			c.HTML(http.StatusOK, "otherTeamPage.html", data)
		}
	}
}

func (h *Handlers) teamProfileUpdate(c *gin.Context) {
	session, _ := cfg.Store.Get(c.Request, cfg.SESSION_ID)
	if len(session.Values) == 0 {
		c.Redirect(http.StatusSeeOther, "/auth/logi")
	} else {
		if c.Request.Method == http.MethodGet {
			id, _ := strconv.ParseInt(c.Query("id"), 10, 64)
			session_id := session.Values["id"]
			session_int := session_id.(int)
			var team entity.Teams
			teamModel.Find(id, &team)
			if id != int64(session_int) {
				fmt.Println("cannot change other teamProfile")
				fmt.Println(session.Values["id"])
				fmt.Println(id)
			} else {
				data := gin.H{
					"team": team,
				}

				c.HTML(http.StatusOK, "profile_update.html", data)
				fmt.Println("same")
				fmt.Println(session.Values["id"])
				fmt.Println(id)
			}
		} else if c.Request.Method == http.MethodPost {
			var team entity.Teams
			id, _ := strconv.Atoi(c.PostForm("id"))
			team.Team_id = int(id)
			team.Team_name = c.PostForm("teamName")
			team.Player1 = c.PostForm("player1")
			team.Player2 = c.PostForm("player2")
			team.Player3 = c.PostForm("player3")
			team.Player4 = c.PostForm("player4")
			team.Player5 = c.PostForm("player5")
			// errorMessages := validation.Struct(team)
			// fmt.Println("validation error")

			teamModel.UpdateProfile(team)
			// c.JSON(http.StatusOK, gin.H{
			// 	"user1": "User data successfully updated",
			// })
			// data := gin.H{
			// 	"team1": "Team updated successfully",
			// }
			// c.HTML(http.StatusOK, "profile_update.html", data)
			c.Redirect(http.StatusFound, "/home")
		}

	}
}
