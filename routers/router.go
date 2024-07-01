package router

import (
    "blog/handlers"
    "github.com/labstack/echo/v4"
)

func New() *echo.Echo {
	e := echo.New()

	e.GET("/posts", handlers.GetAllPosts)
	e.GET("/post/:id", handlers.GetPost)
	e.POST("/posts/:id", handlers.CreatePost)
	e.PUT("/posts/:id", handlers.UpdatePost)
	e.DELETE("/posts/:id", handlers.DeletePost)

	return e
}