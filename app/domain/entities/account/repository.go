package account

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

type Repository interface {
	Create(ctx context.Context, account Account) (Account, error)
	GetBalance(ctx context.Context, id AccountId) (int, error)
	Fetch(ctx context.Context) ([]Account, error)
	GetById(ctx context.Context, id AccountId) (Account, error)
	GetByCpf(ctx context.Context, cpf cpf.Cpf) (Account, error)
	CreditAccount(ctx context.Context, id AccountId, amount money.Money) error
	DebitAccount(ctx context.Context, id AccountId, amount money.Money) error
}
