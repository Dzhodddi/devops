package main

import (
	"devops/internal/store"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"strconv"
)

type postKey string

const path = "/app/public"
const postCtx postKey = "post"

type CreatePostPayload struct {
	Title   string `json:"title" validate:"required,min=1,max=100"`
	Content string `json:"content" validate:"required,min=1,max=1000"`
}

type EditPostPayload struct {
	Content   string `json:"content" validate:"required,min=1,max=1000"`
	ImagePath string `json:"image,omitempty" validate:"min=1,max=1000"`
}

// @Summary Create a new createPost
// @Description Create a new post
// @Tags posts
// @Accept json
// @Produce json
// @Param payload body CreatePostPayload true "Post data"
// @Success 201 {object} store.Post
// @Failure 400 {object} map[string]string "Invalid request format"
// @Failure 422 {object} map[string]string "Validation error"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /post [post]
func (app *application) createPost(c echo.Context) error {
	var req CreatePostPayload
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}
	if err := Validate.Struct(req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
	}

	post := &store.Post{
		Title:       req.Title,
		Content:     req.Content,
		AuthorEmail: os.Getenv("email"),
	}

	if err := app.store.Posts.Create(c.Request().Context(), post); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create post"})
	}
	return c.JSON(http.StatusCreated, post)

}

// @Summary Get a list of getPosts
// @Description Retrieve a list of posts
// @Tags posts
// @Accept json
// @Produce json
// @Success 200 {array} store.Post
// @Failure 500 {object} map[string]string "Internal server error"
// @Router  /post [get]
func (app *application) getPosts(c echo.Context) error {
	posts, err := app.store.Posts.GetList(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve posts"})
	}
	return c.JSON(http.StatusOK, posts)
}

// @Summary Edit an existing post
// @Description Edit an existing post
// @Tags posts
// @Accept multipart/form-data
// @Produce json
// @Param post body EditPostPayload true "Post data"
// @Param id path int true "Post id"
// @Success 200 {object} store.Post
// @Failure 400 {object} map[string]string "Invalid request format"
// @Failure 422 {object} map[string]string "Validation error"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /post/{id} [patch]
func (app *application) editPost(c echo.Context) error {
	post, err := app.getPostFromContext(c)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Post not found"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve post"})
		}
	}
	if err := c.Request().ParseMultipartForm(10 << 20); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid form data"})
	}

	content := c.FormValue("content")
	if content != "" {
		post.Content = content
	}

	file, err := c.FormFile("photo")
	if err == nil {
		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to open uploaded file"})
		}
		defer src.Close()

		// Ensure directory exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := os.MkdirAll(path, 0755); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create upload directory"})
			}
		}

		// Save with unique filename
		dstPath := fmt.Sprintf("%s/%s", path, file.Filename)
		dst, err := os.Create(dstPath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to save file"})
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to write file"})
		}

		// Store relative path for serving
		post.PhotoURL = file.Filename
	}

	// Save changes to DB
	if err = app.store.Posts.Edit(c.Request().Context(), post); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update post"})
	}

	return c.JSON(http.StatusOK, post)

}

// @Summary Delete an existing post
// @Description Delete an existing post
// @Tags posts
// @Accept json
// @Produce json
// @Param id path int true "Post id"
// @Success 204 "No Content"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /post/{id} [delete]
func (app *application) deletePost(c echo.Context) error {
	post, err := app.getPostFromContext(c)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Post not found"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to retrieve post"})
		}
	}

	if err = app.store.Posts.Delete(c.Request().Context(), post.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete post"})
	}
	return c.NoContent(http.StatusNoContent)
}

func (app *application) getPostFromContext(c echo.Context) (*store.Post, error) {
	postID := c.Param("id")
	id, err := strconv.ParseInt(postID, 10, 64)
	if err != nil {
		return nil, err
	}
	post, err := app.store.Posts.GetByID(c.Request().Context(), id)
	if err != nil {
		return nil, err
	}
	return post, nil
}
