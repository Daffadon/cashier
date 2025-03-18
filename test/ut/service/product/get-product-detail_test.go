package service_test

import (
	"testing"
	"tiga-putra-cashier-be/dto"
	"tiga-putra-cashier-be/service"
	testProduct "tiga-putra-cashier-be/test/mocks/product"
	testUtils "tiga-putra-cashier-be/test/mocks/utils"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestGetProductDetail_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)
	barcodeId := "1"
	mockedRepo.On("RetrieveProductByBarcodeId", &barcodeId).Return(dto.ProductWithoutTimeStamp{
		BarcodeId: "1", Image: "image-1", Title: "title-1", Price: decimal.NewFromInt32(1000), Description: "desc-1",
	}, true)
	result, err := ps.GetProductDetailService(&barcodeId)
	assert.Nil(t, err)
	assert.Equal(t, result.BarcodeId, "1")
}

func TestGetProductDetail_RetrieveError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockedRepo := new(testProduct.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)
	barcodeId := "1"
	mockedRepo.On("RetrieveProductByBarcodeId", &barcodeId).Return(dto.ProductWithoutTimeStamp{}, false)
	result, err := ps.GetProductDetailService(&barcodeId)

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), dto.ErrProductDoesntExist.Error())
	assert.Equal(t, result, dto.ProductWithoutTimeStamp{})
}
