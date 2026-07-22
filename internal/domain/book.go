package domain

import "errors"

var ErrNotFound = errors.New("book not found")

type Book struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}
