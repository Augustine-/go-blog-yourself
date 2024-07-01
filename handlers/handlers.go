package handlers

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

func GetAllPosts(c echo.Context) error {
	return c.String(http.StatusOK, "All Posts")
}

func GetPost(c echo.Context) error {
	id := c.Param("id")

	return c.String(http.StatusOK, "Post "+id)
}

func CreatePost(c echo.Context) error {
	return c.String(http.StatusCreated, "Created New Post")
}

func UpdatePost(c echo.Context) error {
	id := c.Param("id")

	return c.String(http.StatusOK, "Updated Post "+id)
}

func DeletePost(c echo.Context) error {
	id := c.Param("id")

	return c.String(http.StatusOK, "Deleted Post "+id)
}