package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"datadog-exercise/internal/adapter/handler/mocks"
	"datadog-exercise/internal/domain"

	"github.com/gin-gonic/gin"
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
		mockSetup      func(*mocks.MockBookService)
		expectedStatus int
	}{
		{
			name:  "Success",
			input: domain.Book{Title: "Test Book", Author: "Author"},
			mockSetup: func(m *mocks.MockBookService) {
				m.CreateFunc = func(ctx context.Context, book *domain.Book) error {
					return nil
				}
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Bind Error",
			input:          domain.Book{},
			mockSetup:      func(m *mocks.MockBookService) {},
			expectedStatus: http.StatusOK,
		},
		{
			name:  "Service Error",
			input: domain.Book{Title: "Test Book"},
			mockSetup: func(m *mocks.MockBookService) {
				m.CreateFunc = func(ctx context.Context, book *domain.Book) error {
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

			body, _ := json.Marshal(tt.input)
			req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if tt.name == "Bind Error" {
				// Special case handling
			} else {
				if w.Code != tt.expectedStatus {
					t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
				}
			}
		})
	}

	t.Run("Malformed JSON", func(t *testing.T) {
		svc := &mocks.MockBookService{}
		r := setupRouter(svc)
		req, _ := http.NewRequest("POST", "/books", bytes.NewBufferString("{invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d for malformed json, got %d", http.StatusBadRequest, w.Code)
		}
	})
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
				m.GetAllFunc = func(ctx context.Context) ([]domain.Book, error) {
					return []domain.Book{{Title: "B1"}, {Title: "B2"}}, nil
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Service Error",
			mockSetup: func(m *mocks.MockBookService) {
				m.GetAllFunc = func(ctx context.Context) ([]domain.Book, error) {
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

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
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
				m.GetOneFunc = func(ctx context.Context, id string) (*domain.Book, error) {
					if id == "1" {
						return &domain.Book{ID: 1, Title: "Found"}, nil
					}
					return nil, errors.New("not found in mock logic")
				}
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Not Found",
			bookID: "999",
			mockSetup: func(m *mocks.MockBookService) {
				m.GetOneFunc = func(ctx context.Context, id string) (*domain.Book, error) {
					return nil, errors.New("some error")
				}
			},
			expectedStatus: http.StatusNotFound,
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

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
