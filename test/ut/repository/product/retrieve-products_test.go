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

func TestRetrieveProducts_Sucess(t *testing.T) {
	db, mock := test.MockDB(t)

	repo := repository.NewProductRepository(db)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE "products"."deleted_at" IS NULL LIMIT $1`)).
		WithArgs(12).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deleted_at", "barcode_id", "title", "image", "price", "description"}).
			AddRow(1, time.Now(), time.Now(), nil, "1", "Product A", "img-1", 1000, "desc-1").
			AddRow(2, time.Now(), time.Now(), nil, "2", "Product B", "img-2", 2000, "desc-2"))

	products, err := repo.RetrieveProductsRepository(12, 0)
	assert.NoError(t, err)
	assert.Len(t, products, 2)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestRetrieveProducts_Error(t *testing.T) {
	db, mock := test.MockDB(t)

	repo := repository.NewProductRepository(db)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "products" WHERE "products"."deleted_at" IS NULL LIMIT $1`)).
		WithArgs(12).
		WillReturnError(db.Error)

	products, err := repo.RetrieveProductsRepository(12, 0)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), dto.ErrISEProducts.Error())
	assert.Len(t, products, 0)

	assert.NoError(t, mock.ExpectationsWereMet())
}
