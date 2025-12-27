package service

import (
	"context"

	"datadog-exercise/internal/domain"
	"datadog-exercise/internal/port"
)

type BookService struct {
	repo port.BookRepository
}

func NewBookService(repo port.BookRepository) port.BookService {
	return &BookService{repo: repo}
}

func (s *BookService) Create(ctx context.Context, book *domain.Book) error {
	return s.repo.Create(ctx, book)
}

func (s *BookService) GetAll(ctx context.Context) ([]domain.Book, error) {
	return s.repo.FindAll(ctx)
}

func (s *BookService) GetOne(ctx context.Context, id string) (*domain.Book, error) {
	return s.repo.FindByID(ctx, id)
}
