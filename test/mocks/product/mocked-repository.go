package test

import (
	"tiga-putra-cashier-be/dto"
	"tiga-putra-cashier-be/entity"

	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) CountProductsRepository() (uint16, error) {
	args := m.Called()
	return args.Get(0).(uint16), args.Error(1)
}
func (m *MockProductRepository) RetrieveProductsRepository(limit, offset uint16) ([]entity.Product, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]entity.Product), args.Error(1)
}
func (m *MockProductRepository) RetrieveProductByBarcodeId(barcodeId *string) (dto.ProductWithoutTimeStamp, bool) {
	args := m.Called(barcodeId)
	return args.Get(0).(dto.ProductWithoutTimeStamp), args.Bool(1)
}
func (m *MockProductRepository) RetrieveProductForSearch(req *dto.SearchProductQuery) ([]entity.Product, error) {
	args := m.Called(req)
	return args.Get(0).([]entity.Product), args.Error(1)
}
func (m *MockProductRepository) RetrieveDeletedProductByBarcodeId(barcodeId *string) (dto.ProductWithoutTimeStamp, bool) {
	args := m.Called(barcodeId)
	return args.Get(0).(dto.ProductWithoutTimeStamp), args.Bool(1)
}
func (m *MockProductRepository) CreateProductRepository(product *entity.Product) error {
	args := m.Called(product)
	return args.Error(0)
}
func (m *MockProductRepository) UpdateProductRepository(barcodeId *string, product *map[string]interface{}) error {
	args := m.Called(barcodeId, product)
	return args.Error(0)
}
func (m *MockProductRepository) UpdateDeletedProductRepository(barcodeId *string) error {
	args := m.Called(barcodeId)
	return args.Error(0)
}
func (m *MockProductRepository) DeleteProductRepository(barcodeId *string) error {
	args := m.Called(barcodeId)
	return args.Error(0)
}
