package products

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	repo "github.com/carlosEA28/ecom/internal/adapters/postgresql/sqlc"
)

type mockService struct {
	listProductsFunc func(ctx context.Context) ([]repo.Product, error)
}

func (m *mockService) ListProducts(ctx context.Context) ([]repo.Product, error) {
	return m.listProductsFunc(ctx)
}

func TestListProducts(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockProducts := []repo.Product{
			{ID: 1, Name: "Product 1"},
			{ID: 2, Name: "Product 2"},
		}

		service := &mockService{
			listProductsFunc: func(ctx context.Context) ([]repo.Product, error) {
				return mockProducts, nil
			},
		}

		handler := NewHandler(service)
		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		rr := httptest.NewRecorder()

		handler.ListProducts(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rr.Code)
		}

		var got []repo.Product
		if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
			t.Fatal(err)
		}

		if len(got) != len(mockProducts) {
			t.Errorf("expected %d products, got %d", len(mockProducts), len(got))
		}
	})

	t.Run("service error", func(t *testing.T) {
		service := &mockService{
			listProductsFunc: func(ctx context.Context) ([]repo.Product, error) {
				return nil, errors.New("something went wrong")
			},
		}

		handler := NewHandler(service)
		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		rr := httptest.NewRecorder()

		handler.ListProducts(rr, req)

		if rr.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", rr.Code)
		}
	})
}
