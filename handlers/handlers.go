package handlers

import (
	"net/http"
	"github.com/labstack/echo/v4"
	"blog/models"
	"blog/database"
	"log"
)

func GetAllPosts(c echo.Context) error {
    var posts []models.Post

    rows, err := database.DB.Query("SELECT id, title, content FROM posts")
    if err != nil {
        log.Println("Error fetching posts: ", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error fetching posts"})
    }
    defer rows.Close()

    for rows.Next() {
        var post models.Post
        if err := rows.Scan(&post.ID, &post.Title, &post.Content); err != nil {
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

    log.Printf("Rendering all posts: %+v", posts)

    err = c.Render(http.StatusOK, "all_posts.html", map[string]interface{}{
        "Posts": posts,
    })
    if err != nil {
        log.Println("Error rendering template: ", err)
        return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Error rendering template"})
    }

    return nil
}


func GetPost(c echo.Context) error {
	id := c.Param("id")
	var post models.Post

	row := database.DB.QueryRow("SELECT id, title, content FROM posts WHERE id = ?", id)
	err := row.Scan(&post.ID, &post.Title, &post.Content)

	if err != nil {
		log.Println("Error fetching post: ", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Post not found."})
	}

	log.Printf("Rendering post: %+v", post) // Log the post data to verify

	return c.Render(http.StatusOK, "view_post.html", post)
}

func CreatePost(c echo.Context) error {
	var newPost models.Post
	newPost.Title = c.FormValue("title")
	newPost.Content = c.FormValue("content")
	log.Println("Form Values:", c.FormValue("title"), c.FormValue("content"))

	if err := c.Bind(&newPost); err != nil {
		log.Println("Error binding post data:", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}


	statement, err := database.DB.Prepare("INSERT INTO posts (title, content) VALUES (?, ?)")
	if err != nil {
		log.Println("Error preparing statement:", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to prepare DB statement."})
	}

	result, err := statement.Exec(newPost.Title, newPost.Content)
	if err != nil {
		log.Println("Error executing statement:", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to execute DB statement."})
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting insert ID:", err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to retrieve ID"})
	}

	newPost.ID = int(id)

	log.Printf("Received post data: %+v", newPost) // confirm we received data
	return c.Redirect(http.StatusSeeOther, "/posts")
}

func UpdatePost(c echo.Context) error {
	id := c.Param("id")

	return c.String(http.StatusOK, "Updated Post "+id)
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
