package main

import (
	"github.com/labstack/echo/v4"
)

func (app *application) getUserFromContext(c echo.Context) string {
	userEmail, _ := c.Request().Context().Value(userCtx).(string)
	return userEmail
}
