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

func TestGetProductDetail_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(test.MockProductService)

	barcodeId := "1"
	var product dto.ProductWithoutTimeStamp = dto.ProductWithoutTimeStamp{
		BarcodeId:   "1",
		Image:       "image-1",
		Title:       "title-1",
		Price:       decimal.NewFromInt(1000),
		Description: "description-1",
	}
	mockService.On("GetProductDetailService", &barcodeId).Return(product, nil)
	pc := controller.NewProductController(mockService)
	req, _ := http.NewRequest("GET", "/v1/product/1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{{Key: "barcode_id", Value: "1"}}
	pc.GetProductDetail(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "200")
	assert.Contains(t, w.Body.String(), dto.MESSAGE_SUCCESS_GET_PRODUCT_DETAIL)

	var actualResponse struct {
		Data dto.ProductWithoutTimeStamp `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
	require.NoError(t, err)

	assert.Equal(t, actualResponse.Data, product)
	mockService.AssertExpectations(t)
}

func TestGetProductDetail_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(test.MockProductService)
	pc := controller.NewProductController(mockService)
	req, _ := http.NewRequest("GET", "/v1/product", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	pc.GetProductDetail(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "400")
	assert.Contains(t, w.Body.String(), dto.ErrBadrequest.Error())
}

func TestGetProductDetail_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(test.MockProductService)
	pc := controller.NewProductController(mockService)

	barcodeId := "1"
	mockService.On("GetProductDetailService", &barcodeId).Return(dto.ProductWithoutTimeStamp{}, dto.ErrProductDoesntExist)
	req, _ := http.NewRequest("GET", "/v1/product/1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{{Key: "barcode_id", Value: "1"}}
	pc.GetProductDetail(ctx)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "404")
	assert.Contains(t, w.Body.String(), dto.ErrProductDoesntExist.Error())
	mockService.AssertExpectations(t)
}
