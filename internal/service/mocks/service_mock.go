package mocks

import (
	"context"

	"mvp_books/internal/domain"
)

// MockBookRepository manually mocks repository.BookRepository for unit tests.
type MockBookRepository struct {
	CreateFunc   func(ctx context.Context, book *domain.Book) error
	FindAllFunc  func(ctx context.Context) ([]domain.Book, error)
	FindByIDFunc func(ctx context.Context, id string) (*domain.Book, error)
}

func (m *MockBookRepository) Create(ctx context.Context, book *domain.Book) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, book)
	}
	return nil
}

func (m *MockBookRepository) FindAll(ctx context.Context) ([]domain.Book, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx)
	}
	return nil, nil
}

func (m *MockBookRepository) FindByID(ctx context.Context, id string) (*domain.Book, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, nil
}
