package dto

import (
	"errors"
	"mime/multipart"

	"github.com/shopspring/decimal"
)

var (
	ErrProductsNotFound   = errors.New("Products Not Found")
	ErrISEProducts        = errors.New("Failed to get products")
	ErrWrongFileExtension = errors.New("File should be jpeg,jpg, or png")
	ErrLimitSizeExceeded  = errors.New("File should be no more than 6MB")
	ErrToAddProduct       = errors.New("Failed to Add Product")
	ErrProductExist       = errors.New("Product with this barcode already Exist")
	ErrProductDoesntExist = errors.New("Product with this barcode doesn't exist")
	ErrNoChangesRequest   = errors.New("There is no updated field in the request")

	MESSAGE_SUCCESS_GET_ALL_PRODUCTS   = "Success Get All product"
	MESSAGE_SUCCESS_GET_PRODUCT_DETAIL = "Success Get Product Detail"
	MESSAGE_SUCCESS_SEARCH_PRODUCTS    = "Success Get All product"
	MESSAGE_SUCCESS_ADD_PRODUCT        = "Success Add Product"
	MESSAGE_SUCCESS_UPDATE_PRODUCT     = "Success Update Product"
	MESSAGE_SUCCESS_DELETE_PRODUCT     = "Success Delete Product"
)

type (
	ProductWithoutTimeStamp struct {
		BarcodeId   string          `json:"barcode_id" binding:"required"`
		Image       string          `json:"image" binding:"required"`
		Title       string          `json:"title" binding:"required"`
		Price       decimal.Decimal `json:"price" binding:"required"`
		Description string          `json:"description" binding:"required"`
	}

	AllProductsWithPagination struct {
		Products     []ProductWithoutTimeStamp `json:"products"`
		PageMetaData PaginationResponse        `json:"page_meta_data"`
	}

	AddProductRequest struct {
		BarcodeId   string                `form:"barcode_id" binding:"required"`
		Image       *multipart.FileHeader `form:"image" binding:"required"`
		Title       string                `form:"title" binding:"required"`
		Price       decimal.Decimal       `form:"price" binding:"required"`
		Description string                `form:"description" binding:"required"`
	}

	ProductBarcodeIdURI struct {
		BarcodeId string `uri:"barcode_id" binding:"required"`
	}

	UpdateProductRequest struct {
		Image       *multipart.FileHeader `form:"image"`
		Title       *string               `form:"title"`
		Price       *decimal.Decimal      `form:"price"`
		Description *string               `form:"description"`
	}

	SearchProductQuery struct {
		Title     *string `form:"title"`
		BarcodeId *string `form:"barcode_id"`
	}
)
