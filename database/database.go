package database

import(
	"database/sql"
	"log"
	"blog/models"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./database/blog.db")

	if err != nil {
		log.Fatal(err)
	}

	createTable()
}

func createTable() {
	createPostTableSQL := `CREATE TABLE IF NOT EXISTS posts (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"title" TEXT NOT NULL,
		"content" TEXT NOT NULL,
		"image_url" TEXT
	);`

	statement, err := DB.Prepare(createPostTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	statement.Exec()
}

func SavePostToDB(post models.Post) error {
	statement, err := DB.Prepare("INSERT INTO posts (title, content, image_url) VALUES (?, ? ,?)")
	if err != nil {
		return err
	}
	defer statement.Close()

	result, err := statement.Exec(post.Title, post.Content, post.ImageURL)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	post.ID = int(id)
	return nil
}

func UpdatePostInDB(post models.Post) error {
    query := "UPDATE posts SET title = ?, content = ?, image_url = ? WHERE id = ?"

    stmt, err := DB.Prepare(query)
    if err != nil {
        log.Println("Error preparing query:", err)
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(post.Title, post.Content, post.ImageURL, post.ID)
    if err != nil {
        log.Println("Error executing query:", err)
        return err
    }

    return nil
}