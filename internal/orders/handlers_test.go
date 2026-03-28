package orders

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	repo "github.com/sikozonpc/ecom/internal/adapters/postgresql/sqlc"
)

type mockService struct {
	placeOrderFunc func(ctx context.Context, tempOrder createOrderParams) (repo.Order, error)
}

func (m *mockService) PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error) {
	return m.placeOrderFunc(ctx, tempOrder)
}

func TestPlaceOrder(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		service := &mockService{
			placeOrderFunc: func(ctx context.Context, tempOrder createOrderParams) (repo.Order, error) {
				return repo.Order{ID: 1, CustomerID: 1}, nil
			},
		}

		handler := NewHandler(service)
		body, _ := json.Marshal(createOrderParams{
			CustomerID: 1,
			Items:      []orderItem{{ProductID: 1, Quantity: 1}},
		})
		req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.PlaceOrder(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", rr.Code)
		}

		var got repo.Order
		if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
			t.Fatal(err)
		}

		if got.ID != 1 {
			t.Errorf("expected order ID 1, got %d", got.ID)
		}
	})

	t.Run("bad request - invalid json", func(t *testing.T) {
		handler := NewHandler(&mockService{})
		req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader([]byte("{invalid json}")))
		rr := httptest.NewRecorder()

		handler.PlaceOrder(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rr.Code)
		}
	})

	t.Run("not found - product not found", func(t *testing.T) {
		service := &mockService{
			placeOrderFunc: func(ctx context.Context, tempOrder createOrderParams) (repo.Order, error) {
				return repo.Order{}, ErrProductNotFound
			},
		}

		handler := NewHandler(service)
		body, _ := json.Marshal(createOrderParams{
			CustomerID: 1,
			Items:      []orderItem{{ProductID: 999, Quantity: 1}},
		})
		req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		handler.PlaceOrder(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rr.Code)
		}
	})
}
