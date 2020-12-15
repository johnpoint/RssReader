package model

import "github.com/labstack/echo/v4"

type SysContext struct {
	echo.Context
	Config string
}
