package apis

import (
	echo2 "github.com/labstack/echo/v4"
	"net/http"
)

const VERSION = "0.0.11"

func Accessible(c echo2.Context) error {
	return c.HTML(http.StatusOK, "<h1>ControlCenter</h1>(´・ω・`) 运行正常<br><hr>Ver: "+VERSION)
}
