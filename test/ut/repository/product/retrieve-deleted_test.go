package repository_test

import (
	"errors"
	"regexp"
	"testing"
	"tiga-putra-cashier-be/repository"
	test "tiga-putra-cashier-be/test/mocks/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveDeletedProduct_Success(t *testing.T) {
	db, mock := test.MockDB(t)

	repo := repository.NewProductRepository(db)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "products"."barcode_id","products"."image","products"."title","products"."price","products"."description" FROM "products" WHERE barcode_id = $1 AND deleted_at IS NOT NULL ORDER BY "products"."id" LIMIT $2`)).
		WithArgs("1", 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"barcode_id", "title", "image", "price", "description"}).
			AddRow("1", "Product A", "img-1", 1000, "desc-1"))

	barcodeId := "1"
	products, ok := repo.RetrieveDeletedProductByBarcodeId(&barcodeId)
	assert.True(t, true, ok)
	assert.Equal(t, products.BarcodeId, "1")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetrieveDeletedProduct_Error(t *testing.T) {
	db, mock := test.MockDB(t)

	repo := repository.NewProductRepository(db)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT "products"."barcode_id","products"."image","products"."title","products"."price","products"."description" FROM "products" WHERE barcode_id = $1 AND deleted_at IS NOT NULL ORDER BY "products"."id" LIMIT $2`)).
		WithArgs("1", 1).
		WillReturnError(errors.New("ISE"))

	barcodeId := "1"
	_, ok := repo.RetrieveDeletedProductByBarcodeId(&barcodeId)
	assert.False(t, false, ok)
	assert.NoError(t, mock.ExpectationsWereMet())
}
