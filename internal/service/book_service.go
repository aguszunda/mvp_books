package service

import (
	"context"

	"mvp_books/internal/domain"
	"mvp_books/internal/repository"
)

type BookService interface {
	Create(ctx context.Context, book *domain.Book) error
	GetAll(ctx context.Context) ([]domain.Book, error)
	GetOne(ctx context.Context, id string) (*domain.Book, error)
}

type bookService struct {
	repo repository.BookRepository
}

func NewBookService(repo repository.BookRepository) BookService {
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
