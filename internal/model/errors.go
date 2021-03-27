package model

type Error struct {
	Name    string   `json:"name"`
	Message string   `json:"message"`
	Tags    []string `json:"tags"`
}
