package controller_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"tiga-putra-cashier-be/controller"
	"tiga-putra-cashier-be/dto"
	test "tiga-putra-cashier-be/test/mocks/product"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProduct_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(test.MockProductService)

	var page uint16 = 1
	var products []dto.ProductWithoutTimeStamp = []dto.ProductWithoutTimeStamp{
		{
			BarcodeId:   "1",
			Image:       "imageurl-1",
			Title:       "data-1",
			Description: "data-1-description",
			Price:       decimal.NewFromInt(1000),
		},
		{
			BarcodeId:   "2",
			Image:       "imageurl-2",
			Title:       "data-2",
			Description: "data-2-description",
			Price:       decimal.NewFromInt(2000),
		},
	}
	var pagination dto.PaginationResponse = dto.PaginationResponse{
		Page:      1,
		PrevPage:  1,
		NextPage:  1,
		TotalPage: 1,
	}
	expectedProducts := dto.AllProductsWithPagination{
		Products:     products,
		PageMetaData: pagination,
	}

	var actualResponse struct {
		Data dto.AllProductsWithPagination `json:"data"`
	}

	mockService.On("GetProductService", &page).Return(expectedProducts, nil)
	pc := controller.NewProductController(mockService)
	req, _ := http.NewRequest("GET", "/v1/product?page=1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	pc.GetProduct(ctx)

	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "200")
	assert.Contains(t, w.Body.String(), dto.MESSAGE_SUCCESS_GET_ALL_PRODUCTS)
	assert.Equal(t, expectedProducts, actualResponse.Data)
	mockService.AssertExpectations(t)
}

func TestGetProduct_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(test.MockProductService)
	var page uint16 = 1
	mockService.On("GetProductService", &page).Return(dto.AllProductsWithPagination{}, dto.ErrProductsNotFound)
	req, _ := http.NewRequest("GET", "/v1/product?page=1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	pc := controller.NewProductController(mockService)
	pc.GetProduct(ctx)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "404")
	assert.Contains(t, w.Body.String(), dto.ErrProductsNotFound.Error())
	mockService.AssertExpectations(t)
}
func TestGetProduct_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockService := new(test.MockProductService)
	pc := controller.NewProductController(mockService)
	req, _ := http.NewRequest("GET", "/v1/product", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	pc.GetProduct(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "400")
	assert.Contains(t, w.Body.String(), dto.ErrBadrequest.Error())
}

func TestGetProduct_InternalServerError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(test.MockProductService)
	var page uint16 = 1
	mockService.On("GetProductService", &page).Return(dto.AllProductsWithPagination{}, dto.ErrISEProducts)
	req, _ := http.NewRequest("GET", "/v1/product?page=1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	pc := controller.NewProductController(mockService)
	pc.GetProduct(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "500")
	assert.Contains(t, w.Body.String(), dto.ErrISEProducts.Error())
	mockService.AssertExpectations(t)
}
