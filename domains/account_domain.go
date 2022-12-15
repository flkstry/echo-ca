package domains

import (
	"context"
	"time"
)

type Account struct {
	ID        string    `json:"id"`
	Email     string    `json:"email" validate:"required,email"`
	Password  string    `json:"password" query:"hashed_password" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// implementation for Account
type AccountRepository interface {
	Fetch(ctx context.Context, q string, args ...interface{}) ([]Account, error)
	Store(ctx context.Context, a *Account) error
	Update(ctx context.Context, a *Account) error
	Delete(ctx context.Context, id int64) error

	GetById(ctx context.Context, id int64) (Account, error)
	GetByEmail(ctx context.Context, email string) (Account, error)
}

// usecase for Account implementation
type AccountUsecase interface {
	Store(ctx context.Context, a *Account) error
	GetByEmail(ctx context.Context, email string) (Account, error)
}
