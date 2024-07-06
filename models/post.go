package models

type Post struct {
	ID 		int `json:id`
	Title 	string `json:"title"`
	Content string `json:"content"`
	ImageURL string `json:"image_url"`
}