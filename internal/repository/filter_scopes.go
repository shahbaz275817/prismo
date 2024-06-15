package repository

import (
	"time"

	"gorm.io/gorm"
)

type FilterRequest struct {
	SortBy      string
	Limit       int
	Offset      int
	CreatedFrom *time.Time
	CreatedTo   *time.Time
	StartedFrom *time.Time
	StartedTo   *time.Time
	RawQuery    string
}

func (fr FilterRequest) CreatedAtRange() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if fr.CreatedFrom != nil && fr.CreatedTo != nil {
			return db.Where("created_at >= (?) and created_at <= (?)", fr.CreatedFrom, fr.CreatedTo)
		}
		return db
	}
}

func (fr FilterRequest) Pagination() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if fr.Limit > 0 {
			return db.Limit(fr.Limit).Offset(fr.Offset)
		}
		return db
	}
}

func (fr FilterRequest) Sort() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if fr.SortBy != "" {
			return db.Order(fr.SortBy)
		}
		return db
	}
}
