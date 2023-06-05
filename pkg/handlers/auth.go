package handlers

import (
	"errors"
	"log"
	"net/http"

	"Diploma/entity"
	"Diploma/pkg/cfg"
	"Diploma/pkg/libraries"
	"Diploma/pkg/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	validation = libraries.NewValidation()
	teamModel  = models.NewTeamModel()
)

// func (h *Handlers) signUp(c *gin.Context) {
// 	var Input Diploma.Teams
// 	if err := c.BindJSON(&Input); err != nil {
// 		newErrorresponse(c, http.StatusBadRequest, err.Error())
// 		return
// 	}

// 	id, err := h.services.Authorization.CreateTeam(&Input)
// 	if err != nil {
// 		newErrorresponse(c, http.StatusInternalServerError, err.Error())
// 	}
// 	c.JSON(http.StatusOK, map[string]interface{}{
// 		"id": id,
// 	})
// }

type signInInput struct {
	Team_name string `validate:"required"`
	Password  string `validate:"required"`
}

type signInInput1 struct {
	Team_name string `json:"teamName" validate:"required"`
	Password  string `json:"password" validate:"required"`
}

func (h *Handlers) signIn1(c *gin.Context) {
	var Input signInInput
	if err := c.BindJSON(&Input); err != nil {
		newErrorresponse(c, http.StatusBadRequest, err.Error())
		return
	}
	token, err := h.services.Authorization.GenerateToken(Input.Team_name, Input.Password)
	if err != nil {
		newErrorresponse(c, http.StatusInternalServerError, err.Error())
	}
	// cookie, err := c.Request.Cookie("Authorization")
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
	// c.HTML(http.StatusOK, "login.html", gin.H{
	// 	"token": token,
	// })
	return
}

func (h *Handlers) signIn(c *gin.Context) {
	if c.Request.Method == "GET" {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString != "" {
			// c.Redirect(http.StatusMovedPermanently, "/home")
			return
		}
		c.HTML(http.StatusOK, "login.html", nil)
	} else if c.Request.Method == "POST" {
		TeamInput := &signInInput{
			Team_name: c.PostForm("teamName"),
			Password:  c.PostForm("password"),
		}
		// if err := c.BindJSON(&TeamInput); err != nil {
		// 	newErrorresponse(c, http.StatusBadRequest, err.Error())
		// 	return
		// }
		errorMessages := validation.Struct(TeamInput)

		if errorMessages != nil {

			data := gin.H{
				"validation": errorMessages,
			}

			c.HTML(http.StatusBadRequest, "login.html", data)

		} else {
			if TeamInput.Team_name == "admin" && TeamInput.Password == "admin" {
				c.Redirect(http.StatusSeeOther, "/admin/teams")
				return
			}
			token, err := h.services.Authorization.GenerateToken(TeamInput.Team_name, TeamInput.Password)
			var message error
			if err != nil {
				// newErrorresponse(c, http.StatusInternalServerError, err.Error())
				message = errors.New("Wrong teamname or password!")
				log.Println("wrong")
				if message != nil {

					data := gin.H{
						"error": message,
					}

					c.HTML(http.StatusBadRequest, "login.html", data)
					// c.Redirect(http.StatusMovedPermanently, "/auth/sign-in")
				} else {
					c.JSON(http.StatusOK, map[string]interface{}{
						"token": token,
					})
					// cookie, err := c.Request.Cookie("Authorization")

					data := gin.H{
						"token":    token,
						"loggedIn": true,
					}

					c.HTML(http.StatusBadRequest, "login.html", data)
					c.Redirect(http.StatusSeeOther, "/")
					return

				}
			}

		}
	}
}

func (h *Handlers) LoginGETHandler(c *gin.Context) {
	if c.Request.Method == http.MethodGet {
		// session, _ := cfg.Store.Get(c.Request, cfg.SESSION_ID)
		// if session.Values["loggedIn"] == true {
		// 	c.Redirect(http.StatusMovedPermanently, "/")
		// 	return
		// }
		// log.Println("blya")
		c.HTML(http.StatusOK, "login.html", nil)
	}
}

// func (h *Handlers) Logi(c *gin.Context) {
// 	if c.Request.Method == http.MethodGet {
// 		session, _ := cfg.Store.Get(c.Request, cfg.SESSION_ID)
// 		if len(session.Values) != 0 {
// 			c.Redirect(http.StatusSeeOther, "/home")
// 		}
// 		// if session.Values["loggedIn"] == true {
// 		// 	log.Println("already logged in")
// 		// 	c.Redirect(http.StatusSeeOther, "/home")

// 		// 	return
// 		// }
// 		c.HTML(http.StatusOK, "login.html", nil)
// 	} else if c.Request.Method == http.MethodPost {
// 		TeamInput := &signInInput1{
// 			Team_name: c.PostForm("teamName"),
// 			Password:  c.PostForm("password"),
// 		}

// 		errorMessages := validation.Struct(TeamInput)
// 		if errorMessages != nil {
// 			data := gin.H{
// 				"validation": errorMessages,
// 			}
// 			c.HTML(http.StatusOK, "register.html", data)
// 		} else {
// 			if TeamInput.Team_name == "admin" && TeamInput.Password == "admin" {
// 				c.Redirect(http.StatusSeeOther, "/admin/teams")
// 			}
// 			var team Diploma.Teams
// 			teamModel.Where(&team, "team_name", TeamInput.Team_name)

// 			var message error

// 			errPassword := bcrypt.CompareHashAndPassword([]byte(team.Password), []byte(TeamInput.Password))
// 			if errPassword != nil {
// 				message = errors.New("Wrong teamname or password!")
// 				log.Println("wrong password")
// 			}
// 			if errorMessages != nil {
// 				data := gin.H{
// 					"error": message,
// 				}
// 				c.HTML(http.StatusOK, "register.html", data)
// 			} else {
// 				// set session
// 				session, _ := cfg.Store.Get(c.Request, cfg.SESSION_ID)

// 				session.Values["loggedIn"] = true
// 				// session.Values["email"] = user.Email
// 				session.Values["teamname"] = team.Team_name
// 				session.Values["player1"] = team.Player1
// 				session.Values["player2"] = team.Player2
// 				session.Values["player3"] = team.Player3
// 				session.Values["player4"] = team.Player4
// 				session.Values["player5"] = team.Player5
// 				// session.Values["full_name"] = user.Full_name
// 				// session.Values["role"] = user.Role
// 				session.Values["id"] = team.Team_id
// 				session.Save(c.Request, c.Writer)

// 				c.JSON(http.StatusOK, gin.H{
// 					//"session":  session,

// 				})
// 				c.Abort()
// 				data := gin.H{
// 					"loggedIn": true,
// 					"user":     session.Values["teamname"],
// 					"user1":    "team logged",
// 				}
// 				c.HTML(http.StatusOK, "register.html", data)

// 			}
// 		}
// 	}
// }

// func (h *Handlers) Logi1(c *gin.Context) {
// 	if c.Request.Method == http.MethodGet {
// 		session, _ := cfg.Store.Get(c.Request, cfg.SESSION_ID)
// 		if len(session.Values) != 0 {
// 			c.Redirect(http.StatusSeeOther, "/home")
// 		}
// 		// if session.Values["loggedIn"] == true {
// 		// 	log.Println("already logged in")
// 		// 	c.Redirect(http.StatusSeeOther, "/home")

// 		// 	return
// 		// }
// 		c.HTML(http.StatusOK, "login.html", nil)
// 	} else if c.Request.Method == http.MethodPost {
// 		TeamInput := &signInInput1{
// 			Team_name: c.PostForm("teamName"),
// 			Password:  c.PostForm("password"),
// 		}

// 		errorMessages := validation.Struct(TeamInput)

// 		if errorMessages != nil {

// 			data := gin.H{
// 				"validation": errorMessages,
// 			}

// 			c.HTML(http.StatusBadRequest, "login.html", data)

// 		} else {
// 			if TeamInput.Team_name == "admin" && TeamInput.Password == "admin" {
// 				c.Redirect(http.StatusSeeOther, "/admin/teams")
// 			}
// 			var team Diploma.Teams
// 			teamModel.Where(&team, "team_name", TeamInput.Team_name)

// 			var message error

// 			errPassword := bcrypt.CompareHashAndPassword([]byte(team.Password), []byte(TeamInput.Password))
// 			if errPassword != nil {
// 				message = errors.New("Wrong teamname or password!")
// 				log.Println("wrong password")
// 			}
// 			if message != nil {

// 				data := gin.H{
// 					"error": message,
// 				}

// 				c.HTML(http.StatusUnauthorized, "login.html", data)
// 			} else {
// 				// set session
// 				session, _ := cfg.Store.Get(c.Request, cfg.SESSION_ID)

// 				session.Values["loggedIn"] = true
// 				// session.Values["email"] = user.Email
// 				session.Values["teamname"] = team.Team_name
// 				session.Values["player1"] = team.Player1
// 				session.Values["player2"] = team.Player2
// 				session.Values["player3"] = team.Player3
// 				session.Values["player4"] = team.Player4
// 				session.Values["player5"] = team.Player5
// 				// session.Values["full_name"] = user.Full_name
// 				// session.Values["role"] = user.Role
// 				session.Values["id"] = team.Team_id
// 				session.Save(c.Request, c.Writer)

// 				c.Redirect(http.StatusSeeOther, "/home")
// 				// c.HTML(http.StatusOK, "index.html", gin.H{})

// 				// c.JSON(http.StatusOK, gin.H{
// 				// 	"session": session,
// 				// })

// 			}
// 		}
// 	}
// }

func (h *Handlers) Register(c *gin.Context) {
	if c.Request.Method == http.MethodGet {

		c.HTML(http.StatusOK, "register.html", nil)

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
		}

		errorMessages := validation.Struct(team)

		if errorMessages != nil {
			data := gin.H{
				"validation": errorMessages,
				"team":       team,
			}
			c.HTML(http.StatusBadRequest, "register.html", data)
			return
		} else {
			// Hash password
			// hash := sha1.New()
			// hash.Write([]byte(team.Password))

			hashPassword, _ := bcrypt.GenerateFromPassword([]byte(team.Password), bcrypt.DefaultCost)
			team.Password = string(hashPassword)
			// Create user in the database
			//h.services.Authorization.CreateTeam111(&team)
			teamModel.CreateTeam111(team)

			// c.JSON(http.StatusOK, gin.H{
			// 	"team1":  "Team added successfully",
			// 	"status": "success",
			// 	"teams":  team,
			// })
			c.Redirect(http.StatusMovedPermanently, "/auth/logi")
		}
	}

}

func (h *Handlers) Logout(c *gin.Context) {
	session, _ := cfg.Store.Get(c.Request, cfg.SESSION_ID)
	// delete session
	session.Options.MaxAge = -1
	err := session.Save(c.Request, c.Writer)
	if err != nil {
		log.Println("session not saved")
	}
	c.Redirect(http.StatusSeeOther, "/auth/logi")
}
