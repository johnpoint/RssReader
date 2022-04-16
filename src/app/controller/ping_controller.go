package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Pong(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{"code": 0, "msg": "pong"})
}
