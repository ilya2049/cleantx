package app

import (
	"context"

	"cleantx/internal/domain"
)

type DoctorCommandExecutor interface {
	TakeShift(ctx context.Context, doctorID int) error
	FinishShift(ctx context.Context, doctorID int) error
}

type doctorCommandExecutor struct {
	repository domain.DoctorRepository
}

func NewDoctorCommandExecutor(repository domain.DoctorRepository) DoctorCommandExecutor {
	return &doctorCommandExecutor{
		repository: repository,
	}
}

func (e *doctorCommandExecutor) TakeShift(ctx context.Context, doctorID int) error {
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

func (e *doctorCommandExecutor) FinishShift(ctx context.Context, doctorID int) error {
	doctorsOnCall, err := e.repository.ListDoctorsOnCall(ctx)
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

	if err := e.repository.Update(ctx, doctor); err != nil {
		return err
	}

	return nil
}
