package password

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	Hash          string
	PlainTextPass string
}

type AuthToken struct {
	Plaintext string    `json:"token"`
	Hash      string    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

func (p *Password) HashPassword(plainTextPass string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPass), 12)

	if err != nil {
		return err
	}

	p.PlainTextPass = plainTextPass
	p.Hash = string(hash)

	return nil
}

func (p *Password) Matches(plainTextPass string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(p.Hash), []byte(plainTextPass))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func (p *Password) GenerateAuthToken(userID int64, ttl time.Duration) (*AuthToken, error) {
	token := &AuthToken{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
	}

	randomBytes := make([]byte, 16)

	// fill byte slice with random bytes from OS
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = string(hash[:])

	return token, nil
}
