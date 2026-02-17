package service

import (
	"context"
	"errors"
	"testing"

	"datadog-exercise/internal/domain"

	"github.com/stretchr/testify/assert"
)

// mockBookRepository manually mocks port.BookRepository
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

func TestBookService_Create(t *testing.T) {
	tests := []struct {
		name      string
		input     *domain.Book
		mockSetup func(*mockBookRepository)
		expectErr bool
	}{
		{
			name:  "Success",
			input: &domain.Book{Title: "T1", Author: "A1"},
			mockSetup: func(m *mockBookRepository) {
				m.createFunc = func(ctx context.Context, book *domain.Book) error {
					return nil
				}
			},
			expectErr: false,
		},
		{
			name:  "Repo Error",
			input: &domain.Book{Title: "T1"},
			mockSetup: func(m *mockBookRepository) {
				m.createFunc = func(ctx context.Context, book *domain.Book) error {
					return errors.New("db error")
				}
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockBookRepository{}
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}
			svc := NewBookService(repo)

			err := svc.Create(context.Background(), tt.input)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBookService_GetAll(t *testing.T) {
	tests := []struct {
		name      string
		mockSetup func(*mockBookRepository)
		expectErr bool
		expectLen int
	}{
		{
			name: "Success",
			mockSetup: func(m *mockBookRepository) {
				m.findAllFunc = func(ctx context.Context) ([]domain.Book, error) {
					return []domain.Book{{Title: "B1"}, {Title: "B2"}}, nil
				}
			},
			expectErr: false,
			expectLen: 2,
		},
		{
			name: "Repo Error",
			mockSetup: func(m *mockBookRepository) {
				m.findAllFunc = func(ctx context.Context) ([]domain.Book, error) {
					return nil, errors.New("db error")
				}
			},
			expectErr: true,
			expectLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockBookRepository{}
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}
			svc := NewBookService(repo)

			books, err := svc.GetAll(context.Background())

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, books, tt.expectLen)
			}
		})
	}
}

func TestBookService_GetOne(t *testing.T) {
	tests := []struct {
		name      string
		inputID   string
		mockSetup func(*mockBookRepository)
		expectErr bool
		expectVal *domain.Book
	}{
		{
			name:    "Success",
			inputID: "1",
			mockSetup: func(m *mockBookRepository) {
				m.findByIDFunc = func(ctx context.Context, id string) (*domain.Book, error) {
					return &domain.Book{ID: 1, Title: "Found"}, nil
				}
			},
			expectErr: false,
			expectVal: &domain.Book{ID: 1, Title: "Found"},
		},
		{
			name:    "Not Found",
			inputID: "999",
			mockSetup: func(m *mockBookRepository) {
				m.findByIDFunc = func(ctx context.Context, id string) (*domain.Book, error) {
					return nil, errors.New("not found")
				}
			},
			expectErr: true,
			expectVal: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockBookRepository{}
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}
			svc := NewBookService(repo)

			book, err := svc.GetOne(context.Background(), tt.inputID)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, book)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectVal, book)
			}
		})
	}
}
