package pgsql

import (
	"cleantx/internal/domain"

	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

type UnitOfWorkProvider struct {
	db *pgx.Conn
}

func NewUnitOfWorkProvider(db *pgx.Conn) *UnitOfWorkProvider {
	return &UnitOfWorkProvider{
		db: db,
	}
}

func (p *UnitOfWorkProvider) NewUnitOfWork(ctx context.Context) (domain.UnitOfWork, error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &UnitOfWork{
		tx: tx,
	}, nil
}

func (s *UnitOfWorkProvider) CommitChanges(ctx context.Context, unitOfWorkContainer any, err error) error {
	unitOfWork := unitOfWorkContainer.(*UnitOfWork)

	if err != nil {
		if rollbackErr := unitOfWork.tx.Rollback(ctx); rollbackErr != nil {
			log.Println("failed to finish shift: failed to rollback tx", err.Error())
		}

		return err
	}

	return unitOfWork.tx.Commit(ctx)
}

type UnitOfWork struct {
	tx pgx.Tx
}

func (s *UnitOfWork) NewDoctorRepository() domain.DoctorRepository {
	return NewDoctorRepository(s.tx)
}
