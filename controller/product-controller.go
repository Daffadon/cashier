package controller

import (
	"net/http"
	"tiga-putra-cashier-be/constant"
	"tiga-putra-cashier-be/dto"
	"tiga-putra-cashier-be/service"
	"tiga-putra-cashier-be/utils"

	"github.com/gin-gonic/gin"
)

type (
	ProductController interface {
		GetProduct(ctx *gin.Context)
		GetProductDetail(ctx *gin.Context)
		SearchProduct(ctx *gin.Context)
		AddProduct(ctx *gin.Context)
		UpdateProduct(ctx *gin.Context)
		DeleteProduct(ctx *gin.Context)
	}
	productController struct {
		productService service.ProductService
	}
)

func NewProductController(productService service.ProductService) ProductController {
	return &productController{productService}
}

func (p *productController) GetProduct(ctx *gin.Context) {
	var req dto.PaginationRequest
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		res := utils.ReturnResponseError(400, dto.ErrBadrequest.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	products, err := p.productService.GetProductService(&req.Page)
	if err == dto.ErrProductsNotFound {
		res := utils.ReturnResponseError(404, err.Error())
		ctx.AbortWithStatusJSON(http.StatusNotFound, res)
		return
	} else if err == dto.ErrISEProducts {
		res := utils.ReturnResponseError(500, err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}
	res := utils.ReturnResponseSuccess(200, dto.MESSAGE_SUCCESS_GET_ALL_PRODUCTS, products)
	ctx.JSON(http.StatusOK, res)
}

func (p *productController) GetProductDetail(ctx *gin.Context) {
	var req dto.ProductBarcodeIdURI
	if err := ctx.ShouldBindUri(&req); err != nil {
		res := utils.ReturnResponseError(400, dto.ErrBadrequest.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	product, err := p.productService.GetProductDetailService(&req.BarcodeId)
	if err != nil {
		if err == dto.ErrProductDoesntExist {
			res := utils.ReturnResponseError(404, err.Error())
			ctx.AbortWithStatusJSON(http.StatusNotFound, res)
			return
		}
	}
	res := utils.ReturnResponseSuccess(200, dto.MESSAGE_SUCCESS_GET_PRODUCT_DETAIL, product)
	ctx.JSON(http.StatusOK, res)
}

func (p *productController) SearchProduct(ctx *gin.Context) {
	var req dto.SearchProductQuery
	ctx.ShouldBindQuery(&req)
	if req.BarcodeId == nil && req.Title == nil {
		res := utils.ReturnResponseError(400, dto.ErrBadrequest.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	products, err := p.productService.SearchProductService(&req)
	if err != nil {
		if err == dto.ErrProductsNotFound {
			res := utils.ReturnResponseError(404, err.Error())
			ctx.AbortWithStatusJSON(http.StatusNotFound, res)
			return
		}
		res := utils.ReturnResponseError(500, err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}
	res := utils.ReturnResponseSuccess(200, dto.MESSAGE_SUCCESS_SEARCH_PRODUCTS, products)
	ctx.JSON(http.StatusOK, res)
}

func (p *productController) AddProduct(ctx *gin.Context) {
	var req dto.AddProductRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := utils.ReturnResponseError(400, dto.ErrBadrequest.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	if req.Image.Size > constant.MaxUploadSize {
		res := utils.ReturnResponseError(400, dto.ErrLimitSizeExceeded.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	ext := utils.GetFileNameExtension(req.Image.Filename)
	if ext != "jpg" && ext != "jpeg" && ext != "png" {
		res := utils.ReturnResponseError(400, dto.ErrWrongFileExtension.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	if err := p.productService.CreateProductService(req); err != nil {
		if err == dto.ErrProductExist {
			res := utils.ReturnResponseError(409, err.Error())
			ctx.AbortWithStatusJSON(http.StatusConflict, res)
			return
		}
		res := utils.ReturnResponseError(500, err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}
	res := utils.ReturnResponseSuccess(200, dto.MESSAGE_SUCCESS_ADD_PRODUCT)
	ctx.JSON(http.StatusOK, res)
}

func (p *productController) UpdateProduct(ctx *gin.Context) {
	var barcodeId dto.ProductBarcodeIdURI
	if err := ctx.ShouldBindUri(&barcodeId); err != nil {
		res := utils.ReturnResponseError(400, dto.ErrBadrequest.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	var req dto.UpdateProductRequest
	ctx.ShouldBind(&req)
	err := p.productService.UpdateProductService(barcodeId.BarcodeId, req)
	if err != nil {
		if err == dto.ErrProductDoesntExist {
			res := utils.ReturnResponseError(404, err.Error())
			ctx.AbortWithStatusJSON(http.StatusNotFound, res)
			return
		}
		if err == dto.ErrNoChangesRequest {
			res := utils.ReturnResponseError(304, err.Error())
			ctx.AbortWithStatusJSON(http.StatusNotModified, res)
			return
		}
		res := utils.ReturnResponseError(500, err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}
	res := utils.ReturnResponseSuccess(200, dto.MESSAGE_SUCCESS_UPDATE_PRODUCT)
	ctx.JSON(http.StatusOK, res)
}

func (p *productController) DeleteProduct(ctx *gin.Context) {
	var req dto.ProductBarcodeIdURI
	if err := ctx.ShouldBindUri(&req); err != nil {
		res := utils.ReturnResponseError(400, dto.ErrBadrequest.Error())
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}
	if err := p.productService.DeleteProductService(&req.BarcodeId); err != nil {
		if err == dto.ErrProductDoesntExist {
			res := utils.ReturnResponseError(404, err.Error())
			ctx.AbortWithStatusJSON(http.StatusNotFound, res)
			return
		}
		res := utils.ReturnResponseError(500, err.Error())
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, res)
		return
	}
	res := utils.ReturnResponseSuccess(200, dto.MESSAGE_SUCCESS_DELETE_PRODUCT)
	ctx.JSON(http.StatusOK, res)
}
