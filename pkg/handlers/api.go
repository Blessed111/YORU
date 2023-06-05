package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handlers) CreateBrackets(c *gin.Context) {
	id, _ := c.Get(teamCtx)
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}
