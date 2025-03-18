package service_test

import (
	"errors"
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

func TestSearchProduct_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)
	barcodeId := "1"
	req := dto.SearchProductQuery{
		BarcodeId: &barcodeId,
	}
	mockedRepo.On("RetrieveProductForSearch", &req).Return([]entity.Product{
		{BarcodeId: "1", Title: "Product 1", Image: "img1", Price: decimal.NewFromInt32(1000), Description: "Desc 1"},
	}, nil)

	result, err := ps.SearchProductService(&req)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, len(result), 1)
	assert.Equal(t, result[0].BarcodeId, "1")
	mockedRepo.AssertExpectations(t)
}

func TestSearchProduct_RetrieveISE(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)
	barcodeId := "1"
	req := dto.SearchProductQuery{
		BarcodeId: &barcodeId,
	}
	mockedRepo.On("RetrieveProductForSearch", &req).Return([]entity.Product{}, errors.New("ISE"))

	result, err := ps.SearchProductService(&req)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "ISE")
	assert.Equal(t, len(result), 0)
}

func TestSearchProduct_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)
	barcodeId := "1"
	req := dto.SearchProductQuery{
		BarcodeId: &barcodeId,
	}
	mockedRepo.On("RetrieveProductForSearch", &req).Return([]entity.Product{}, nil)

	result, err := ps.SearchProductService(&req)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), dto.ErrProductsNotFound.Error())
	assert.Equal(t, len(result), 0)
	mockedRepo.AssertExpectations(t)
}
