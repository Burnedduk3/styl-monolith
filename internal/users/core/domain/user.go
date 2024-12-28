package domain

import "time"

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusDeleted  Status = "deleted"
)

type User struct {
	ID          uint
	Username    string
	Name        string
	Email       string
	Country     string
	Phone       string
	CountryCode string
	Role        Role
	Status      Status
	LastIp      string
	LastLogin   time.Time
	Created     time.Time
	Updated     time.Time
	Deleted     time.Time
}
