package domain

import "context"

type UnitOfWorkProvider interface {
	Atomic(
		ctx context.Context,
		doAtomically func(context.Context, UnitOfWork) error,
	) error
}

type UnitOfWork interface {
	NewDoctorRepository() DoctorRepository
}
