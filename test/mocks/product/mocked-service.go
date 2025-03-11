package test

import (
	"tiga-putra-cashier-be/dto"

	"github.com/stretchr/testify/mock"
)

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) GetProductService(page *uint16) (dto.AllProductsWithPagination, error) {
	args := m.Called(page)
	return args.Get(0).(dto.AllProductsWithPagination), args.Error(1)
}
func (m *MockProductService) GetProductDetailService(barcodeId *string) (dto.ProductWithoutTimeStamp, error) {
	args := m.Called(barcodeId)
	return args.Get(0).(dto.ProductWithoutTimeStamp), args.Error(1)
}
func (m *MockProductService) SearchProductService(req *dto.SearchProductQuery) ([]dto.ProductWithoutTimeStamp, error) {
	args := m.Called(req)
	return args.Get(0).([]dto.ProductWithoutTimeStamp), args.Error(1)
}
func (m *MockProductService) CreateProductService(product dto.AddProductRequest) error {
	args := m.Called(product)
	return args.Error(0)
}
func (m *MockProductService) UpdateProductService(barcodeId string, product dto.UpdateProductRequest) error {
	args := m.Called(barcodeId, product)
	return args.Error(0)
}
func (m *MockProductService) DeleteProductService(barcodeId *string) error {
	args := m.Called(barcodeId)
	return args.Error(0)
}
