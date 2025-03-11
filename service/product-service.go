package service

import (
	"fmt"
	"math"
	"tiga-putra-cashier-be/constant"
	"tiga-putra-cashier-be/dto"
	"tiga-putra-cashier-be/entity"
	"tiga-putra-cashier-be/repository"
	"tiga-putra-cashier-be/utils"

	"github.com/google/uuid"
)

type (
	ProductService interface {
		GetProductService(page *uint16) (dto.AllProductsWithPagination, error)
		GetProductDetailService(barcodeId *string) (dto.ProductWithoutTimeStamp, error)
		SearchProductService(req *dto.SearchProductQuery) ([]dto.ProductWithoutTimeStamp, error)
		CreateProductService(product dto.AddProductRequest) error
		UpdateProductService(barcodeId string, product dto.UpdateProductRequest) error
		DeleteProductService(barcodeId *string) error
	}
	productService struct {
		producRepository repository.ProductRepository
	}
)

func NewProductService(productRepository repository.ProductRepository) ProductService {
	return &productService{productRepository}
}

func (p *productService) GetProductService(page *uint16) (dto.AllProductsWithPagination, error) {
	totalProducts, err := p.producRepository.CountProductsRepository()
	if err != nil {
		return dto.AllProductsWithPagination{}, err
	}
	if totalProducts == 0 {
		return dto.AllProductsWithPagination{}, dto.ErrProductsNotFound
	}
	var itemPerPage uint16 = 12
	totalPage := uint16(math.Ceil(float64(totalProducts) / float64(itemPerPage)))
	if *page == 0 {
		*page = 1
	} else if *page > totalPage {
		*page = totalPage
	}

	offset := itemPerPage * (*page - 1)
	allProducts, err := p.producRepository.RetrieveProductsRepository(itemPerPage, offset)
	if err != nil {
		return dto.AllProductsWithPagination{}, err
	}
	var finalProducts []dto.ProductWithoutTimeStamp

	for _, product := range allProducts {
		finalProducts = append(finalProducts, dto.ProductWithoutTimeStamp{
			BarcodeId:   product.BarcodeId,
			Title:       product.Title,
			Image:       product.Image,
			Price:       product.Price,
			Description: product.Description,
		})
	}

	var nextPage, prevPage uint16
	if *page == totalPage {
		nextPage = *page
	} else {
		nextPage = *page + 1
	}

	if *page == 1 {
		prevPage = *page
	} else {
		prevPage = *page - 1
	}

	return dto.AllProductsWithPagination{
		Products: finalProducts,
		PageMetaData: dto.PaginationResponse{
			Page:      *page,
			PrevPage:  prevPage,
			NextPage:  nextPage,
			TotalPage: totalPage,
		},
	}, nil
}

func (p *productService) GetProductDetailService(barcodeId *string) (dto.ProductWithoutTimeStamp, error) {
	productExist, ok := p.producRepository.RetrieveProductByBarcodeId(barcodeId)
	if !ok {
		return dto.ProductWithoutTimeStamp{}, dto.ErrProductDoesntExist
	}
	return productExist, nil
}

func (p *productService) SearchProductService(req *dto.SearchProductQuery) ([]dto.ProductWithoutTimeStamp, error) {
	products, err := p.producRepository.RetrieveProductForSearch(req)
	if err != nil {
		return []dto.ProductWithoutTimeStamp{}, err
	}
	if len(products) < 1 {
		return []dto.ProductWithoutTimeStamp{}, dto.ErrProductsNotFound
	}
	var finalProducts []dto.ProductWithoutTimeStamp
	for _, product := range products {
		finalProducts = append(finalProducts, dto.ProductWithoutTimeStamp{
			BarcodeId:   product.BarcodeId,
			Image:       product.Image,
			Title:       product.Title,
			Price:       product.Price,
			Description: product.Description,
		})
	}
	return finalProducts, nil
}

func (p *productService) CreateProductService(product dto.AddProductRequest) error {
	_, ok := p.producRepository.RetrieveDeletedProductByBarcodeId(&product.BarcodeId)
	if ok {
		fmt.Println(ok)
		if err := p.producRepository.UpdateDeletedProductRepository(&product.BarcodeId); err != nil {
			return err
		}
		return nil
	} else {
		_, ok := p.producRepository.RetrieveProductByBarcodeId(&product.BarcodeId)
		if ok {
			return dto.ErrProductExist
		}
		ext := utils.GetFileNameExtension(product.Image.Filename)
		var newFileName = fmt.Sprintf("%s.%s", uuid.New().String(), ext)
		if ext != "" {
			pathDir := constant.ImageDir
			if err := utils.UploadFile(product.Image, newFileName, pathDir); err != nil {
				return err
			}
		}
		newProduct := entity.Product{
			BarcodeId:   product.BarcodeId,
			Image:       newFileName,
			Title:       product.Title,
			Price:       product.Price,
			Description: product.Description,
		}
		if err := p.producRepository.CreateProductRepository(&newProduct); err != nil {
			return err
		}
		return nil
	}
}

func (p *productService) UpdateProductService(barcodeId string, product dto.UpdateProductRequest) error {
	productExist, ok := p.producRepository.RetrieveProductByBarcodeId(&barcodeId)
	if !ok {
		return dto.ErrProductDoesntExist
	}
	updates := make(map[string]interface{})
	if product.Title != nil {
		updates["title"] = *product.Title
	}
	if product.Price != nil {
		updates["price"] = *product.Price
	}
	if product.Description != nil {
		updates["description"] = *product.Description
	}

	if product.Image != nil {
		ext := utils.GetFileNameExtension(product.Image.Filename)
		var newFileName = fmt.Sprintf("%s.%s", uuid.New().String(), ext)
		if ext != "" {
			pathDir := constant.ImageDir
			if err := utils.UploadFile(product.Image, newFileName, pathDir); err != nil {
				return err
			}
			pathDeletedImage := fmt.Sprintf("%s/%s", constant.ImageDir, productExist.Image)
			err := utils.DeleteFile(pathDeletedImage)
			if err != nil {
				return err
			}
			updates["image"] = newFileName
		}
	}
	if len(updates) > 0 {
		if err := p.producRepository.UpdateProductRepository(&barcodeId, &updates); err != nil {
			return err
		}
		return nil
	} else {
		return dto.ErrNoChangesRequest
	}
}

func (p *productService) DeleteProductService(barcodeId *string) error {
	_, ok := p.producRepository.RetrieveProductByBarcodeId(barcodeId)
	if !ok {
		return dto.ErrProductDoesntExist
	}
	if err := p.producRepository.DeleteProductRepository(barcodeId); err != nil {
		return err
	}
	return nil
}
