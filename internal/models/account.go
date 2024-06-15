package models

type Account struct {
	AccountID      int64  `gorm:"primaryKey;autoIncrement" json:"account_id"`
	DocumentNumber string `gorm:"type:varchar(15);not null" json:"document_number"`
}
