package repository_test

import (
	"errors"
	"regexp"
	"testing"
	"tiga-putra-cashier-be/dto"
	"tiga-putra-cashier-be/entity"
	"tiga-putra-cashier-be/repository"
	test "tiga-putra-cashier-be/test/mocks/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestCreateProduct_Success(t *testing.T) {
	db, mock := test.MockDB(t)

	repo := repository.NewProductRepository(db)
	prod := &entity.Product{
		BarcodeId:   "1",
		Title:       "title-1",
		Image:       "img-1",
		Price:       decimal.NewFromInt32(1000),
		Description: "desc1",
	}
	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "products" ("created_at","updated_at","deleted_at","barcode_id","image","title","price","description")
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			prod.BarcodeId,
			prod.Image,
			prod.Title,
			prod.Price,
			prod.Description,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectCommit()
	err := repo.CreateProductRepository(prod)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateProduct_Error(t *testing.T) {
	db, mock := test.MockDB(t)

	repo := repository.NewProductRepository(db)
	prod := &entity.Product{
		BarcodeId:   "1",
		Title:       "title-1",
		Image:       "img-1",
		Price:       decimal.NewFromInt32(1000),
		Description: "desc1",
	}
	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(
		`INSERT INTO "products" ("created_at","updated_at","deleted_at","barcode_id","image","title","price","description")
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING "id"`)).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			prod.BarcodeId,
			prod.Image,
			prod.Title,
			prod.Price,
			prod.Description,
		).
		WillReturnError(errors.New("error"))
	mock.ExpectRollback()
	err := repo.CreateProductRepository(prod)
	assert.Error(t, err)
	assert.Equal(t, err, dto.ErrToAddProduct)
	assert.NoError(t, mock.ExpectationsWereMet())
}
