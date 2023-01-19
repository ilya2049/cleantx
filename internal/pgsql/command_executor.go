package pgsql

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type commandExecutor interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
}

type commandExecutorKey string

const ctxCommandExecutorKey commandExecutorKey = "tx"

func beginTxAndInjectInCtx(ctx context.Context, db *pgx.Conn,
) (context.Context, pgx.Tx, error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return ctx, nil, err
	}

	return context.WithValue(ctx, ctxCommandExecutorKey, tx), tx, nil
}

func closeTx(ctx context.Context, tx pgx.Tx, err error) error {
	if err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			log.Println("failed to finish shift: failed to rollback tx", err.Error())
		}

		return err
	}

	return tx.Commit(ctx)
}

func getCommandExecutorFromCtxOrDefault(
	ctx context.Context,
	defaultCommandExecutor commandExecutor,
) commandExecutor {
	ctxCommandExecutorContainer := ctx.Value(ctxCommandExecutorKey)
	if ctxCommandExecutorContainer == nil {
		return defaultCommandExecutor
	}

	ctxCommandExecutor, ok := ctxCommandExecutorContainer.(commandExecutor)
	if !ok {
		return defaultCommandExecutor
	}

	return ctxCommandExecutor
}
