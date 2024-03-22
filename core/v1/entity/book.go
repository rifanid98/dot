package entity

import (
	"time"
)

type Book struct {
	Id       string     `json:"id"`
	Name     string     `json:"name"`
	Author   string     `json:"author"`
	Created  *time.Time `json:"created"`
	Modified *time.Time `json:"modified"`
}
