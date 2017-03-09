package models

import (
	"time"
)

// Favorite is the model for a user
type Favorite struct {
	Word    string    `json:"word"`
	Created time.Time `json:"created"`
}
