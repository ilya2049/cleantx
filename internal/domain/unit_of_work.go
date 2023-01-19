package domain

import "context"

type UnitOfWorkProvider interface {
	NewUnitOfWork(context.Context) (UnitOfWork, error)
	CommitChanges(ctx context.Context, unitOfWork any, err error) error
}

type UnitOfWork interface {
	NewDoctorRepository() DoctorRepository
}
