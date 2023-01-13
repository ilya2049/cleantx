package sqldb

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

type TxBeginner interface {
	Begin(context.Context) (Tx, error)
}

type Tx any

type txBeginner struct {
	db *pgx.Conn
}

func NewTxBeginner(db *pgx.Conn) TxBeginner {
	return &txBeginner{
		db: db,
	}
}

func (b *txBeginner) Begin(ctx context.Context) (Tx, error) {
	return b.db.Begin(ctx)
}

func CloseTx(ctx context.Context, appTx Tx, err error) error {
	tx := appTx.(pgx.Tx)

	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			log.Println("failed to finish shift: failed to rollback tx", err.Error())
		}

		return err
	}

	return tx.Commit(ctx)
}
