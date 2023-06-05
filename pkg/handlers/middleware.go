package handlers

import (
	"net/http"

	"Diploma/pkg/cfg"

	"github.com/gin-gonic/gin"
)

const (
	Authorization = "Authorization"
	teamCtx       = "teamId"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session, err := cfg.Store.Get(ctx.Request, cfg.SESSION_ID)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Проверка авторизации пользователя
		if session.Values["loggedIn"] != true {
			ctx.Redirect(http.StatusSeeOther, "/auth/logi")
			ctx.Abort()
			return
		}

		// Если пользователь авторизован, продолжаем выполнение следующего обработчика
		ctx.Next()
	}
}
func noAccessMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session, err := cfg.Store.Get(ctx.Request, cfg.SESSION_ID)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Проверка роли пользователя
		if session.Values["role"] != "admin" {
			ctx.Redirect(http.StatusSeeOther, "/no_access")
			ctx.Abort()
			return
		}

		// Если пользователь авторизован, продолжаем выполнение следующего обработчика
		ctx.Next()
	}
}

func NonAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		session, err := cfg.Store.Get(ctx.Request, cfg.SESSION_ID)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// Проверка авторизации пользователя
		if session.Values["loggedIn"] == true {
			ctx.Redirect(http.StatusSeeOther, "/home")
			ctx.Abort()
			return
		}

		// Если пользователь авторизован, продолжаем выполнение следующего обработчика
		ctx.Next()
	}
}


