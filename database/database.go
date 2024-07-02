package database

import(
	"database/sql"
	"log"
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
		"content" text nnot null
	);`

	statement, err := DB.Prepare(createPostTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	statement.Exec()
}