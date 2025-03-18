package repository_test

import (
	"errors"
	"regexp"
	"testing"
	"tiga-putra-cashier-be/repository"
	test "tiga-putra-cashier-be/test/mocks/db"
	"tiga-putra-cashier-be/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateProduct_Success(t *testing.T) {
	db, mock := test.MockDB(t)
	repo := repository.NewProductRepository(db)

	barcodeId := "1"
	productUpdates := map[string]any{
		"title":       "Updated Title",
		"image":       "updated-img.png",
		"price":       2000,
		"description": "Updated description",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "products" SET "description"=$1,"image"=$2,"price"=$3,"title"=$4,"updated_at"=$5 WHERE barcode_id = $6 AND "products"."deleted_at" IS NULL`)).
		WithArgs(
			"Updated description",
			"updated-img.png",
			2000,
			"Updated Title",
			utils.AnyTime{},
			"1",
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.UpdateProductRepository(&barcodeId, &productUpdates)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateProduct_Error(t *testing.T) {
	db, mock := test.MockDB(t)
	repo := repository.NewProductRepository(db)

	barcodeId := "1"
	productUpdates := map[string]any{
		"title":       "Updated Title",
		"image":       "updated-img.png",
		"price":       2000,
		"description": "Updated description",
	}

	mock.ExpectBegin()
	mock.MatchExpectationsInOrder(false)
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "products" SET "description"=$1,"image"=$2,"price"=$3,"title"=$4,"updated_at"=$5 WHERE barcode_id = $6 AND "products"."deleted_at" IS NULL`)).
		WithArgs(
			"Updated description",
			"updated-img.png",
			2000,
			"Updated Title",
			utils.AnyTime{},
			"1",
		).
		WillReturnError(errors.New("ISE"))
	mock.ExpectRollback()

	err := repo.UpdateProductRepository(&barcodeId, &productUpdates)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "ISE")
	assert.NoError(t, mock.ExpectationsWereMet())
}
