package main

import (
	"blog/router"
	"blog/database"
	"blog/renderer"
	"os"
)



func main() {
    database.InitDB()

    e := router.New()

    e.Renderer = renderer.NewRenderer()

    port := os.Getenv("PORT")
    if port == "" {
        port = "1323" // Default to 1323 if PORT is not set
    }

    e.Logger.Fatal(e.Start(":" + port))
}