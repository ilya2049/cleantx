package pgsql

import (
	"cleantx/internal/app"

	"context"

	"github.com/jackc/pgx/v5"
)

type TxDoctorCommandExecutor struct {
	app.DoctorCommandExecutor

	db *pgx.Conn
}

func DecorateDoctorCommandExecutorWithTx(
	db *pgx.Conn,
	appDoctorCommandExecutor app.DoctorCommandExecutor,
) app.DoctorCommandExecutor {
	return &TxDoctorCommandExecutor{
		DoctorCommandExecutor: appDoctorCommandExecutor,
		db:                    db,
	}
}

func (e *TxDoctorCommandExecutor) TakeShift(ctx context.Context, doctorID int) error {
	ctx, tx, err := beginTxAndInjectInCtx(ctx, e.db)
	if err != nil {
		return err
	}

	err = e.DoctorCommandExecutor.TakeShift(ctx, doctorID)

	return closeTx(ctx, tx, err)
}

func (e *TxDoctorCommandExecutor) FinishShift(ctx context.Context, doctorID int) error {
	ctx, tx, err := beginTxAndInjectInCtx(ctx, e.db)
	if err != nil {
		return err
	}

	err = e.DoctorCommandExecutor.FinishShift(ctx, doctorID)

	return closeTx(ctx, tx, err)
}
