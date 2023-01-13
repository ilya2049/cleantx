package app

import (
	"context"

	"cleantx/internal/domain"
	"cleantx/internal/pkg/sqldb"
)

type DoctorCommandExecutor struct {
	txBeginner sqldb.TxBeginner
	repository domain.DoctorRepository
}

func NewDoctorCommandExecutor(
	txBeginner sqldb.TxBeginner,
	repository domain.DoctorRepository,
) *DoctorCommandExecutor {
	return &DoctorCommandExecutor{
		txBeginner: txBeginner,
		repository: repository,
	}
}

func (e *DoctorCommandExecutor) TakeShift(ctx context.Context, doctorID int) error {
	doctor, err := e.repository.Get(ctx, doctorID)
	if err != nil {
		return err
	}

	doctor.TakeShift()

	if err := e.repository.Update(ctx, doctor); err != nil {
		return err
	}

	return nil
}

func (e *DoctorCommandExecutor) FinishShift(ctx context.Context, doctorID int) (err error) {
	tx, err := e.txBeginner.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() { err = sqldb.CloseTx(ctx, tx, err) }()

	doctorsOnCall, err := e.repository.WithTx(tx).ListDoctorsOnCall(ctx)
	if err != nil {
		return err
	}

	doctor := doctorsOnCall.Get(doctorID)
	if doctor == nil {
		return domain.ErrDoctorNotFound
	}

	if err := doctor.FinishShift(len(doctorsOnCall)); err != nil {
		return err
	}

	if err := e.repository.WithTx(tx).Update(ctx, doctor); err != nil {
		return err
	}

	return nil
}
