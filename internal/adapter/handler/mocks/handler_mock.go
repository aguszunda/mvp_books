package mocks

import (
	"context"

	"datadog-exercise/internal/domain"
)

// MockBookService manually mocks port.BookService
type MockBookService struct {
	CreateFunc func(ctx context.Context, book *domain.Book) error
	GetAllFunc func(ctx context.Context) ([]domain.Book, error)
	GetOneFunc func(ctx context.Context, id string) (*domain.Book, error)
}

func (m *MockBookService) Create(ctx context.Context, book *domain.Book) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, book)
	}
	return nil
}

func (m *MockBookService) GetAll(ctx context.Context) ([]domain.Book, error) {
	if m.GetAllFunc != nil {
		return m.GetAllFunc(ctx)
	}
	return nil, nil
}

func (m *MockBookService) GetOne(ctx context.Context, id string) (*domain.Book, error) {
	if m.GetOneFunc != nil {
		return m.GetOneFunc(ctx, id)
	}
	return nil, nil
}
