package postgres

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Name      string
	Email     string
	Country   string
	Role      string
	Status    string
	LastIp    string
	LastLogin time.Time
}
