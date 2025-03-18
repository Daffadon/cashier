package controller_test

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"tiga-putra-cashier-be/controller"
	"tiga-putra-cashier-be/dto"
	test "tiga-putra-cashier-be/test/mocks/product"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateProduct_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(test.MockProductService)

	reqBody := &bytes.Buffer{}
	formWriter := multipart.NewWriter(reqBody)
	_ = formWriter.WriteField("title", "title-2")
	formWriter.Close()

	request := httptest.NewRequest(http.MethodPatch, "/v1/product/1", reqBody)
	request.Header.Set("Content-Type", formWriter.FormDataContentType())

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = request
	ctx.Params = gin.Params{{Key: "barcode_id", Value: "1"}}

	titleExpected := "title-2"
	mockService.On("UpdateProductService", "1", mock.MatchedBy(func(req dto.UpdateProductRequest) bool {
		return *req.Title == titleExpected
	})).Return(nil)
	pc := controller.NewProductController(mockService)
	pc.UpdateProduct(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "200")
	assert.Contains(t, w.Body.String(), dto.MESSAGE_SUCCESS_UPDATE_PRODUCT)

	mockService.AssertExpectations(t)
}

func TestUpdateProduct_BadRequestURI(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(test.MockProductService)

	reqBody := &bytes.Buffer{}
	formWriter := multipart.NewWriter(reqBody)
	_ = formWriter.WriteField("title", "title-2")
	formWriter.Close()

	request := httptest.NewRequest(http.MethodPatch, "/v1/product/1", reqBody)
	request.Header.Set("Content-Type", formWriter.FormDataContentType())

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = request

	pc := controller.NewProductController(mockService)
	pc.UpdateProduct(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "400")
	assert.Contains(t, w.Body.String(), dto.ErrBadrequest.Error())

	mockService.AssertExpectations(t)
}

func TestUpdateProduct_BadRequestBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(test.MockProductService)

	reqBody := &bytes.Buffer{}
	formWriter := multipart.NewWriter(reqBody)
	formWriter.Close()

	request := httptest.NewRequest(http.MethodPatch, "/v1/product/1", reqBody)
	request.Header.Set("Content-Type", formWriter.FormDataContentType())

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = request
	ctx.Params = gin.Params{{Key: "barcode_id", Value: "1"}}

	mockService.On("UpdateProductService", "1", dto.UpdateProductRequest{}).Return(dto.ErrNoChangesRequest)

	pc := controller.NewProductController(mockService)
	pc.UpdateProduct(ctx)

	assert.Equal(t, http.StatusNotModified, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateProduct_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(test.MockProductService)

	reqBody := &bytes.Buffer{}
	formWriter := multipart.NewWriter(reqBody)
	_ = formWriter.WriteField("title", "title-2")
	formWriter.Close()

	request := httptest.NewRequest(http.MethodPatch, "/v1/product/1", reqBody)
	request.Header.Set("Content-Type", formWriter.FormDataContentType())

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = request
	ctx.Params = gin.Params{{Key: "barcode_id", Value: "1"}}

	titleExpected := "title-2"
	mockService.On("UpdateProductService", "1", mock.MatchedBy(func(req dto.UpdateProductRequest) bool {
		return *req.Title == titleExpected
	})).Return(dto.ErrProductDoesntExist)
	pc := controller.NewProductController(mockService)
	pc.UpdateProduct(ctx)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "404")
	assert.Contains(t, w.Body.String(), dto.ErrProductDoesntExist.Error())

	mockService.AssertExpectations(t)
}

func TestUpdateProduct_ISE(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(test.MockProductService)

	reqBody := &bytes.Buffer{}
	formWriter := multipart.NewWriter(reqBody)
	_ = formWriter.WriteField("title", "title-2")
	formWriter.Close()

	request := httptest.NewRequest(http.MethodPatch, "/v1/product/1", reqBody)
	request.Header.Set("Content-Type", formWriter.FormDataContentType())

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = request
	ctx.Params = gin.Params{{Key: "barcode_id", Value: "1"}}

	titleExpected := "title-2"
	mockService.On("UpdateProductService", "1", mock.MatchedBy(func(req dto.UpdateProductRequest) bool {
		return *req.Title == titleExpected
	})).Return(errors.New("ISE"))
	pc := controller.NewProductController(mockService)
	pc.UpdateProduct(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "500")
	assert.Contains(t, w.Body.String(), "ISE")

	mockService.AssertExpectations(t)
}
