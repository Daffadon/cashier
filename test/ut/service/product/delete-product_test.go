package service_test

import (
	"errors"
	"testing"
	"tiga-putra-cashier-be/dto"
	"tiga-putra-cashier-be/service"
	testRepo "tiga-putra-cashier-be/test/mocks/product"
	testUtils "tiga-putra-cashier-be/test/mocks/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDeleteProduct_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockedRepo := new(testRepo.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)

	barcodeId := "1"
	mockedRepo.On("RetrieveProductByBarcodeId", &barcodeId).Return(dto.ProductWithoutTimeStamp{}, true)
	mockedRepo.On("DeleteProductRepository", &barcodeId).Return(nil)

	err := ps.DeleteProductService(&barcodeId)
	assert.Nil(t, err)
	mockedRepo.AssertExpectations(t)
}

func TestDeleteProduct_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockedRepo := new(testRepo.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)

	barcodeId := "1"
	mockedRepo.On("RetrieveProductByBarcodeId", &barcodeId).Return(dto.ProductWithoutTimeStamp{}, false)

	err := ps.DeleteProductService(&barcodeId)
	assert.Error(t, err)
	assert.Equal(t, err, dto.ErrProductDoesntExist)
	mockedRepo.AssertExpectations(t)
}

func TestDeleteProduct_ISEDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockedRepo := new(testRepo.MockProductRepository)
	mockedUtils := new(testUtils.MockFileManagement)
	ps := service.NewProductService(mockedRepo, mockedUtils)

	barcodeId := "1"
	mockedRepo.On("RetrieveProductByBarcodeId", &barcodeId).Return(dto.ProductWithoutTimeStamp{}, true)
	mockedRepo.On("DeleteProductRepository", &barcodeId).Return(errors.New("ISE"))

	err := ps.DeleteProductService(&barcodeId)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "ISE")
	mockedRepo.AssertExpectations(t)
}
