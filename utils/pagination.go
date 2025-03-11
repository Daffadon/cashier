package utils

import "gorm.io/gorm"

func Paginate(limit, offset uint16) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(int(limit)).Offset(int(offset))
	}
}
