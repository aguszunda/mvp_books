package repository

import (
	"context"

	"gorm.io/gorm"

	"datadog-exercise/internal/domain"
	"datadog-exercise/internal/port"
)

type MysqlBookRepository struct {
	db *gorm.DB
}

func NewMysqlBookRepository(db *gorm.DB) port.BookRepository {
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
		return nil, err
	}
	return &book, nil
}
