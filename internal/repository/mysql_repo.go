package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"datadog-exercise/internal/domain"
)

type BookRepository interface {
	Create(ctx context.Context, book *domain.Book) error
	FindAll(ctx context.Context) ([]domain.Book, error)
	FindByID(ctx context.Context, id string) (*domain.Book, error)
}

type MysqlBookRepository struct {
	db *gorm.DB
}

func NewMysqlBookRepository(db *gorm.DB) BookRepository {
	return &MysqlBookRepository{db: db}
}

func (r *MysqlBookRepository) Create(ctx context.Context, book *domain.Book) error {
	return r.db.WithContext(ctx).Create(book).Error
}

func (r *MysqlBookRepository) FindAll(ctx context.Context) ([]domain.Book, error) {
	var books []domain.Book
	err := r.db.WithContext(ctx).Find(&books).Error
	return books, err
}

func (r *MysqlBookRepository) FindByID(ctx context.Context, id string) (*domain.Book, error) {
	var book domain.Book
	if err := r.db.WithContext(ctx).First(&book, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &book, nil
}
