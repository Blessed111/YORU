package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type error1 struct {
	Message string `json:"message"`
}

func newErrorresponse(c *gin.Context, Statuscode int, message string) {
	logrus.Error(message)
	fmt.Println("invalid charac")
	c.AbortWithStatusJSON(Statuscode, error1{message})
}
