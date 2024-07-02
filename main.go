package main

import (
	"blog/router"
	"blog/database"
	"blog/renderer"
)

func main() {
	database.InitDB()

	e := router.New()

	e.Renderer = renderer.NewRenderer()

	e.Logger.Fatal(e.Start(":1323"))
}