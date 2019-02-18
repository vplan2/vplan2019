package auth

import (
	"github.com/zekroTJA/vplan2019/internal/database"
)

// TokenManager is for managing, creating and
// checking user API tokens
type TokenManager struct {
	db *database.Driver
}

// NewTokenManager creates a new instance of TokenManager
// with the specified database driver as storage module
func NewTokenManager(db *database.Driver) *TokenManager {
	return &TokenManager{
		db: db,
	}
}

// TODO:
//   - add check token function
//   - add set token function
//   - add remove token function
