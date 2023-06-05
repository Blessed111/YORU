package handlers

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"

	"Diploma/entity"
	"Diploma/pkg/cfg"

	"github.com/gin-gonic/gin"
)

// --------------------ADMIN CASE--------------------------------//
func (h *Handlers) allTournaments(ctx *gin.Context) {
	tournaments, err := h.services.TournamentServiceI.GetAllTournament(nil, nil)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	session, err := cfg.Store.Get(ctx.Request, cfg.SESSION_ID)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	ctx.HTML(http.StatusOK, "admin_bracket_new.html", gin.H{
		"Tournament": tournaments,
		"Name":       session.Values["teamname"],
		"ID":         session.Values["id"],
	})
}

// создаем турнир
func (h *Handlers) newTournament(ctx *gin.Context) {
	var tournament entity.Tournament
	startTime, err := time.Parse("2006-01-02", ctx.PostForm("start_date"))
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": http.StatusText(http.StatusBadRequest),
		})
		return
	}

	endTime, err := time.Parse("2006-01-02", ctx.PostForm("end_date"))
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": http.StatusText(http.StatusBadRequest),
		})
		return
	}
	teamCount, err := strconv.Atoi(ctx.PostForm("teams_count"))
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": http.StatusText(http.StatusBadRequest),
		})
		return
	}
	if teamCount < 4 || teamCount > 512 {
		// ctx.JSON(http.StatusBadRequest, gin.H{
		// 	"error": "need minimum 4 team count or maximum 512",
		// })
		ctx.HTML(http.StatusBadRequest, "admin_warning.html", gin.H{
			"error": "need minimum 4 team count or maximum 512",
		})
		return
	}
	active := ctx.PostForm("status_active")
	if active == "on" {
		tournament.StatusActive = true
	}

	tournament.TournamentName = ctx.PostForm("tournament_name")
	tournament.Description = ctx.PostForm("description")
	tournament.StartDate = startTime
	tournament.EndDate = endTime
	tournament.TeamsCount = teamCount
	tournament.TotalRoundNumber = int(math.Log2(float64(teamCount)))

	if err := h.services.TournamentServiceI.CreateTournament(tournament); err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": http.StatusText(http.StatusBadRequest),
		})
		return
	}
	ctx.Redirect(http.StatusFound, "/admin/tournament")
}

func (h *Handlers) processTournament(ctx *gin.Context) {
	tournamentID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	matches, err := h.services.TournamentServiceI.GetProcessTour(tournamentID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	tour, err := h.services.GetTournamentByID(tournamentID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	switch ctx.Request.URL.Path {
	case "/teams/tournament/grid/" + ctx.Param("id"):
		ctx.HTML(http.StatusOK, "team_bracket.html", gin.H{
			"TotalRound": tour.TotalRoundNumber,
			"Match":      matches,
		})
		return
	case "/admin/tournament/process/" + ctx.Param("id"):
		ctx.HTML(http.StatusOK, "matches_process.html", matches)
		return

	default:
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"invalid path": ctx.Request.URL.Path,
		})
		return
	}
}

func (h *Handlers) activityTournament(ctx *gin.Context) {
	strID := ctx.PostForm("id")
	tournamentID, err := strconv.Atoi(strID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	strID = ctx.PostForm("round_number")
	roundNumber, err := strconv.Atoi(strID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	switch ctx.Request.URL.Path {
	case "/admin/tournament/start":
		err := h.services.TournamentServiceI.StartTournament(roundNumber, tournamentID)
		if err != nil {
			ctx.HTML(http.StatusBadRequest, "admin_warning.html", gin.H{
				"error": err.Error(),
			})
			return
		}
	case "/admin/tournament/active":
		err := h.services.TournamentServiceI.ActiveTournament(tournamentID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	case "/admin/tournament/inactive":
		err := h.services.TournamentServiceI.InactiveTournament(tournamentID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	case "/admin/tournament/finish":
		err := h.services.TournamentServiceI.FinishTournament(tournamentID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	default:
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error url path": ctx.Request.URL.Path,
		})
		return
	}

	ctx.Redirect(http.StatusFound, "/admin/tournament")
}

func (h *Handlers) updateMatch(ctx *gin.Context) {
	tournamentID := ctx.PostForm("tour_id")

	switch ctx.Request.URL.Path {
	case "/admin/tournament/match":

		matchID, err := strconv.Atoi(ctx.PostForm("match_id"))
		if err != nil {

			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		winnerID, err := strconv.Atoi(ctx.PostForm("winner_id"))
		if err != nil {

			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		loserID, err := strconv.Atoi(ctx.PostForm("loser_id"))
		if err != nil {

			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = h.services.TournamentServiceI.UpdateWinnerMatch(matchID, winnerID, loserID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
	}
	ctx.Redirect(http.StatusFound, fmt.Sprintf("/admin/tournament/process/%s", tournamentID))
}

// -----------------------------USER CASES------------------------------------------------------------//

func (h *Handlers) activeTournaments(ctx *gin.Context) {
	active := true
	session, _ := cfg.Store.Get(ctx.Request, cfg.SESSION_ID)
	teamID := (session.Values["id"]).(int)
	tournaments, err := h.services.TournamentServiceI.GetAllTournament(&active, &teamID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.HTML(http.StatusOK, "team_tour_info.html", gin.H{
		"Tournament": tournaments,
		"Name":       session.Values["teamname"],
		"ID":         session.Values["id"],
	})
}

func (h *Handlers) welcomeList(ctx *gin.Context) {
	active := true
	tournaments, err := h.services.TournamentServiceI.GetAllTournament(&active, nil)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	session, _ := cfg.Store.Get(ctx.Request, cfg.SESSION_ID)
	ctx.HTML(http.StatusOK, "welcome_team_tour_info.html", gin.H{
		"Tournament":      tournaments,
		"isnotAuthorized": session.Values["loggedIn"] != true,
		"teamName":        session.Values["teamname"],
		"Team_id":         session.Values["id"],
	})
}

func (h *Handlers) registerTournament(ctx *gin.Context) {
	strID := ctx.PostForm("id")
	tournamentID, err := strconv.Atoi(strID)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	session, _ := cfg.Store.Get(ctx.Request, cfg.SESSION_ID)
	teamID := (session.Values["id"]).(int)

	switch ctx.Request.URL.Path {
	case "/teams/register":
		if err = h.services.RegisterTeam(tournamentID, teamID); err != nil {
			ctx.HTML(http.StatusBadRequest, "admin_warning.html", gin.H{
				"error": err.Error(),
			})
		}
	case "/teams/unregister":
		if err = h.services.UnregisterTeam(tournamentID, teamID); err != nil {
			ctx.HTML(http.StatusBadRequest, "admin_warning.html", gin.H{
				"error": err.Error(),
			})

		}
	default:
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"invalid path url": ctx.Request.URL.Path,
		})
	}

	ctx.Redirect(http.StatusFound, "/teams/tournaments")
}
