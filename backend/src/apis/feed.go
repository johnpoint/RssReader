package apis

import (
	"github.com/johnpoint/rssreader-server/src/model"
	echo2 "github.com/labstack/echo/v4"
	"net/http"
)



func FeedTodo(c echo2.Context) error {
	return c.JSON(http.StatusOK, model.Response{Code: 0, Message: "TODO"})
}
