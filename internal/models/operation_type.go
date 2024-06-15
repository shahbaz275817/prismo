package models

type OperationsType struct {
	OperationTypeID int64  `gorm:"column:operationtype_id;primaryKey;autoIncrement" json:"operationtype_id"`
	Description     string `gorm:"type:varchar(50);not null" json:"description"`
}

func (OperationsType) TableName() string {
	return "operationstypes"
}
