package domain

import "time"

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusDeleted  Status = "deleted"
)

type User struct {
	ID        int
	Name      string
	Email     string
	Country   string
	Role      string
	Status    Status
	LastIp    string
	LastLogin time.Time
	Created   time.Time
	Updated   time.Time
	Deleted   time.Time
}
