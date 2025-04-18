package models

import (
	"time"
)

type User struct {
	ID           uint
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type LoyaltyAccount struct {
	ID        uint
	UserID    uint
	User      User
	Balance   int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Transaction struct {
	ID        uint
	AccountID uint
	Account   LoyaltyAccount
	Type      string
	Points    int
	Timestamp time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
