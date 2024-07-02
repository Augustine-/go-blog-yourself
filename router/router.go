package router

import (
    "blog/handlers"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "log"
)

func New() *echo.Echo {
    e := echo.New()

    // Log to confirm middleware setup
    log.Println("Setting up method override middleware")

    // Add method override middleware
    e.Use(middleware.MethodOverrideWithConfig(middleware.MethodOverrideConfig{
        Getter: func(c echo.Context) string {
            method := c.FormValue("_method")
            c.Logger().Infof("Method override: %s", method)
            return method
        },
    }))

    // Define routes
    e.GET("/posts", handlers.GetAllPosts)
    e.GET("/posts/:id", handlers.GetPost)
    e.GET("/posts/new", handlers.NewPostForm)
    e.POST("/posts", handlers.CreatePost)
    e.PUT("/posts/:id", handlers.UpdatePost)
    e.DELETE("/posts/:id", handlers.DeletePost)

    return e
}
