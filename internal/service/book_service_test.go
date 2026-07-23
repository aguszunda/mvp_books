package service

import (
	"context"
	"errors"
	"testing"

	"mvp_books/internal/domain"
	"mvp_books/internal/service/mocks"

	"github.com/stretchr/testify/assert"
)

func TestBookService_Create(t *testing.T) {
	tests := []struct {
		name      string
		book      *domain.Book
		mockSetup func(*mocks.MockBookRepository)
		expectErr bool
	}{
		{
			name: "Success",
			book: &domain.Book{Title: "T1", Author: "A1"},
			mockSetup: func(m *mocks.MockBookRepository) {
				m.CreateFunc = func(_ context.Context, _ *domain.Book) error {
					return nil
				}
			},
			expectErr: false,
		},
		{
			name: "Repo Error",
			book: &domain.Book{Title: "T1"},
			mockSetup: func(m *mocks.MockBookRepository) {
				m.CreateFunc = func(_ context.Context, _ *domain.Book) error {
					return errors.New("db error")
				}
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.MockBookRepository{}
			if tt.mockSetup != nil {
				tt.mockSetup(repo)
			}
			svc := NewBookService(repo)

			err := svc.Create(context.Background(), tt.book)

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
		mockSetup func(*mocks.MockBookRepository)
		expectErr bool
		expectLen int
	}{
		{
			name: "Success",
			mockSetup: func(m *mocks.MockBookRepository) {
				m.FindAllFunc = func(_ context.Context) ([]domain.Book, error) {
					return []domain.Book{{Title: "B1"}, {Title: "B2"}}, nil
				}
			},
			expectErr: false,
			expectLen: 2,
		},
		{
			name: "Repo Error",
			mockSetup: func(m *mocks.MockBookRepository) {
				m.FindAllFunc = func(_ context.Context) ([]domain.Book, error) {
					return nil, errors.New("db error")
				}
			},
			expectErr: true,
			expectLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.MockBookRepository{}
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
		mockSetup func(*mocks.MockBookRepository)
		expectErr bool
		expectVal *domain.Book
	}{
		{
			name:    "Success",
			inputID: "1",
			mockSetup: func(m *mocks.MockBookRepository) {
				m.FindByIDFunc = func(_ context.Context, _ string) (*domain.Book, error) {
					return &domain.Book{ID: 1, Title: "Found"}, nil
				}
			},
			expectErr: false,
			expectVal: &domain.Book{ID: 1, Title: "Found"},
		},
		{
			name:    "Not Found",
			inputID: "999",
			mockSetup: func(m *mocks.MockBookRepository) {
				m.FindByIDFunc = func(_ context.Context, _ string) (*domain.Book, error) {
					return nil, errors.New("not found")
				}
			},
			expectErr: true,
			expectVal: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mocks.MockBookRepository{}
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
