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

func TestUpdateDeletedProduct_Success(t *testing.T) {
	db, mock := test.MockDB(t)
	repo := repository.NewProductRepository(db)

	barcodeId := "1"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "products" SET "deleted_at"=$1,"updated_at"=$2 WHERE barcode_id = $3`)).
		WithArgs(
			nil,
			utils.AnyTime{},
			"1",
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := repo.UpdateDeletedProductRepository(&barcodeId)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateDeletedProduct_Error(t *testing.T) {
	db, mock := test.MockDB(t)
	repo := repository.NewProductRepository(db)

	barcodeId := "1"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "products" SET "deleted_at"=$1,"updated_at"=$2 WHERE barcode_id = $3`)).
		WithArgs(
			nil,
			utils.AnyTime{},
			"1",
		).
		WillReturnError(errors.New("ISE"))
	mock.ExpectRollback()

	err := repo.UpdateDeletedProductRepository(&barcodeId)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "ISE")
	assert.NoError(t, mock.ExpectationsWereMet())
}
