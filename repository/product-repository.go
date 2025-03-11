package repository

import (
	"context"
	"fmt"
	"tiga-putra-cashier-be/dto"
	"tiga-putra-cashier-be/entity"
	"tiga-putra-cashier-be/utils"
	"time"

	"gorm.io/gorm"
)

type (
	ProductRepository interface {
		CountProductsRepository() (uint16, error)
		RetrieveProductsRepository(limit, offset uint16) ([]entity.Product, error)
		RetrieveProductByBarcodeId(barcodeId *string) (dto.ProductWithoutTimeStamp, bool)
		RetrieveProductForSearch(req *dto.SearchProductQuery) ([]entity.Product, error)
		RetrieveDeletedProductByBarcodeId(barcodeId *string) (dto.ProductWithoutTimeStamp, bool)
		CreateProductRepository(product *entity.Product) error
		UpdateProductRepository(barcodeId *string, product *map[string]interface{}) error
		UpdateDeletedProductRepository(barcodeId *string) error
		DeleteProductRepository(barcodeId *string) error
	}
	productRepository struct {
		db *gorm.DB
	}
)

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (p *productRepository) CountProductsRepository() (uint16, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var totalProduct int64
	err := p.db.WithContext(ctx).Model(&entity.Product{}).Count(&totalProduct).Error
	if err != nil {
		return 0, dto.ErrISEProducts
	}
	return uint16(totalProduct), nil
}

func (p *productRepository) RetrieveProductsRepository(limit, offset uint16) ([]entity.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var allProducts []entity.Product
	err := p.db.WithContext(ctx).Scopes(utils.Paginate(limit, offset)).Find(&allProducts).Error
	if err != nil {
		return []entity.Product{}, dto.ErrProductsNotFound
	}
	return allProducts, nil
}

func (p *productRepository) RetrieveProductByBarcodeId(barcodeId *string) (dto.ProductWithoutTimeStamp, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var product dto.ProductWithoutTimeStamp
	err := p.db.WithContext(ctx).Model(&entity.Product{}).Where("barcode_id = ?", *barcodeId).First(&product).Error
	if err != nil {
		if err.Error() == "record not found" {
			return dto.ProductWithoutTimeStamp{}, false
		}
		return dto.ProductWithoutTimeStamp{}, false
	}
	return product, true
}

func (p *productRepository) RetrieveProductForSearch(req *dto.SearchProductQuery) ([]entity.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var products []entity.Product
	query := p.db.WithContext(ctx).Model(&entity.Product{})

	if req.Title != nil && *req.Title != "" {
		query = query.Where("title ILIKE ?", "%"+*req.Title+"%")
	}

	if req.BarcodeId != nil && *req.BarcodeId != "" {
		query = query.Where("barcode_id = ?", *req.BarcodeId)
	}

	err := query.Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (p *productRepository) RetrieveDeletedProductByBarcodeId(barcodeId *string) (dto.ProductWithoutTimeStamp, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var product dto.ProductWithoutTimeStamp
	err := p.db.WithContext(ctx).Unscoped().Model(&entity.Product{}).Where("barcode_id = ? AND deleted_at IS NOT NULL", *barcodeId).First(&product).Error
	if err != nil {
		if err.Error() == "record not found" {
			return dto.ProductWithoutTimeStamp{}, false
		}
		return dto.ProductWithoutTimeStamp{}, false
	}
	fmt.Print(product)
	return product, true
}

func (p *productRepository) CreateProductRepository(product *entity.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := p.db.WithContext(ctx).Create(product).Error
	if err != nil {
		return dto.ErrToAddProduct
	}
	return nil
}

func (p *productRepository) UpdateProductRepository(barcodeId *string, product *map[string]interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := p.db.WithContext(ctx).Model(&entity.Product{}).Where("barcode_id = ?", *barcodeId).Updates(&product).Error
	if err != nil {
		return err
	}
	return nil
}

func (p *productRepository) UpdateDeletedProductRepository(barcodeId *string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := p.db.WithContext(ctx).Unscoped().Model(&entity.Product{}).Where("barcode_id = ?", *barcodeId).Update("deleted_at", nil).Error
	if err != nil {
		fmt.Print(err.Error())
		return err
	}
	return nil
}

func (p *productRepository) DeleteProductRepository(barcodeId *string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := p.db.WithContext(ctx).Where("barcode_id = ?", barcodeId).Delete(&entity.Product{})
	if err != nil {
		return err.Error
	}
	return nil
}
