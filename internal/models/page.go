package models

type LoadMoreRequest struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
}
