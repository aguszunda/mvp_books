package service

import (
	"context"

	"datadog-exercise/internal/domain"
)

// mockBookRepository manually mocks repository.BookRepository for unit tests.
type mockBookRepository struct {
	createFunc   func(ctx context.Context, book *domain.Book) error
	findAllFunc  func(ctx context.Context) ([]domain.Book, error)
	findByIDFunc func(ctx context.Context, id string) (*domain.Book, error)
}

func (m *mockBookRepository) Create(ctx context.Context, book *domain.Book) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, book)
	}
	return nil
}

func (m *mockBookRepository) FindAll(ctx context.Context) ([]domain.Book, error) {
	if m.findAllFunc != nil {
		return m.findAllFunc(ctx)
	}
	return nil, nil
}

func (m *mockBookRepository) FindByID(ctx context.Context, id string) (*domain.Book, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, nil
}
