package router

import (
    "blog/handlers"
    "blog/authentication"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
    e := echo.New()

    e.Pre(middleware.RemoveTrailingSlash())

    // Public routes
    e.GET("/posts", handlers.GetAllPosts)
    e.GET("/posts/:id", handlers.GetPost)

    // Protected routes
    e.GET("/posts/new", handlers.NewPostForm, authentication.AuthMiddleware())
    e.POST("/posts", handlers.CreatePost, authentication.AuthMiddleware())
    e.PUT("/posts/:id", handlers.UpdatePost, authentication.AuthMiddleware())
    e.DELETE("/posts/:id", handlers.DeletePost, authentication.AuthMiddleware())

    return e
}
