package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"datadog-exercise/internal/controller/mocks"
	"datadog-exercise/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(service *mocks.MockBookService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := NewBookHandler(service)
	h.RegisterRoutes(r)
	return r
}

func TestBookHandler_CreateBook(t *testing.T) {
	tests := []struct {
		name           string
		input          domain.Book
		rawBody        string
		mockSetup      func(*mocks.MockBookService)
		expectedStatus int
	}{
		{
			name:  "Success",
			input: domain.Book{Title: "Test Book", Author: "Author"},
			mockSetup: func(m *mocks.MockBookService) {
				m.CreateFunc = func(_ context.Context, _ *domain.Book) error {
					return nil
				}
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Invalid JSON",
			rawBody:        "{invalid",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:  "Service Error",
			input: domain.Book{Title: "Test Book"},
			mockSetup: func(m *mocks.MockBookService) {
				m.CreateFunc = func(_ context.Context, _ *domain.Book) error {
					return errors.New("database error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mocks.MockBookService{}
			if tt.mockSetup != nil {
				tt.mockSetup(svc)
			}
			r := setupRouter(svc)

			var body []byte
			if tt.rawBody != "" {
				body = []byte(tt.rawBody)
			} else {
				body, _ = json.Marshal(tt.input)
			}
			req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestBookHandler_GetBooks(t *testing.T) {
	tests := []struct {
		name           string
		mockSetup      func(*mocks.MockBookService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			mockSetup: func(m *mocks.MockBookService) {
				m.GetAllFunc = func(_ context.Context) ([]domain.Book, error) {
					return []domain.Book{{Title: "B1"}, {Title: "B2"}}, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Service Error",
			mockSetup: func(m *mocks.MockBookService) {
				m.GetAllFunc = func(_ context.Context) ([]domain.Book, error) {
					return nil, errors.New("db error")
				}
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mocks.MockBookService{}
			if tt.mockSetup != nil {
				tt.mockSetup(svc)
			}
			r := setupRouter(svc)

			req, _ := http.NewRequest("GET", "/books", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestBookHandler_GetBook(t *testing.T) {
	tests := []struct {
		name           string
		bookID         string
		mockSetup      func(*mocks.MockBookService)
		expectedStatus int
	}{
		{
			name:   "Success",
			bookID: "1",
			mockSetup: func(m *mocks.MockBookService) {
				m.GetOneFunc = func(_ context.Context, id string) (*domain.Book, error) {
					if id == "1" {
						return &domain.Book{ID: 1, Title: "Found"}, nil
					}
					return nil, domain.ErrNotFound
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Not Found",
			bookID: "999",
			mockSetup: func(m *mocks.MockBookService) {
				m.GetOneFunc = func(_ context.Context, _ string) (*domain.Book, error) {
					return nil, domain.ErrNotFound
				}
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "Invalid ID",
			bookID: "abc",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "DB Error",
			bookID: "1",
			mockSetup: func(m *mocks.MockBookService) {
				m.GetOneFunc = func(_ context.Context, _ string) (*domain.Book, error) {
					return nil, errors.New("connection refused")
				}
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := &mocks.MockBookService{}
			if tt.mockSetup != nil {
				tt.mockSetup(svc)
			}
			r := setupRouter(svc)

			req, _ := http.NewRequest("GET", "/books/"+tt.bookID, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
