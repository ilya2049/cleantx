package pgsql

import (
	"cleantx/internal/domain"

	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

type UnitOfWorkProvider struct {
	db commandExecutor
}

func NewUnitOfWorkProvider(db commandExecutor) *UnitOfWorkProvider {
	return &UnitOfWorkProvider{
		db: db,
	}
}

func (p *UnitOfWorkProvider) Atomic(
	ctx context.Context,
	doAtomically func(context.Context, domain.UnitOfWork) error,
) (err error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				log.Println("failed to finish shift: failed to rollback tx", err.Error())
			}

			return
		}

		err = tx.Commit(ctx)
	}()

	return doAtomically(ctx, &UnitOfWork{tx: tx})
}

type UnitOfWork struct {
	tx pgx.Tx
}

func (s *UnitOfWork) NewDoctorRepository() domain.DoctorRepository {
	return NewDoctorRepository(s.tx)
}
