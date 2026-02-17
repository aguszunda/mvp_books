package domain

import (
	"context"
)

type BookRepository interface {
	Create(ctx context.Context, book *Book) error
	FindAll(ctx context.Context) ([]Book, error)
	FindByID(ctx context.Context, id string) (*Book, error)
}
