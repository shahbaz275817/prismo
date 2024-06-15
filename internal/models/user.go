package models

import "time"

type User struct {
	ID        *int64
	HubID     int64
	Name      string
	Email     string
	IsActive  bool
	CreatedAt time.Time
	CreatedBy *int64
	UpdatedAt time.Time
	UpdatedBy *int64
}
