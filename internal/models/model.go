package models

type Link struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Link  string `json:"link"`
}

type CreateLinkRequest struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

type UpdateLinkRequest struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

type APIResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}
