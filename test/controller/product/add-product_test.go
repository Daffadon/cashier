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
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddProduct_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(test.MockProductService)

	reqBody := &bytes.Buffer{}
	formWriter := multipart.NewWriter(reqBody)
	_ = formWriter.WriteField("barcode_id", "1")
	_ = formWriter.WriteField("title", "title-1")
	_ = formWriter.WriteField("price", "1000")
	_ = formWriter.WriteField("description", "desc-1")

	fileWriter, _ := formWriter.CreateFormFile("image", "test.jpg")
	fileWriter.Write([]byte("fake image data"))
	formWriter.Close()

	request := httptest.NewRequest("POST", "/v1/product", reqBody)
	request.Header.Set("Content-Type", formWriter.FormDataContentType())

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = request

	mockService.On("CreateProductService", mock.MatchedBy(func(req dto.AddProductRequest) bool {
		return req.BarcodeId == "1" &&
			req.Title == "title-1" &&
			req.Price.Equal(decimal.NewFromInt(1000)) &&
			req.Description == "desc-1" &&
			req.Image.Filename == "test.jpg"
	})).Return(nil)

	pc := controller.NewProductController(mockService)
	pc.AddProduct(ctx)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "200")
	assert.Contains(t, w.Body.String(), dto.MESSAGE_SUCCESS_ADD_PRODUCT)

	mockService.AssertExpectations(t)
}

func TestAddProduct_BadRequestMissing(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(test.MockProductService)

	reqBody := &bytes.Buffer{}
	formWriter := multipart.NewWriter(reqBody)
	_ = formWriter.WriteField("barcode_id", "1")
	_ = formWriter.WriteField("title", "title-1")
	_ = formWriter.WriteField("price", "1000")

	fileWriter, _ := formWriter.CreateFormFile("image", "test.jpg")
	fileWriter.Write([]byte("fake image data"))
	formWriter.Close()

	request := httptest.NewRequest("POST", "/v1/product", reqBody)
	request.Header.Set("Content-Type", formWriter.FormDataContentType())

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = request

	pc := controller.NewProductController(mockService)
	pc.AddProduct(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "400")
	assert.Contains(t, w.Body.String(), dto.ErrBadrequest.Error())
}

func TestAddProduct_BadRequestOverSize(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(test.MockProductService)

	reqBody := &bytes.Buffer{}
	formWriter := multipart.NewWriter(reqBody)
	_ = formWriter.WriteField("barcode_id", "1")
	_ = formWriter.WriteField("title", "title-1")
	_ = formWriter.WriteField("price", "1000")
	_ = formWriter.WriteField("description", "desc-1")

	fileWriter, _ := formWriter.CreateFormFile("image", "test.jpg")
	largeData := make([]byte, 7*1024*1024) // 7MB
	fileWriter.Write(largeData)
	formWriter.Close()

	request := httptest.NewRequest("POST", "/v1/product", reqBody)
	request.Header.Set("Content-Type", formWriter.FormDataContentType())

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = request

	pc := controller.NewProductController(mockService)
	pc.AddProduct(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "400")
	assert.Contains(t, w.Body.String(), dto.ErrLimitSizeExceeded.Error())
}

func TestAddProduct_BadRequestWrongExtension(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(test.MockProductService)

	reqBody := &bytes.Buffer{}
	formWriter := multipart.NewWriter(reqBody)
	_ = formWriter.WriteField("barcode_id", "1")
	_ = formWriter.WriteField("title", "title-1")
	_ = formWriter.WriteField("price", "1000")
	_ = formWriter.WriteField("description", "desc-1")

	fileWriter, _ := formWriter.CreateFormFile("image", "test.webp")
	fileWriter.Write([]byte("fake image data"))
	formWriter.Close()

	request := httptest.NewRequest("POST", "/v1/product", reqBody)
	request.Header.Set("Content-Type", formWriter.FormDataContentType())

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = request

	pc := controller.NewProductController(mockService)
	pc.AddProduct(ctx)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "400")
	assert.Contains(t, w.Body.String(), dto.ErrWrongFileExtension.Error())
}

func TestAddProduct_Conflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(test.MockProductService)

	reqBody := &bytes.Buffer{}
	formWriter := multipart.NewWriter(reqBody)
	_ = formWriter.WriteField("barcode_id", "1")
	_ = formWriter.WriteField("title", "title-1")
	_ = formWriter.WriteField("price", "1000")
	_ = formWriter.WriteField("description", "desc-1")

	fileWriter, _ := formWriter.CreateFormFile("image", "test.jpg")
	fileWriter.Write([]byte("fake image data"))
	formWriter.Close()

	request := httptest.NewRequest("POST", "/v1/product", reqBody)
	request.Header.Set("Content-Type", formWriter.FormDataContentType())

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = request

	mockService.On("CreateProductService", mock.MatchedBy(func(req dto.AddProductRequest) bool {
		return req.BarcodeId == "1" &&
			req.Title == "title-1" &&
			req.Price.Equal(decimal.NewFromInt(1000)) &&
			req.Description == "desc-1" &&
			req.Image.Filename == "test.jpg"
	})).Return(dto.ErrProductExist)

	pc := controller.NewProductController(mockService)
	pc.AddProduct(ctx)

	assert.Equal(t, http.StatusConflict, w.Code)
	assert.Contains(t, w.Body.String(), "409")
	assert.Contains(t, w.Body.String(), dto.ErrProductExist.Error())

	mockService.AssertExpectations(t)
}

func TestAddProduct_ISE(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(test.MockProductService)

	reqBody := &bytes.Buffer{}
	formWriter := multipart.NewWriter(reqBody)
	_ = formWriter.WriteField("barcode_id", "1")
	_ = formWriter.WriteField("title", "title-1")
	_ = formWriter.WriteField("price", "1000")
	_ = formWriter.WriteField("description", "desc-1")

	fileWriter, _ := formWriter.CreateFormFile("image", "test.jpg")
	fileWriter.Write([]byte("fake image data"))
	formWriter.Close()

	request := httptest.NewRequest("POST", "/v1/product", reqBody)
	request.Header.Set("Content-Type", formWriter.FormDataContentType())

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = request

	mockService.On("CreateProductService", mock.MatchedBy(func(req dto.AddProductRequest) bool {
		return req.BarcodeId == "1" &&
			req.Title == "title-1" &&
			req.Price.Equal(decimal.NewFromInt(1000)) &&
			req.Description == "desc-1" &&
			req.Image.Filename == "test.jpg"
	})).Return(errors.New("ISE"))

	pc := controller.NewProductController(mockService)
	pc.AddProduct(ctx)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "500")
	assert.Contains(t, w.Body.String(), "ISE")
	mockService.AssertExpectations(t)
}
