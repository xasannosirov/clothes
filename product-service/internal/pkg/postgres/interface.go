package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type Tx interface {
	pgx.Tx
}

type RepoTx interface {
	Beging(ctx context.Context) (Tx, error)
	TxRollBack(ctx context.Context, tx Tx, err error)
}
