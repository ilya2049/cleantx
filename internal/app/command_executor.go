package app

import (
	"context"

	"cleantx/internal/domain"
)

type DoctorCommandExecutor struct {
	repository domain.DoctorRepository
}

func NewDoctorCommandExecutor(repository domain.DoctorRepository) *DoctorCommandExecutor {
	return &DoctorCommandExecutor{
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

func (e *DoctorCommandExecutor) FinishShift(ctx context.Context, doctorID int) error {
	return e.repository.FinishShift(ctx, doctorID)
}
