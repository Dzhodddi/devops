package main

import (
	"context"
	"devops/internal/auth"
	"devops/internal/store"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type userKey string

const userCtx userKey = "user"

func (app *application) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user, err := auth.GetUserFromSession(c.Request())
		if err != nil || user.UserID == "" {
			return c.Redirect(http.StatusTemporaryRedirect, "auth/google")
		}
		ctx := context.WithValue(c.Request().Context(), userCtx, user.Email)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}

}

func (app *application) PostContextMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		postID := c.Param("id")
		id, err := strconv.ParseInt(postID, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid post ID"})
		}
		post, err := app.store.Posts.GetByID(c.Request().Context(), id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				return c.JSON(http.StatusNotFound, map[string]string{"error": "Post not found"})
			default:
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
			}
		}

		ctx := context.WithValue(c.Request().Context(), postCtx, post)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
