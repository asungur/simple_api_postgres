package models

type Todo struct {
	Id    int64  `json:"id"`
	Title string `json:"title"`
	Done  string `json:"done"`
}
