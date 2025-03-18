package repository_test

import (
	"errors"
	"regexp"
	"testing"
	"tiga-putra-cashier-be/dto"
	"tiga-putra-cashier-be/repository"
	test "tiga-putra-cashier-be/test/mocks/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCountProduct_Success(t *testing.T) {
	db, mock := test.MockDB(t)

	repo := repository.NewProductRepository(db)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "products" WHERE "products"."deleted_at" IS NULL`)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))
	count, err := repo.CountProductsRepository()
	assert.NoError(t, err)
	assert.Equal(t, uint16(count), uint16(10))

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCountProduct_Error(t *testing.T) {
	db, mock := test.MockDB(t)

	repo := repository.NewProductRepository(db)
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "products" WHERE "products"."deleted_at" IS NULL`)).
		WillReturnError(errors.New("ISE"))
	_, err := repo.CountProductsRepository()
	assert.Error(t, err)
	assert.Equal(t, err.Error(), dto.ErrISEProducts.Error())

	assert.NoError(t, mock.ExpectationsWereMet())
}
