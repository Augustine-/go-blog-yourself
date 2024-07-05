package authentication

import (
	"os"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func BasicAuth(username, password string, c echo.Context) (bool, error) {
	expectedUsername := os.Getenv("BLOG_U")
	expectedPassword := os.Getenv("BLOG_P")
	if username == expectedUsername && password == expectedPassword {
		c.Set("isLoggedIn", true)
		return true, nil
	}
	return false, nil
}

func AuthMiddleware() echo.MiddlewareFunc {
	return middleware.BasicAuth(BasicAuth)
}
