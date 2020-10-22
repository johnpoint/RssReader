package apis

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"rssreader/src/model"
)

func FeedTodo(c echo.Context) error {
	return c.JSON(http.StatusOK, model.Response{Code: 0, Message: "TODO"})
}
