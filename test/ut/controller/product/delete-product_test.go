package controller_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"tiga-putra-cashier-be/controller"
	"tiga-putra-cashier-be/dto"
	test "tiga-putra-cashier-be/test/mocks/product"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDeleteProduct_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(test.MockProductService)

	barcodeId := "1"
	mockService.On("DeleteProductService", &barcodeId).Return(nil)
	pc := controller.NewProductController(mockService)
	req, _ := http.NewRequest(http.MethodDelete, "/v1/product/1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{{Key: "barcode_id", Value: "1"}}
	pc.DeleteProduct(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "200")
	assert.Contains(t, w.Body.String(), dto.MESSAGE_SUCCESS_DELETE_PRODUCT)
	mockService.AssertExpectations(t)
}

func TestDeleteProduct_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(test.MockProductService)

	pc := controller.NewProductController(mockService)
	req, _ := http.NewRequest(http.MethodDelete, "/v1/product/1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	pc.DeleteProduct(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "400")
	assert.Contains(t, w.Body.String(), dto.ErrBadrequest.Error())
}

func TestDeleteProduct_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(test.MockProductService)

	barcodeId := "1"
	mockService.On("DeleteProductService", &barcodeId).Return(dto.ErrProductDoesntExist)
	pc := controller.NewProductController(mockService)
	req, _ := http.NewRequest(http.MethodDelete, "/v1/product/1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{{Key: "barcode_id", Value: "1"}}
	pc.DeleteProduct(ctx)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "404")
	assert.Contains(t, w.Body.String(), dto.ErrProductDoesntExist.Error())
	mockService.AssertExpectations(t)
}

func TestDeleteProduct_ISE(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(test.MockProductService)

	barcodeId := "1"
	mockService.On("DeleteProductService", &barcodeId).Return(errors.New("ISE"))
	pc := controller.NewProductController(mockService)
	req, _ := http.NewRequest(http.MethodDelete, "/v1/product/1", nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = req
	ctx.Params = gin.Params{{Key: "barcode_id", Value: "1"}}
	pc.DeleteProduct(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "500")
	assert.Contains(t, w.Body.String(), "ISE")
	mockService.AssertExpectations(t)
}
