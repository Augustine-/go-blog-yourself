package authentication

import (
	"os"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo-contrib/session"
	"log"
)

func BasicAuth(username, password string, c echo.Context) (bool, error) {
	expectedUsername := os.Getenv("BLOG_U")
	expectedPassword := os.Getenv("BLOG_P")
	if username == expectedUsername && password == expectedPassword {
		sess, _ := session.Get("session", c)
		sess.Values["isAuthenticated"] = true
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			log.Println("Failed to save session:", err)
			return false, err
		}
		log.Println("User authenticated successfully")
		return true, nil
	}
	log.Println("Authentication failed")
	return false, nil
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get("session", c)
		if auth, ok := sess.Values["isAuthenticated"].(bool); !ok || !auth {
			return c.String(401, "Unauthorized")
		}
		return next(c)
	}
}

func BasicAuthMiddleware() echo.MiddlewareFunc {
	return middleware.BasicAuth(BasicAuth)
}
