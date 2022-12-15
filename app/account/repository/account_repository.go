package repository

import (
	"context"
	"database/sql"
	"login-system/server/domains"

	"github.com/sirupsen/logrus"
)

type postgreAccountRepository struct {
	DB *sql.DB
}

func NewPsqlAccountRepository(Conn *sql.DB) domains.AccountRepository {
	return &postgreAccountRepository{Conn}
}

func (psql *postgreAccountRepository) Fetch(ctx context.Context, query string, args ...interface{}) (psqles []domains.Account, err error) {
	rows, err := psql.DB.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	res := make([]domains.Account, 0)
	for rows.Next() {
		t := domains.Account{}
		err = rows.Scan(
			&t.ID,
			&t.Email,
			&t.Password,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		res = append(res, t)
	}

	return res, nil
}

func (psql *postgreAccountRepository) Store(ctx context.Context, a *domains.Account) (err error) {
	q := `INSERT INTO account (
		id,
		email,
		hashed_password,
		created_at,
		updated_at
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5
	) RETURNING id;`

	stm, err := psql.DB.PrepareContext(ctx, q)
	if err != nil {
		return
	}

	lastInsertedId := ""

	err = stm.QueryRowContext(ctx, a.ID, a.Email, a.Password, a.CreatedAt, a.UpdatedAt).Scan(&lastInsertedId)
	if err != nil {
		return
	}

	a.ID = lastInsertedId
	return
}

func (psql *postgreAccountRepository) Update(ctx context.Context, a *domains.Account) (err error) {

	return
}

func (psql *postgreAccountRepository) Delete(ctx context.Context, id int64) (err error) {
	return
}
func (psql *postgreAccountRepository) GetById(ctx context.Context, id int64) (res domains.Account, err error) {
	return
}

func (psql *postgreAccountRepository) GetByEmail(ctx context.Context, email string) (res domains.Account, err error) {
	q := `SELECT * FROM account WHERE email=$1;`

	list, err := psql.Fetch(ctx, q, email)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domains.ErrNotFound
	}

	return
}
