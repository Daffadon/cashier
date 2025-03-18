package repository_test

import (
	"regexp"
	"testing"
	"tiga-putra-cashier-be/dto"
	"tiga-putra-cashier-be/repository"
	test "tiga-putra-cashier-be/test/mocks/db"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSearchProduct_Success(t *testing.T) {
	db, mock := test.MockDB(t)

	repo := repository.NewProductRepository(db)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE title ILIKE $1 AND barcode_id = $2 AND "products"."deleted_at" IS NULL`)).
		WithArgs("%product-1%", "1").
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deleted_at", "barcode_id", "title", "image", "price", "description"}).
			AddRow(1, time.Now(), time.Now(), nil, "1", "product-1", "img-1", 1000, "desc-1"))

	productName := "product-1"
	barcodeId := "1"
	search := dto.SearchProductQuery{
		Title:     &productName,
		BarcodeId: &barcodeId,
	}
	products, err := repo.RetrieveProductForSearch(&search)
	assert.NoError(t, err)
	assert.Len(t, products, 1)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSearchProduct_Error(t *testing.T) {
	db, mock := test.MockDB(t)

	repo := repository.NewProductRepository(db)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE title ILIKE $1 AND "products"."deleted_at" IS NULL`)).
		WithArgs("%product-1%").
		WillReturnError(db.Error)

	productName := "product-1"
	search := dto.SearchProductQuery{
		Title: &productName,
	}
	_, err := repo.RetrieveProductForSearch(&search)
	assert.Error(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
