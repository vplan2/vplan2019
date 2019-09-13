package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
	"time"

	"github.com/vplan2/vplan2019/internal/database"
)

// TokenManager is for managing, creating and
// checking user API tokens
type TokenManager struct {
	db       database.Driver
	lifetime time.Duration
}

// NewTokenManager creates a new instance of TokenManager
// with the specified database driver as storage module
func NewTokenManager(db database.Driver, tokenLifetime time.Duration) *TokenManager {
	return &TokenManager{
		db:       db,
		lifetime: tokenLifetime,
	}
}

// Check checks if the passed token is valid. If the token
// is invalid, the user has no token registered or the
// token exeeded, this returns false.
// Only if the db access failes, it returns an error.
func (t *TokenManager) Check(token string) (string, error) {
	token = strSHA256SumToken(token)

	ident, expire, err := t.db.GetAPIToken(token)
	if err != nil {
		return "", err
	}

	if ident == "" {
		return "", nil
	}

	if time.Now().After(expire) {
		err = t.Delete(ident)
		return "", err
	}

	return ident, nil
}

// Set creates a new random token with a random number plus a current
// unix timestamp encoded witn base64 to a string. This token
// will be returned with the time of expiration. The token will be
// sha256-hashed and then saved in the database.
func (t *TokenManager) Set(ident string) (string, time.Time, error) {
	rInt, err := rand.Int(rand.Reader, big.NewInt(9999999999))
	if err != nil {
		return "", time.Time{}, err
	}

	tokenBlob := []byte(fmt.Sprintf("%d%d", rInt, time.Now().UnixNano()))
	token := base64.StdEncoding.EncodeToString(tokenBlob)
	hashedToken := strSHA256SumToken(token)

	expire := time.Now().Add(t.lifetime)
	err = t.db.SetUserAPIToken(ident, hashedToken, expire)
	if err != nil {
		return "", time.Time{}, err
	}

	return token, expire, nil
}

// Delete removes a token entry from the database of
// the passed user ident.
// This only returns an error, if the database access
// failes and not if there was no token registered before.
func (t *TokenManager) Delete(ident string) error {
	return t.db.DeleteUserAPIToken(ident)
}

// strSHA256SumToken creates a sha256-hashed
// hexadecimal string from the source string
func strSHA256SumToken(token string) string {
	bSum := sha256.Sum256([]byte(token))
	return fmt.Sprintf("%x", bSum)
}
