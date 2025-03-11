package controller_test

import (
	"encoding/json"
	"errors"
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

func TestSearchProduct_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(test.MockProductService)

	var title = "title-1"
	var reqQuery dto.SearchProductQuery = dto.SearchProductQuery{
		Title: &title,
	}
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
	var actualResponse struct {
		Data []dto.ProductWithoutTimeStamp `json:"data"`
	}
	mockService.On("SearchProductService", &reqQuery).Return(products, nil)
	pc := controller.NewProductController(mockService)
	req, _ := http.NewRequest("GET", "/v1/product/search?title=title-1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	pc.SearchProduct(ctx)

	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "200")
	assert.Contains(t, w.Body.String(), dto.MESSAGE_SUCCESS_SEARCH_PRODUCTS)
	assert.Equal(t, products, actualResponse.Data)
	mockService.AssertExpectations(t)
}

func TestSearchProduct_BadRequestNoQuery(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(test.MockProductService)

	pc := controller.NewProductController(mockService)
	req, _ := http.NewRequest("GET", "/v1/product/search", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	pc.SearchProduct(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "400")
	assert.Contains(t, w.Body.String(), dto.ErrBadrequest.Error())
}

func TestSearchProduct_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(test.MockProductService)

	var title = "title-1"
	var reqQuery dto.SearchProductQuery = dto.SearchProductQuery{
		Title: &title,
	}

	mockService.On("SearchProductService", &reqQuery).Return([]dto.ProductWithoutTimeStamp{}, dto.ErrProductsNotFound)
	pc := controller.NewProductController(mockService)
	req, _ := http.NewRequest("GET", "/v1/product/search?title=title-1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	pc.SearchProduct(ctx)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "404")
	assert.Contains(t, w.Body.String(), dto.ErrProductsNotFound.Error())
	mockService.AssertExpectations(t)
}

func TestSearchProduct_ISE(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(test.MockProductService)

	var title = "title-1"
	var reqQuery dto.SearchProductQuery = dto.SearchProductQuery{
		Title: &title,
	}

	mockService.On("SearchProductService", &reqQuery).Return([]dto.ProductWithoutTimeStamp{}, errors.New("ISE"))
	pc := controller.NewProductController(mockService)
	req, _ := http.NewRequest("GET", "/v1/product/search?title=title-1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	pc.SearchProduct(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "500")
	assert.Contains(t, w.Body.String(), "ISE")
	mockService.AssertExpectations(t)
}
