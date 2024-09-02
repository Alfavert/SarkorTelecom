package handler

import (
	"SarkorTelekom/pkg/repository"
	"SarkorTelekom/pkg/service"
	mock_service "SarkorTelekom/pkg/service/mocks"
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestHandler_CreateProduct(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockProduct, products repository.Products)
	tests := []struct {
		name                 string
		inputBody            string
		inputUser            repository.Products
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ok",
			inputBody: `{"name":"test", "description":"test description", "price":1, "quantity": 10}`,
			inputUser: repository.Products{
				Name:        "test",
				Description: "test description",
				Price:       1,
				Quantity:    10,
			},
			mockBehavior: func(r *mock_service.MockProduct, products repository.Products) {
				r.EXPECT().CreateProduct(products).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"id":1}`,
		}, {
			name:                 "Wrong Input",
			inputBody:            `{"name": "wrong input"}`,
			inputUser:            repository.Products{},
			mockBehavior:         func(r *mock_service.MockProduct, products repository.Products) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		}, {
			name:      "Service Error",
			inputBody: `{"name":"test", "description":"test description", "price":1, "quantity": 10}`,
			inputUser: repository.Products{
				Name:        "test",
				Description: "test description",
				Price:       1,
				Quantity:    10,
			},
			mockBehavior: func(r *mock_service.MockProduct, products repository.Products) {
				r.EXPECT().CreateProduct(products).Return(0, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockProduct(c)
			test.mockBehavior(repo, test.inputUser)

			services := &service.Service{Product: repo}
			handler := Handler{services}

			// Init Endpoint
			r := gin.New()
			r.POST("/products", handler.CreateProduct)

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/products",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_getProduct(t *testing.T) {

	type mockBehavior func(r *mock_service.MockProduct, products repository.Products)

	tests := []struct {
		name               string
		id                 string
		mockBehavior       func(s *mock_service.MockProduct)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "Valid ID",
			id:   "1",
			mockBehavior: func(s *mock_service.MockProduct) {
				s.EXPECT().GetById(1).Return(repository.Products{
					Id:          1,
					Name:        "Test Product",
					Description: "Test Description",
					Price:       10.99,
					Quantity:    100,
				}, nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"id":1,"name":"Test Product","description":"Test Description","price":10.99,"quantity":100}`,
		},
		{
			name:               "Invalid ID",
			id:                 "abc",
			mockBehavior:       func(s *mock_service.MockProduct) {},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"message":"invalid id"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock_service.NewMockProduct(ctrl)
			test.mockBehavior(repo)

			services := &service.Service{Product: repo}
			handler := &Handler{services: services}

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.GET("/product", handler.getProduct)

			req, _ := http.NewRequest(http.MethodGet, "/product?id="+test.id, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponse, w.Body.String())
		})
	}
}

func TestHandler_UpdateProduct(t *testing.T) {
	type mockBehavior func(r *mock_service.MockProduct, id int, product repository.UpdateProducts)

	tests := []struct {
		name                 string
		inputBody            string
		id                   int
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "ok",
			inputBody: `{"name":"test", "description":"test description", "price":1, "quantity": 10}`,
			id:        1,
			mockBehavior: func(r *mock_service.MockProduct, id int, product repository.UpdateProducts) {
				r.EXPECT().Update(id, repository.UpdateProducts{
					Name:        strPtr("test"),
					Description: strPtr("test description"),
					Price:       float64Ptr(1),
					Quantity:    intPtr(10),
				}).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"success"}`,
		},
		{
			name:                 "Invalid Input Body",
			inputBody:            `{"name": "test", "price": 1}`, // missing "quantity"
			id:                   1,
			mockBehavior:         func(r *mock_service.MockProduct, id int, product repository.UpdateProducts) {},
			expectedStatusCode:   400,
			expectedResponseBody: `{"message":"invalid input body"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"name":"test", "description":"test description", "price":1, "quantity": 10}`,
			id:        1,
			mockBehavior: func(r *mock_service.MockProduct, id int, product repository.UpdateProducts) {
				r.EXPECT().Update(id, repository.UpdateProducts{
					Name:        strPtr("test"),
					Description: strPtr("test description"),
					Price:       float64Ptr(1),
					Quantity:    intPtr(10),
				}).Return(errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockProduct(c)
			test.mockBehavior(repo, test.id, repository.UpdateProducts{
				Name:        strPtr("test"),
				Description: strPtr("test description"),
				Price:       float64Ptr(1),
				Quantity:    intPtr(10),
			})

			services := &service.Service{Product: repo}
			handler := &Handler{services: services}

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.PUT("/product", handler.updateProduct)

			reqBody := bytes.NewBufferString(test.inputBody)
			req, _ := http.NewRequest(http.MethodPut, "/product?id="+strconv.Itoa(test.id), reqBody)
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_DeleteProduct(t *testing.T) {
	type mockBehavior func(r *mock_service.MockProduct, id int)

	tests := []struct {
		name                 string
		id                   string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:                 "ok",
			id:                   "1",
			mockBehavior:         func(r *mock_service.MockProduct, id int) { r.EXPECT().Delete(id).Return(nil) },
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"success"}`,
		},
		{
			name: "Service Error",
			id:   "1",
			mockBehavior: func(r *mock_service.MockProduct, id int) {
				r.EXPECT().Delete(id).Return(errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockProduct(c)
			if test.id != "invalid" {
				id, _ := strconv.Atoi(test.id)
				test.mockBehavior(repo, id)
			}

			services := &service.Service{Product: repo}
			handler := &Handler{services: services}

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.DELETE("/product", handler.deleteProduct)

			req, _ := http.NewRequest(http.MethodDelete, "/product?id="+test.id, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.Equal(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func TestHandler_GetProducts(t *testing.T) {
	type mockBehavior func(r *mock_service.MockProduct)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "ok",
			mockBehavior: func(r *mock_service.MockProduct) {
				r.EXPECT().GetAll().Return([]repository.Products{
					{Id: 1, Name: "Product1", Description: "Description1", Price: 10.0, Quantity: 100},
					{Id: 2, Name: "Product2", Description: "Description2", Price: 20.0, Quantity: 200},
				}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"id":1,"name":"Product1","description":"Description1","price":10.0,"quantity":100},{"id":2,"name":"Product2","description":"Description2","price":20.0,"quantity":200}]`,
		},
		{
			name: "Service Error",
			mockBehavior: func(r *mock_service.MockProduct) {
				r.EXPECT().GetAll().Return(nil, errors.New("something went wrong"))
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"message":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockProduct(c)
			test.mockBehavior(repo)

			services := &service.Service{Product: repo}
			handler := &Handler{services: services}

			gin.SetMode(gin.TestMode)
			r := gin.New()
			r.GET("/products", handler.getProducts)

			req, _ := http.NewRequest(http.MethodGet, "/products", nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
			assert.JSONEq(t, test.expectedResponseBody, w.Body.String())
		})
	}
}

func strPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}

func intPtr(i int) *int {
	return &i
}
