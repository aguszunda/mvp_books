package port

import (
	"context"

	"datadog-exercise/internal/domain"
)

type BookService interface {
	Create(ctx context.Context, book *domain.Book) error
	GetAll(ctx context.Context) ([]domain.Book, error)
	GetOne(ctx context.Context, id string) (*domain.Book, error)
}
