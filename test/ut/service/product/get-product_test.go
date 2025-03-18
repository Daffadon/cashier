package service_test

import (
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

func TestGetProductService_SuccessPageZero(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)
	mockedRepo.On("CountProductsRepository").Return(uint16(30), nil)
	mockedRepo.On("RetrieveProductsRepository", uint16(12), uint16(0)).Return([]entity.Product{
		{BarcodeId: "123", Title: "Product 1", Image: "img1", Price: decimal.NewFromInt32(1000), Description: "Desc 1"},
		{BarcodeId: "345", Title: "Product 2", Image: "img2", Price: decimal.NewFromInt32(2000), Description: "Desc 2"},
	}, nil)
	page := uint16(0)

	result, err := ps.GetProductService(&page)

	assert.NoError(t, err)
	assert.Equal(t, len(result.Products), 2)
	assert.Equal(t, uint16(result.PageMetaData.Page), uint16(1))
	assert.Equal(t, uint16(result.PageMetaData.PrevPage), uint16(1))
	assert.Equal(t, uint16(result.PageMetaData.NextPage), uint16(2))
	assert.Equal(t, uint16(result.PageMetaData.TotalPage), uint16(3))
	mockedRepo.AssertExpectations(t)
}

func TestGetProductService_SuccessPageOverTotal(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)
	mockedRepo.On("CountProductsRepository").Return(uint16(30), nil)
	mockedRepo.On("RetrieveProductsRepository", uint16(12), uint16(24)).Return([]entity.Product{
		{BarcodeId: "123", Title: "Product 1", Image: "img1", Price: decimal.NewFromInt32(1000), Description: "Desc 1"},
		{BarcodeId: "345", Title: "Product 2", Image: "img2", Price: decimal.NewFromInt32(2000), Description: "Desc 2"},
	}, nil)
	page := uint16(5)

	result, err := ps.GetProductService(&page)

	assert.NoError(t, err)
	assert.Equal(t, len(result.Products), 2)
	assert.Equal(t, uint16(result.PageMetaData.Page), uint16(3))
	assert.Equal(t, uint16(result.PageMetaData.PrevPage), uint16(2))
	assert.Equal(t, uint16(result.PageMetaData.NextPage), uint16(3))
	assert.Equal(t, uint16(result.PageMetaData.TotalPage), uint16(3))
	mockedRepo.AssertExpectations(t)
}

func TestGetProductService_CountError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)
	mockedRepo.On("CountProductsRepository").Return(uint16(0), dto.ErrISEProducts)
	page := uint16(1)

	result, err := ps.GetProductService(&page)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), dto.ErrISEProducts.Error())
	assert.Equal(t, len(result.Products), 0)
	mockedRepo.AssertExpectations(t)
}

func TestGetProductService_RetrieveError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)
	mockedRepo.On("CountProductsRepository").Return(uint16(30), nil)
	mockedRepo.On("RetrieveProductsRepository", uint16(12), uint16(0)).Return([]entity.Product{}, dto.ErrISEProducts)
	page := uint16(1)

	result, err := ps.GetProductService(&page)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), dto.ErrISEProducts.Error())
	assert.Equal(t, len(result.Products), 0)
	mockedRepo.AssertExpectations(t)
}

func TestGetProductService_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)
	mockedRepo.On("CountProductsRepository").Return(uint16(0), nil)
	page := uint16(1)

	result, err := ps.GetProductService(&page)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), dto.ErrProductsNotFound.Error())
	assert.Equal(t, len(result.Products), 0)
	mockedRepo.AssertExpectations(t)
}
