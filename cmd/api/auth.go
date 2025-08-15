package main

import (
	"devops/internal/auth"
	"devops/internal/store"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

// @Summary Auth getAuthCallback
// @Description Callback for authentication
// @Tags Auth
// @Accept json
// @Produce json
// @Param provider path string true "Provider name"
// @Success 200 {object} object
// @Failure 500 {object} error
// @Router /auth/{provider} [get]
func (app *application) startAuthHandler(c echo.Context) error {
	provider := c.Param("provider")
	auth.BeginAuth(c.Response(), c.Request(), provider)
	return nil
}

// @Summary Auth start auth handler
// @Description Callback for authentication
// @Tags Auth
// @Accept json
// @Produce json
// @Param provider path string true "Provider name"
// @Success 200 {object} object
// @Failure 500 {object} error
// @Router /auth/{provider}/callback [get]
func (app *application) getAuthCallback(c echo.Context) error {
	provider := c.Param("provider")
	user, err := auth.CompleteAuth(c.Response(), c.Request(), provider)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	}

	_, err = app.store.Users.CreateUser(user.FirstName, user.Email)
	if err != nil {
		switch {
		case errors.Is(err, store.ViolatePK):
			return c.JSON(http.StatusTemporaryRedirect, map[string]string{
				"message": "user already exists, redirecting to login"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
	return c.JSON(http.StatusOK, user)
}

// @Summary Logout handler
func (app *application) logout(c echo.Context) error {
	provider := c.Param("provider")

	if err := auth.Logout(c.Response(), c.Request()); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message":  "logged out",
		"provider": provider,
	})
}
