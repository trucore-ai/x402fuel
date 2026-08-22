package types

import (
	"crypto/ecdsa"
	"time"
)

type Wallet struct {
	ID        string    `json:"id"`
	Address   string    `json:"address"`
	KeyPath   string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	Label     string    `json:"label"`
}

type SigningKey struct {
	PrivateKey *ecdsa.PrivateKey
	Address    string
}