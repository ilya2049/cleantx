package app

import (
	"context"

	"cleantx/internal/domain"
)

type DoctorCommandExecutor struct {
	unitOfWorkProvider domain.UnitOfWorkProvider
	repository         domain.DoctorRepository
}

func NewDoctorCommandExecutor(
	unitOfWorkProvider domain.UnitOfWorkProvider,
	repository domain.DoctorRepository,
) *DoctorCommandExecutor {
	return &DoctorCommandExecutor{
		unitOfWorkProvider: unitOfWorkProvider,
		repository:         repository,
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
	return e.unitOfWorkProvider.Atomic(ctx, func(ctx context.Context, uow domain.UnitOfWork) error {
		doctorRepository := uow.NewDoctorRepository()

		doctorsOnCall, err := doctorRepository.ListDoctorsOnCall(ctx)
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
	})
}
