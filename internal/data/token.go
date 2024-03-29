package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"

	"github.com/terajari/ipdb/internal/validator"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserId    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

func generateToken(userId int64, ttl time.Duration, scope string) (*Token, error) {

	token := &Token{
		UserId: userId,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)

	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

func ValidatePlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

type TokenModel struct {
	db *sql.DB
}

type IToken interface {
	Insert(*Token) error
	DeleteAll(string, int64) error
	New(int64, time.Duration, string) (*Token, error)
}

func NewTokenModel(db *sql.DB) IToken {
	return &TokenModel{db: db}
}

func (tm *TokenModel) New(userId int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userId, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = tm.Insert(token)

	return token, err
}

func (tm *TokenModel) Insert(t *Token) error {

	query := `
		INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4)
	`
	args := []interface{}{t.Hash, t.UserId, t.Expiry, t.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := tm.db.ExecContext(ctx, query, args...)

	if err != nil {
		return err
	}

	return nil
}

func (tm *TokenModel) DeleteAll(scope string, userId int64) error {

	stmt := `
		DELETE FROM tokens
		WHERE scope = $1 AND user_id = $2
	`
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := tm.db.ExecContext(ctx, stmt, scope, userId)

	if err != nil {
		return err
	}

	return nil
}
