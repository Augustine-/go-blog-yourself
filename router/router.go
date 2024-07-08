package router

import (
    "blog/handlers"
    "blog/authentication"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "github.com/gorilla/sessions"
    echoSession "github.com/labstack/echo-contrib/session"
)

func New() *echo.Echo {
    e := echo.New()

    e.Static("/static", "static")

    e.Pre(middleware.RemoveTrailingSlash())

    e.Use(middleware.MethodOverrideWithConfig(middleware.MethodOverrideConfig{
        Getter: middleware.MethodFromForm("_method"),
    }))

    // Session middleware
    e.Use(echoSession.Middleware(sessions.NewCookieStore([]byte("secret"))))

    // Public routes
    e.GET("/", handlers.GetAllPosts)
    e.GET("/posts", handlers.GetAllPosts)
    e.GET("/posts/:id", handlers.GetPost)

    // Login route
    e.GET("/login", handlers.ShowLoginForm)
    e.POST("/login", handlers.Login, authentication.BasicAuthMiddleware())

    // Protected routes
    e.GET("/posts/new", handlers.NewPostForm, authentication.AuthMiddleware)
    e.GET("/posts/edit/:id", handlers.EditPostForm, authentication.AuthMiddleware)
    e.POST("/posts", handlers.CreatePost, authentication.AuthMiddleware)
    e.POST("/posts/:id", handlers.UpdatePost, authentication.AuthMiddleware)
    e.POST("/posts/:id/delete", handlers.DeletePost, authentication.AuthMiddleware)

    return e
}
