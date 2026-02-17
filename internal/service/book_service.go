package service

import (
	"context"

	"datadog-exercise/internal/domain"
)

type bookService struct {
	repo domain.BookRepository
}

func NewBookService(repo domain.BookRepository) BookService {
	return &bookService{repo: repo}
}

func (s *bookService) Create(ctx context.Context, book *domain.Book) error {
	return s.repo.Create(ctx, book)
}

func (s *bookService) GetAll(ctx context.Context) ([]domain.Book, error) {
	return s.repo.FindAll(ctx)
}

func (s *bookService) GetOne(ctx context.Context, id string) (*domain.Book, error) {
	return s.repo.FindByID(ctx, id)
}
