package main

import (
	"devops/docs"
	"devops/internal/store"
	_ "github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type application struct {
	config config
	logger *zap.SugaredLogger
	store  *store.Storage
}

type dbConfig struct {
	addr               string
	maxOpenConnections int
	maxIdleConnections int
	maxIdleTime        string
}

type config struct {
	addr string
	db   dbConfig
	env  string
}

func (app *application) run(mux http.Handler) error {
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.addr
	docs.SwaggerInfo.BasePath = "/v1"
	srv := echo.New()
	srv.Any("/*", echo.WrapHandler(mux))
	err := srv.Start(app.config.addr)
	if err != nil {
		return err
	}
	return nil
}

func (app *application) mount() http.Handler {
	e := echo.New()
	e.Use(middleware.RequestID())
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} [${status}] ${method} ${path} ${latency_human}\n",
	}))
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 60 * time.Second,
	}))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000", "http://0.0.0.0:3000", "http://localhost:5173"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodDelete, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodHead, http.MethodConnect, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))
	e.Static("/public", "/app/public")
	v1 := e.Group("/v1")
	v1.GET("/health", app.healthCheckHandler)
	v1.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/auth/:provider/callback", app.getAuthCallback)
	v1.GET("/auth/:provider", app.startAuthHandler)
	v1.GET("/auth/logout/:provider", app.logout)

	posts := v1.Group("/post")

	posts.POST("", app.createPost)
	posts.GET("", app.getPosts)
	postsID := posts.Group("/:id")
	postsID.PATCH("", app.editPost)
	postsID.DELETE("", app.deletePost)
	return e
}
