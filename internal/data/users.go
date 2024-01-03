package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/terajari/ipdb/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var ErrDuplicateEmail = errors.New("duplicate email")

type User struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email" binding:"required,email"`
	Password  password  `json:"-" binding:"required,min=8,max=72"`
	Activated bool      `json:"active"`
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}

type UserModel struct {
	Db *sql.DB
}

type IUser interface {
	Insert(user *User) error
	GetByEmail(email string) (*User, error)
	Update(user *User) error
}

func NewUserModel(db *sql.DB) IUser {
	return &UserModel{Db: db}
}

type password struct {
	Plaintext *string
	Hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.Plaintext = &plaintextPassword
	p.Hash = hash
	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextPassword))
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

func ValidateUser(v *validator.Validator, u *User) {
	v.Check(u.Name != "", "name", "must be provided")
	v.Check(len(u.Name) <= 500, "name", "must not be more than 500 bytes long")
	v.Check(u.Email != "", "email", "must be provided")
	v.Check(u.Password.Plaintext != nil, "password", "must be provided")
	v.Check(len(*u.Password.Plaintext) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(*u.Password.Plaintext) <= 72, "password", "must not be more than 72 bytes long")

	if u.Password.Hash == nil {
		panic("missing password hash for user")
	}
}

func (m UserModel) Insert(user *User) error {

	query := `
		INSERT INTO users (name, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, version
	`

	args := []any{user.Name, user.Email, user.Password.Hash}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.Db.QueryRowContext(ctx, query, args...).Scan(&user.Id, &user.CreatedAt, &user.Version)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

// GetByEmail implements IUser.
func (m *UserModel) GetByEmail(email string) (*User, error) {

	query := `
		SELECT id, created_at, name, email, password_hash, activated, version
		FROM users
		WHERE email = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User

	err := m.Db.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			ErrRecordNotFound := errors.New("record not found")
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

// Update implements IUser.
func (m *UserModel) Update(user *User) error {

	query := `
		UPDATE users SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		user.Name,
		user.Email,
		user.Password.Hash,
		user.Activated,
		user.Id,
		user.Version,
	}

	err := m.Db.QueryRowContext(ctx, query, args...).Scan(&user.Version)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			ErrDuplicateEmail := errors.New("duplicate email")
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			ErrEditConflict := errors.New("edit conflict")
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}
