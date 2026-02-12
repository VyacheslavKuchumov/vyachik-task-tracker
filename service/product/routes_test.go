package product

import (
	"VyacheslavKuchumov/test-backend/types"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestProductServiceHandlers(t *testing.T) {
	productStore := &mockProductStore{}
	handler := NewHandler(productStore)

	t.Run("should fail if the product payload is invalid", func(t *testing.T) {
		payload := types.CreateProductPayload{
			Name:        "incorrect product",
			Description: "dfg",
			Image:       "dfg",
			Price:       -1,
			Quantity:    -10,
		}

		marshaled, _ := json.Marshal(payload)
		req, err := http.NewRequest(http.MethodPost, "/create", bytes.NewBuffer(marshaled))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.HandleFunc("/create", handler.HandleCreateProduct)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("should get all products correctly", func(t *testing.T) {

		req, err := http.NewRequest(http.MethodGet, "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router := chi.NewRouter()

		router.HandleFunc("/", handler.HandleGetProducts)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
		}
	})
}

type mockProductStore struct{}

func (m *mockProductStore) CreateProduct(types.CreateProductPayload) error {
	return nil
}

func (m *mockProductStore) GetProducts() ([]*types.Product, error) {
	return nil, nil
}
