package models

import (
	"time"
)

type Transaction struct {
	TransactionID   int64     `gorm:"primaryKey;autoIncrement" json:"transaction_id"`
	AccountID       int64     `gorm:"not null" json:"account_id"`
	OperationTypeID int64     `gorm:"column:operationtype_id;not null" json:"operationtype_id"`
	Amount          float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	EventDate       time.Time `gorm:"column:eventdate;type:timestamp;not null" json:"event_date"`

	Account        Account        `gorm:"foreignKey:AccountID;references:AccountID"`
	OperationsType OperationsType `gorm:"foreignKey:OperationTypeID;references:OperationTypeID"`
}
