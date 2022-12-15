package usecase

import (
	"context"
	"login-system/server/domains"
	"login-system/server/internals"
	"time"

	"github.com/google/uuid"
)

type AccountUsecase struct {
	accountRepository domains.AccountRepository
	contextTimeout    time.Duration
	hashingParams     internals.Params
}

func NewAccountUsecase(a domains.AccountRepository, to time.Duration, p internals.Params) domains.AccountUsecase {
	return &AccountUsecase{
		accountRepository: a,
		contextTimeout:    to,
		hashingParams:     p,
	}
}

func (a *AccountUsecase) Store(c context.Context, acc *domains.Account) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	// check existed account with email
	existed, _ := a.accountRepository.GetByEmail(ctx, acc.Email)
	if (existed != domains.Account{}) {
		return domains.ErrConflict
	}

	// generate uuid
	id := uuid.Must(uuid.NewRandom()).String()
	acc.ID = id

	// set current time
	timeNow := time.Now()
	acc.CreatedAt = timeNow
	acc.UpdatedAt = timeNow

	// generate hashing password
	hashed_password, err := internals.GenerateHashedPassword(acc.Password, &a.hashingParams)
	if err != nil {
		return
	}

	acc.Password = string(hashed_password[:])

	err = a.accountRepository.Store(ctx, acc)
	return
}

func (a *AccountUsecase) GetByEmail(c context.Context, email string) (res domains.Account, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.accountRepository.GetByEmail(ctx, email)
	if err != nil {
		return
	}

	return
}
