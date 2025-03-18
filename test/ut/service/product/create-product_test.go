package service_test

import (
	"errors"
	"mime/multipart"
	"testing"
	"tiga-putra-cashier-be/dto"
	"tiga-putra-cashier-be/entity"
	"tiga-putra-cashier-be/service"
	testProduct "tiga-putra-cashier-be/test/mocks/product"
	testUtils "tiga-putra-cashier-be/test/mocks/utils"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct_SuccessAdd(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)

	req := dto.AddProductRequest{
		BarcodeId: "1",
		Image: &multipart.FileHeader{
			Filename: "image-1.jpg",
			Size:     1000,
		},
		Title:       "title-1",
		Price:       decimal.NewFromInt32(1000),
		Description: "desc-1",
	}
	newFilename := "generated-1.jpg"
	pathDir := "assets/image"

	mockedUtils.On("GetFileNameExtension", req.Image.Filename).Return("jpg")
	mockedUtils.On("GenerateNewFileName", "jpg").Return("generated-1.jpg")
	mockedRepo.On("RetrieveDeletedProductByBarcodeId", &req.BarcodeId).Return(dto.ProductWithoutTimeStamp{}, false)
	mockedRepo.On("RetrieveProductByBarcodeId", &req.BarcodeId).Return(dto.ProductWithoutTimeStamp{}, false)
	mockedUtils.On("UploadFile", req.Image, newFilename, pathDir).Return(nil)

	newReq := entity.Product{
		BarcodeId:   "1",
		Image:       newFilename,
		Title:       "title-1",
		Price:       decimal.NewFromInt32(1000),
		Description: "desc-1",
	}
	mockedRepo.On("CreateProductRepository", &newReq).Return(nil)

	err := ps.CreateProductService(req)

	assert.Nil(t, err)
	mockedRepo.AssertExpectations(t)
	mockedUtils.AssertExpectations(t)
}

func TestCreateProduct_SuccessRestore(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)

	req := dto.AddProductRequest{
		BarcodeId: "1",
		Image: &multipart.FileHeader{
			Filename: "image-1.jpg",
			Size:     1000,
		},
		Title:       "title-1",
		Price:       decimal.NewFromInt32(1000),
		Description: "desc-1",
	}
	mockedRepo.On("RetrieveDeletedProductByBarcodeId", &req.BarcodeId).Return(dto.ProductWithoutTimeStamp{}, true)
	mockedRepo.On("UpdateDeletedProductRepository", &req.BarcodeId).Return(nil)
	err := ps.CreateProductService(req)

	assert.Nil(t, err)
	mockedRepo.AssertExpectations(t)
}

func TestCreateProduct_ISEUpdate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)

	req := dto.AddProductRequest{
		BarcodeId: "1",
		Image: &multipart.FileHeader{
			Filename: "image-1.jpg",
			Size:     1000,
		},
		Title:       "title-1",
		Price:       decimal.NewFromInt32(1000),
		Description: "desc-1",
	}
	mockedRepo.On("RetrieveDeletedProductByBarcodeId", &req.BarcodeId).Return(dto.ProductWithoutTimeStamp{}, true)
	mockedRepo.On("UpdateDeletedProductRepository", &req.BarcodeId).Return(errors.New("ISE"))
	err := ps.CreateProductService(req)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "ISE")
	mockedRepo.AssertExpectations(t)
}

func TestCreateProduct_ISESaveFile(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)

	req := dto.AddProductRequest{
		BarcodeId: "1",
		Image: &multipart.FileHeader{
			Filename: "image-1.jpg",
			Size:     1000,
		},
		Title:       "title-1",
		Price:       decimal.NewFromInt32(1000),
		Description: "desc-1",
	}
	newFilename := "generated-1.jpg"
	pathDir := "assets/image"

	mockedUtils.On("GetFileNameExtension", req.Image.Filename).Return("jpg")
	mockedUtils.On("GenerateNewFileName", "jpg").Return("generated-1.jpg")
	mockedRepo.On("RetrieveDeletedProductByBarcodeId", &req.BarcodeId).Return(dto.ProductWithoutTimeStamp{}, false)
	mockedRepo.On("RetrieveProductByBarcodeId", &req.BarcodeId).Return(dto.ProductWithoutTimeStamp{}, false)
	mockedUtils.On("UploadFile", req.Image, newFilename, pathDir).Return(dto.ErrToSaveFile)

	err := ps.CreateProductService(req)

	assert.Error(t, err)
	assert.Equal(t, err, dto.ErrToSaveFile)
	mockedRepo.AssertExpectations(t)
	mockedUtils.AssertExpectations(t)
}

func TestCreateProduct_ISEAddProduct(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)

	req := dto.AddProductRequest{
		BarcodeId: "1",
		Image: &multipart.FileHeader{
			Filename: "image-1.jpg",
			Size:     1000,
		},
		Title:       "title-1",
		Price:       decimal.NewFromInt32(1000),
		Description: "desc-1",
	}
	newFilename := "generated-1.jpg"
	pathDir := "assets/image"

	mockedUtils.On("GetFileNameExtension", req.Image.Filename).Return("jpg")
	mockedUtils.On("GenerateNewFileName", "jpg").Return("generated-1.jpg")
	mockedRepo.On("RetrieveDeletedProductByBarcodeId", &req.BarcodeId).Return(dto.ProductWithoutTimeStamp{}, false)
	mockedRepo.On("RetrieveProductByBarcodeId", &req.BarcodeId).Return(dto.ProductWithoutTimeStamp{}, false)
	mockedUtils.On("UploadFile", req.Image, newFilename, pathDir).Return(nil)

	newReq := entity.Product{
		BarcodeId:   "1",
		Image:       newFilename,
		Title:       "title-1",
		Price:       decimal.NewFromInt32(1000),
		Description: "desc-1",
	}
	mockedRepo.On("CreateProductRepository", &newReq).Return(dto.ErrToAddProduct)

	err := ps.CreateProductService(req)

	assert.Error(t, err)
	assert.Equal(t, err, dto.ErrToAddProduct)
	mockedRepo.AssertExpectations(t)
	mockedUtils.AssertExpectations(t)
}

func TestCreateProduct_BadRequestExtension(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)

	req := dto.AddProductRequest{
		BarcodeId: "1",
		Image: &multipart.FileHeader{
			Filename: "image-1.webp",
			Size:     1000,
		},
		Title:       "title-1",
		Price:       decimal.NewFromInt32(1000),
		Description: "desc-1",
	}
	mockedRepo.On("RetrieveDeletedProductByBarcodeId", &req.BarcodeId).Return(dto.ProductWithoutTimeStamp{}, false)
	mockedUtils.On("GetFileNameExtension", req.Image.Filename).Return("webp")
	err := ps.CreateProductService(req)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), dto.ErrWrongFileExtension.Error())
	mockedRepo.AssertExpectations(t)
	mockedUtils.AssertExpectations(t)
}

func TestCreateProduct_BadRequestSize(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)

	req := dto.AddProductRequest{
		BarcodeId: "1",
		Image: &multipart.FileHeader{
			Filename: "image-1.jpg",
			Size:     7 * 1024 * 1024,
		},
		Title:       "title-1",
		Price:       decimal.NewFromInt32(1000),
		Description: "desc-1",
	}
	mockedRepo.On("RetrieveDeletedProductByBarcodeId", &req.BarcodeId).Return(dto.ProductWithoutTimeStamp{}, false)
	mockedUtils.On("GetFileNameExtension", req.Image.Filename).Return("jpg")
	err := ps.CreateProductService(req)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), dto.ErrLimitSizeExceeded.Error())
	mockedRepo.AssertExpectations(t)
	mockedUtils.AssertExpectations(t)
}

func TestCreateProduct_ConflictExist(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)

	req := dto.AddProductRequest{
		BarcodeId: "1",
		Image: &multipart.FileHeader{
			Filename: "image-1.jpg",
			Size:     1000,
		},
		Title:       "title-1",
		Price:       decimal.NewFromInt32(1000),
		Description: "desc-1",
	}
	mockedRepo.On("RetrieveDeletedProductByBarcodeId", &req.BarcodeId).Return(dto.ProductWithoutTimeStamp{}, false)
	mockedUtils.On("GetFileNameExtension", req.Image.Filename).Return("jpg")
	mockedRepo.On("RetrieveProductByBarcodeId", &req.BarcodeId).Return(dto.ProductWithoutTimeStamp{}, true)
	err := ps.CreateProductService(req)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), dto.ErrProductExist.Error())
	mockedRepo.AssertExpectations(t)
	mockedUtils.AssertExpectations(t)
}
