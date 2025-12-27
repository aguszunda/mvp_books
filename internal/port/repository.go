package port

import (
	"context"

	"datadog-exercise/internal/domain"
)

type BookRepository interface {
	Create(ctx context.Context, book *domain.Book) error
	FindAll(ctx context.Context) ([]domain.Book, error)
	FindByID(ctx context.Context, id string) (*domain.Book, error)
}
