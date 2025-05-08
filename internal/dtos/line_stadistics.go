package dtos

import "time"

type LinePickingStats struct {
	Item        string    `json:"itemName"`
	Ean         string    `json:"ean"`
	Worker      string    `json:"worker"`
	CurrentTime time.Time `json:"currentTime"`
	Quantity    int       `json:"quantity"`
}
