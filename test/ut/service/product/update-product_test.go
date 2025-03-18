package service_test

import (
	"errors"
	"mime/multipart"
	"testing"
	"tiga-putra-cashier-be/dto"
	"tiga-putra-cashier-be/service"
	testProduct "tiga-putra-cashier-be/test/mocks/product"
	testUtils "tiga-putra-cashier-be/test/mocks/utils"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestUpdateProduct_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)

	barcodeId := "1"
	title := "title-update"
	image := multipart.FileHeader{
		Filename: "image,jpg",
		Size:     1000,
	}
	price := decimal.NewFromInt32(1000)
	description := "desc-1"
	product := dto.UpdateProductRequest{
		Title:       &title,
		Image:       &image,
		Price:       &price,
		Description: &description,
	}

	newFilename := "generated-1.jpg"
	pathDir := "assets/image"
	pathDeletedImage := "assets/image/deleted-1.jpg"

	updates := make(map[string]interface{})
	updates["title"] = *product.Title
	updates["price"] = *product.Price
	updates["description"] = *product.Description
	updates["image"] = newFilename
	mockedRepo.On("RetrieveProductByBarcodeId", &barcodeId).Return(dto.ProductWithoutTimeStamp{
		Image: "deleted-1.jpg",
	}, true)
	mockedUtils.On("GetFileNameExtension", product.Image.Filename).Return("jpg")
	mockedUtils.On("GenerateNewFileName", "jpg").Return("generated-1.jpg")
	mockedUtils.On("UploadFile", product.Image, newFilename, pathDir).Return(nil)
	mockedUtils.On("DeleteFile", pathDeletedImage).Return(nil)
	mockedRepo.On("UpdateProductRepository", &barcodeId, &updates).Return(nil)

	err := ps.UpdateProductService(barcodeId, product)

	assert.Nil(t, err)
	mockedRepo.AssertExpectations(t)
	mockedUtils.AssertExpectations(t)
}

func TestUpdateProduct_NoChanges(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)

	barcodeId := "1"
	product := dto.UpdateProductRequest{}
	mockedRepo.On("RetrieveProductByBarcodeId", &barcodeId).Return(dto.ProductWithoutTimeStamp{}, true)

	err := ps.UpdateProductService(barcodeId, product)

	assert.Error(t, err)
	assert.Equal(t, err, dto.ErrNoChangesRequest)
	mockedRepo.AssertExpectations(t)
}

func TestUpdateProduct_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)

	barcodeId := "1"
	product := dto.UpdateProductRequest{}
	mockedRepo.On("RetrieveProductByBarcodeId", &barcodeId).Return(dto.ProductWithoutTimeStamp{}, false)

	err := ps.UpdateProductService(barcodeId, product)

	assert.Error(t, err)
	assert.Equal(t, err, dto.ErrProductDoesntExist)
	mockedRepo.AssertExpectations(t)
}

func TestUpdateProduct_BadRequestExtension(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)

	barcodeId := "1"
	product := dto.UpdateProductRequest{
		Image: &multipart.FileHeader{
			Filename: "img-update.webp",
			Size:     1000,
		},
	}
	mockedRepo.On("RetrieveProductByBarcodeId", &barcodeId).Return(dto.ProductWithoutTimeStamp{}, true)
	mockedUtils.On("GetFileNameExtension", product.Image.Filename).Return("webp")

	err := ps.UpdateProductService(barcodeId, product)

	assert.Error(t, err)
	assert.Equal(t, err, dto.ErrWrongFileExtension)
	mockedRepo.AssertExpectations(t)
	mockedUtils.AssertExpectations(t)
}

func TestUpdateProduct_BadRequestSize(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)

	barcodeId := "1"
	product := dto.UpdateProductRequest{
		Image: &multipart.FileHeader{
			Filename: "img-update.jpg",
			Size:     7 * 1024 * 1024,
		},
	}
	mockedRepo.On("RetrieveProductByBarcodeId", &barcodeId).Return(dto.ProductWithoutTimeStamp{}, true)
	mockedUtils.On("GetFileNameExtension", product.Image.Filename).Return("jpg")

	err := ps.UpdateProductService(barcodeId, product)

	assert.Error(t, err)
	assert.Equal(t, err, dto.ErrLimitSizeExceeded)
	mockedRepo.AssertExpectations(t)
	mockedUtils.AssertExpectations(t)
}

func TestUpdateProduct_ISESaveFile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)

	barcodeId := "1"
	product := dto.UpdateProductRequest{
		Image: &multipart.FileHeader{
			Filename: "img-update.jpg",
			Size:     1000,
		},
	}

	newFilename := "generated-1.jpg"
	pathDir := "assets/image"

	mockedRepo.On("RetrieveProductByBarcodeId", &barcodeId).Return(dto.ProductWithoutTimeStamp{}, true)
	mockedUtils.On("GetFileNameExtension", product.Image.Filename).Return("jpg")
	mockedUtils.On("GenerateNewFileName", "jpg").Return("generated-1.jpg")
	mockedUtils.On("UploadFile", product.Image, newFilename, pathDir).Return(dto.ErrToSaveFile)

	err := ps.UpdateProductService(barcodeId, product)

	assert.Error(t, err)
	assert.Equal(t, err, dto.ErrToSaveFile)
	mockedRepo.AssertExpectations(t)
	mockedUtils.AssertExpectations(t)
}

func TestUpdateProduct_ISEDeleteFile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)

	barcodeId := "1"
	product := dto.UpdateProductRequest{
		Image: &multipart.FileHeader{
			Filename: "img-update.jpg",
			Size:     1000,
		},
	}

	newFilename := "generated-1.jpg"
	pathDir := "assets/image"
	pathDeletedImage := "assets/image/deleted-1.jpg"

	mockedRepo.On("RetrieveProductByBarcodeId", &barcodeId).Return(dto.ProductWithoutTimeStamp{
		Image: "deleted-1.jpg",
	}, true)
	mockedUtils.On("GetFileNameExtension", product.Image.Filename).Return("jpg")
	mockedUtils.On("GenerateNewFileName", "jpg").Return("generated-1.jpg")
	mockedUtils.On("UploadFile", product.Image, newFilename, pathDir).Return(nil)
	mockedUtils.On("DeleteFile", pathDeletedImage).Return(errors.New("ISE"))
	err := ps.UpdateProductService(barcodeId, product)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "ISE")
	mockedRepo.AssertExpectations(t)
	mockedUtils.AssertExpectations(t)
}

func TestUpdateProduct_ISEUpdateProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)

	barcodeId := "1"
	title := "title-update"
	product := dto.UpdateProductRequest{
		Title: &title,
	}
	updates := make(map[string]interface{})
	updates["title"] = *product.Title
	mockedRepo.On("RetrieveProductByBarcodeId", &barcodeId).Return(dto.ProductWithoutTimeStamp{}, true)
	mockedRepo.On("UpdateProductRepository", &barcodeId, &updates).Return(errors.New("ISE"))

	err := ps.UpdateProductService(barcodeId, product)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "ISE")
	mockedRepo.AssertExpectations(t)
	mockedUtils.AssertExpectations(t)
}
