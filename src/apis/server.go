package apis

import (
	"net/http"

	echo2 "github.com/labstack/echo/v4"
)

const VERSION = "0.16"

func Accessible(c echo2.Context) error {
	return c.HTML(http.StatusOK, "<h1>RssReader api</h1>(´・ω・`) 运行正常<br><hr>Ver: "+VERSION)
}
