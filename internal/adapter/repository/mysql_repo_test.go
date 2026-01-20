package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"datadog-exercise/internal/domain"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock
}

func TestMysqlBookRepository_Create(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewMysqlBookRepository(db)

	t.Run("Success", func(t *testing.T) {
		book := &domain.Book{Title: "Test Book", Author: "Test Author"}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `books` (`title`,`author`) VALUES (?,?)")).
			WithArgs(book.Title, book.Author).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.Create(context.Background(), book)
		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DB Error", func(t *testing.T) {
		book := &domain.Book{Title: "Test Book", Author: "Test Author"}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `books`")).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		err := repo.Create(context.Background(), book)
		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMysqlBookRepository_FindAll(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewMysqlBookRepository(db)

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "author"}).
			AddRow(1, "Book 1", "Author 1").
			AddRow(2, "Book 2", "Author 2")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `books`")).
			WillReturnRows(rows)

		books, err := repo.FindAll(context.Background())
		assert.NoError(t, err)
		assert.Len(t, books, 2)
		assert.Equal(t, "Book 1", books[0].Title)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DB Error", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `books`")).
			WillReturnError(errors.New("db error"))

		books, err := repo.FindAll(context.Background())
		assert.Error(t, err)
		assert.Nil(t, books)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestMysqlBookRepository_FindByID(t *testing.T) {
	db, mock := newMockDB(t)
	repo := NewMysqlBookRepository(db)

	t.Run("Success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "author"}).
			AddRow(1, "Book 1", "Author 1")

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `books` WHERE `books`.`id` = ? ORDER BY `books`.`id` LIMIT ?")).
			WithArgs("1", 1).
			WillReturnRows(rows)

		book, err := repo.FindByID(context.Background(), "1")
		assert.NoError(t, err)
		assert.Equal(t, "Book 1", book.Title)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `books` WHERE `books`.`id` = ?")).
			WithArgs("999", 1).
			WillReturnError(gorm.ErrRecordNotFound)

		book, err := repo.FindByID(context.Background(), "999")
		assert.Error(t, err) // GORM returns error on not found if checking .Error (which repo does)
		assert.Nil(t, book)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
