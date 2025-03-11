package entity

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	BarcodeId   string `gorm:"unique"`
	Image       string
	Title       string
	Price       decimal.Decimal
	Description string
}
