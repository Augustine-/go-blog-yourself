package handlers

import (
	"net/http"
    "html/template"
	"github.com/labstack/echo/v4"
    "github.com/labstack/echo-contrib/session"
	"blog/models"
	"blog/database"
	"log"
	"os"
	"io"
	"path/filepath"
	"strconv"
)

func GetAllPosts(c echo.Context) error {
    sess, _ := session.Get("session", c)
    isAuthenticated, _ := sess.Values["isAuthenticated"].(bool)
    var posts []models.Post

    rows, err := database.DB.Query("SELECT id, title, content, image_url FROM posts")
    if err != nil {
        log.Println("Error fetching posts: ", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error fetching posts"})
    }
    defer rows.Close()

    for rows.Next() {
        var post models.Post
        if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.ImageURL); err != nil {
            log.Println("Error scanning post: ", err)
            return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error scanning post"})
        }
        posts = append(posts, post)
    }

    if err = rows.Err(); err != nil {
        log.Println("Rows error: ", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Rows error"})
    }

    if len(posts) == 0 {
        log.Println("No posts found")
        return c.JSON(http.StatusNotFound, echo.Map{"error": "No posts found"})
    }

    err = c.Render(http.StatusOK, "all_posts.html", map[string]interface{}{
        "Posts": posts,
        "isAuthenticated": isAuthenticated,
    })

    if err != nil {
        log.Println("Error rendering template: ", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error rendering template"})
    }

    return nil
}

func GetPost(c echo.Context) error {
    sess, _ := session.Get("session", c)
    isAuthenticated, _ := sess.Values["isAuthenticated"].(bool)
    id := c.Param("id")
    var post models.Post

    row := database.DB.QueryRow("SELECT id, title, content, image_url FROM posts WHERE id = ?", id)
    err := row.Scan(&post.ID, &post.Title, &post.Content, &post.ImageURL)

    if err != nil {
        log.Println("Error fetching post: ", err)
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Post not found."})
    }

    // Mark content as safe HTML
    safeContent := template.HTML(post.Content)

    return c.Render(http.StatusOK, "view_post.html", map[string]interface{}{
        "Post":           post,
        "isAuthenticated": isAuthenticated,
        "SafeContent":    safeContent,
    })
}


func CreatePost(c echo.Context) error {
	var newPost models.Post
	newPost.Title = c.FormValue("title")
	newPost.Content = c.FormValue("content")

	file, err := c.FormFile("image")
	if err != nil {
		log.Println("Error retrieving file:", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error":"Failed to retrieve image"})
	}
	src, err := file.Open()
	if err != nil {
		log.Println("Error opening file:", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error":"Failed to open image"})
	}
	defer src.Close()

	dstPath := filepath.Join("static/images", file.Filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		log.Println("Error creating destination file:", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error":"Failed to create destination file"})
	}
	defer dst.Close()

	// copy file from src to dst
	if _, err = io.Copy(dst, src); err != nil {
		log.Println("Error copying file:", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error":"Failed to copy file"})
	}

	newPost.ImageURL = "/static/images/" + file.Filename

	if err := database.SavePostToDB(newPost); err != nil {
		log.Println("Error copying file:", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error":"Failed to copy file"})
	}
	return c.Redirect(http.StatusSeeOther, "/posts")
}

func UpdatePost(c echo.Context) error {
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        log.Println("Invalid post ID:", err)
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid post ID"})
    }

    var updatedPost models.Post
    updatedPost.ID = id
    updatedPost.Title = c.FormValue("title")
    updatedPost.Content = c.FormValue("content")

    file, err := c.FormFile("image")
    if err == nil {
        src, err := file.Open()
        if err != nil {
            log.Println("Error opening file:", err)
            return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to open image"})
        }
        defer src.Close()

        dstPath := filepath.Join("static/images", file.Filename)
        dst, err := os.Create(dstPath)
        if err != nil {
            log.Println("Error creating destination file:", err)
            return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to create destination file"})
        }
        defer dst.Close()

        if _, err = io.Copy(dst, src); err != nil {
            log.Println("Error copying file:", err)
            return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to copy file"})
        }

        updatedPost.ImageURL = "/static/images/" + file.Filename
    } else if err == http.ErrMissingFile {
        // Retrieve the current image URL from the database
        row := database.DB.QueryRow("SELECT image_url FROM posts WHERE id = ?", id)
        err = row.Scan(&updatedPost.ImageURL)
        if err != nil {
            log.Println("Error fetching current image URL:", err)
            return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to retrieve current image URL"})
        }
    } else {
        log.Println("Error retrieving file:", err)
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to retrieve image"})
    }

    if err := database.UpdatePostInDB(updatedPost); err != nil {
        log.Println("Error updating post:", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update post"})
    }

    return c.Redirect(http.StatusSeeOther, "/posts/" + idStr)
}

func DeletePost(c echo.Context) error {
	id := c.Param("id")
	statement, err := database.DB.Prepare("DELETE FROM posts WHERE id = ?")

	if err != nil {
		log.Println("Error preparing statement:", err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to execute delete statement"})
	}

	_, err = statement.Exec(id)
	if err != nil {
		log.Println("Error executing statement:", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to execute DB statement."})
	}

	log.Println("Received delete data: " + id) // confirm deletion
	return c.Redirect(http.StatusSeeOther, "/posts")
}

func NewPostForm(c echo.Context) error {
	return c.Render(http.StatusOK, "new_post.html", nil)
}

func EditPostForm(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
        log.Println("Invalid post ID:", err)
        return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid post ID"})
    }

	var post models.Post

	row := database.DB.QueryRow("SELECT id, title, content, image_url FROM posts WHERE id = ?", id)
	err = row.Scan(&post.ID, &post.Title, &post.Content, &post.ImageURL)
	if err != nil {
        log.Println("Error fetching post: ", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error fetching post"})
    }


	return c.Render(http.StatusOK, "edit_post.html", post)
}

func Login(c echo.Context) error {
	sess, _ := session.Get("session", c)
	if auth, ok := sess.Values["isAuthenticated"].(bool); ok && auth {
		return c.Redirect(http.StatusFound, "/posts")
	}
	return c.String(http.StatusUnauthorized, "Invalid username or password")
}

func ShowLoginForm(c echo.Context) error {
    return c.Render(http.StatusOK, "login.html", nil)
}