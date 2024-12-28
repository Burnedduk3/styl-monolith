package domain

import "time"

type Role struct {
	ID          uint
	Name        string
	Description string
	Created     time.Time
	Updated     time.Time
	Deleted     time.Time
}

type Permissions struct {
	ID      uint
	Name    string
	Created time.Time
	Updated time.Time
	Deleted time.Time
}
