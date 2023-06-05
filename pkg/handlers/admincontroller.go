package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"Diploma/entity"
	"Diploma/pkg/cfg"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// var model = models.NewTeamModel()
func (H *Handlers) HomePage(c *gin.Context) {
	session, _ := cfg.Store.Get(c.Request, cfg.SESSION_ID)
	fmt.Println("session:", session)
	// var team models.TeamModel
	if c.Request.Method == http.MethodGet {
		data := gin.H{
			"isnotAuthorized": session.Values["loggedIn"] != true,
			"role":            session.Values["role"],
			"teamname":        session.Values["teamname"],
			"Team_id":         session.Values["id"],
			//"team":     team,
		}
		c.HTML(http.StatusOK, "index.html", data)

	}
}

// func (h *Handlers) AddTeam(c *gin.Context) {
// 	if c.Request.Method == http.MethodGet {
// 		c.HTML(http.StatusOK, "admin_add.html", nil)
// 	} else if c.Request.Method == http.MethodPost {
// 		team := Diploma.Teams{
// 			Team_name: c.PostForm("teamName"),
// 			Player1:   c.PostForm("player1"),
// 			Player2:   c.PostForm("player2"),
// 			Player3:   c.PostForm("player3"),
// 			Player4:   c.PostForm("player4"),
// 			Player5:   c.PostForm("player5"),
// 			Password:  c.PostForm("password"),
// 			Role: c.PostForm("role"),

// 		}
// 		if err := c.ShouldBind(&team); err != nil {
// 			c.HTML(http.StatusBadRequest, "admin_add.html", gin.H{"error": err.Error()})
// 			return
// 		}

// 		errorMessages := validation.Struct(team)

// 		if errorMessages != nil {
// 			data := gin.H{
// 				"validation": errorMessages,
// 				"team":       team,
// 			}
// 			c.HTML(http.StatusBadRequest, "admin_add.html", data)
// 			return

// 		} else {

// 			// Hash password
// 			// hash := sha1.New()
// 			// hash.Write([]byte(team.Password))

// 			hashPassword, _ := bcrypt.GenerateFromPassword([]byte(team.Password), bcrypt.DefaultCost)
// 			team.Password = string(hashPassword)
// 			// Create user in the database
// 			teamModel.CreateTeam111(team)

// 			data := gin.H{
// 				"team1": "Team added successfully",
// 			}
// 			c.HTML(http.StatusOK, "admin_add.html", data)

// 			c.JSON(http.StatusOK, gin.H{
// 				"team1":  "Team added successfully",
// 				"status": "success",
// 				"teams":  team,
// 			})
// 			c.Redirect(http.StatusMovedPermanently, "/admin/teams")
// 		}
// 	}
// }

func (h *Handlers) AddTeam(c *gin.Context) {
	// session, _ := cfg.Store.Get(c.Request, cfg.SESSION_ID)
	// if session.Values["role"] != "admin" {
	// 	c.Redirect(http.StatusSeeOther, "/no_access")
	// } else {
	if c.Request.Method == http.MethodGet {

		c.HTML(http.StatusOK, "admin_add.html", nil)

	} else if c.Request.Method == http.MethodPost {
		team := entity.Teams{
			Team_name: c.PostForm("teamName"),
			Player1:   c.PostForm("player1"),
			Player2:   c.PostForm("player2"),
			Player3:   c.PostForm("player3"),
			Player4:   c.PostForm("player4"),
			Player5:   c.PostForm("player5"),
			Password:  c.PostForm("password"),
			Cpassword: c.PostForm("cpassword"),
			Role:      c.PostForm("role"),
		}

		errorMessages := validation.Struct(team)
		fmt.Println("validation error")

		if errorMessages != nil {
			data := gin.H{
				"validation": errorMessages,
				"team":       team,
			}
			c.HTML(http.StatusBadRequest, "admin_add.html", data)
			return
		} else {
			hashPassword, _ := bcrypt.GenerateFromPassword([]byte(team.Password), bcrypt.DefaultCost)
			team.Password = string(hashPassword)

			teamModel.CreateTeam111(team)
			c.Redirect(http.StatusMovedPermanently, "/admin/teams")
			// data := gin.H{
			// 	// "team1": "Team added successfully",
			// }
			// c.HTML(http.StatusOK, "admin_add.html", data)

			// c.JSON(http.StatusOK, gin.H{
			// 	"team1":  "Team added successfully",
			// 	"status": "success",
			// 	"teams":  team,
			// })

		}
		//}
	}

}

func (h *Handlers) getTeamsList(c *gin.Context) {
	session, _ := cfg.Store.Get(c.Request, cfg.SESSION_ID)
	if session.Values["role"] != "admin" {
		c.Redirect(http.StatusSeeOther, "/no_access")
	} else {

		teams, _ := teamModel.FindAll()
		var team entity.Teams
		c.HTML(200, "admin.html", gin.H{
			"teams":    teams,
			"teamname": session.Values["teamname"],
			"Team_id":  session.Values["id"],

			"team": team,
		})
		// c.JSON(http.StatusOK, gin.H{
		// 	"status": "success",
		// 	"teams":  teams,
		// })
	}
}

func (h *Handlers) updateTeam(c *gin.Context) {
	// session, _ := cfg.Store.Get(c.Request, cfg.SESSION_ID)
	// if session.Values["role"] != "admin" {
	// 	c.Redirect(http.StatusSeeOther, "/no_access")
	// } else {
	if c.Request.Method == http.MethodGet {
		id, _ := strconv.ParseInt(c.Query("id"), 10, 64)

		var team entity.Teams
		teamModel.Find(id, &team)

		data := gin.H{
			"team": team,
		}

		c.HTML(http.StatusOK, "admin_update.html", data)
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
		team.Role = c.PostForm("role")

		teamModel.Update(team)

		// c.JSON(http.StatusOK, gin.H{
		// 	"user1": "User data successfully updated",
		// })
		//}

		c.Redirect(http.StatusMovedPermanently, "/admin/teams")

		//}
	}
}

func (h *Handlers) deleteTeam(c *gin.Context) {
	idStr := c.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	teamModel.Delete(id)

	// c.JSON(http.StatusOK, gin.H{
	// 	"status":  "success",
	// 	"message": "User deleted successfully",
	// })

	c.Redirect(http.StatusMovedPermanently, "/admin/teams")
}

// TEST

// func ForgotPassword(c *gin.Context) {
// 	email := c.PostForm("email")

// 	// Check if email exists in the database
// 	team, err := teamModel.GetByEmail(email)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Email not found"})
// 		return
// 	}

// 	// Generate a new password
// 	newPassword := models.generatePassword()

// 	// Update the user's password in the database
// 	team = hashPassword(newPassword)
// 	if err := TeaModel.UpdatePassword(user); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
// 		return
// 	}

// 	// Send the new password to the user's email
// 	if err := sendNewPasswordEmail(user.Email, newPassword); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "New password has been sent to your email"})
// }

func (h *Handlers) noAccess(c *gin.Context) {
	data := gin.H{
		"error": "no rights",
	}
	c.HTML(http.StatusForbidden, "no_access.html", data)
}
