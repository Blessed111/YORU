package handlers

import (
	"net/http"

	"Diploma/pkg/cfg"
	"Diploma/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	services *service.Service
}

func NewHandler(service *service.Service) *Handlers {
	return &Handlers{services: service}
}

func (h *Handlers) Init() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("../templates/*")
	router.Static("/assets", "../assets")
	router.GET("/home", h.HomePage)
	router.GET("/about", h.aboutPage)
	router.GET("/no_access", h.noAccess)

	team := router.Group("/teams")
	//router.Use(AuthMiddleware())
	router.GET("/tournament-list", h.welcomeList)
	team.Use(AuthMiddleware())
	{
		team.GET("/team/profile", h.teamProfile)
		team.GET("/team", h.otherTeamProfile)
		team.GET("/tournaments", h.activeTournaments)
		team.POST("/register", h.registerTournament)
		team.POST("/unregister", h.registerTournament)
	}
	router.GET("/profile/update", h.teamProfileUpdate)
	router.POST("/profile/update", h.teamProfileUpdate)

	router.GET("/teams/tournament/grid/:id", h.processTournament)

	auth := router.Group("/auth")
	auth.Use(NonAuthMiddleware())
	{
		auth.POST("/sign-up", h.Register)
		auth.GET("/sign-up", h.Register)
		auth.POST("/logi", h.Logi1)
		auth.GET("/logi", h.Logi1)
	}
	router.GET("/lout", h.Logout)

	//ADMIN
	admin := router.Group("/admin")
	admin.Use(noAccessMiddleware())
	{
		admin.GET("/teams", h.getTeamsList)
		admin.GET("/teams/delete", h.deleteTeam)
		admin.POST("/teams/add", h.AddTeam)
		admin.GET("/teams/add", h.AddTeam)
		admin.POST("/teams/update", h.updateTeam)
		admin.GET("/teams/update", h.updateTeam)

		// admin.GET("/tournament", func(c *gin.Context) {
		// 	c.HTML(http.StatusOK, "admin_bracket.html",
		// 		nil,
		// 	)
		// })
		admin.GET("/tournament/new", func(ctx *gin.Context) {
			session, err := cfg.Store.Get(ctx.Request, cfg.SESSION_ID)
			if err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			ctx.HTML(http.StatusOK, "form.html", gin.H{
				"Name": session.Values["teamname"],
				"ID":   session.Values["id"],
			})
		})
		admin.POST("/tournament/start", h.activityTournament)
		admin.POST("/tournament/finish", h.activityTournament)
		admin.POST("/tournament/inactive", h.activityTournament)
		admin.POST("/tournament/active", h.activityTournament)
		admin.GET("/tournament/process/:id", h.processTournament)

		admin.GET("/tournament", h.allTournaments)
		admin.POST("/create-tournament", h.newTournament)

		admin.POST("/tournament/match", h.updateMatch)
	}
	return router
}
