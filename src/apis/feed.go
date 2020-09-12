package apis

import (
	echo2 "github.com/labstack/echo/v4"
	"net/http"
	"rssreader/src/model"
)

func FeedTodo(c echo2.Context) error {
	return c.JSON(http.StatusOK, model.Response{Code: 0, Message: "TODO"})
}
