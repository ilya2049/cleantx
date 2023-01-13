package pgsql

import (
	"context"
	"errors"

	"cleantx/internal/domain"
	"cleantx/internal/pkg/sqldb"

	"github.com/jackc/pgx/v5"
)

type DoctorRepository struct {
	db commandExecutor
}

func NewDoctorRepository(db commandExecutor) *DoctorRepository {
	return &DoctorRepository{
		db: db,
	}
}

func (r *DoctorRepository) WithTx(tx sqldb.Tx) domain.DoctorRepository {
	return NewDoctorRepository(tx.(commandExecutor))
}

func (r *DoctorRepository) Get(ctx context.Context, id int) (*domain.Doctor, error) {
	doctor := domain.Doctor{ID: id}

	row := r.db.QueryRow(ctx, `select name, on_call from doctors where id = $1;`, id)

	if err := row.Scan(&doctor.Name, &doctor.OnCall); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrDoctorNotFound
		}

		return nil, err
	}

	return &doctor, nil
}

func (r *DoctorRepository) Update(ctx context.Context, doctor *domain.Doctor) error {
	_, err := r.db.Exec(ctx, `update doctors set name=$1, on_call=$2 where id=$3`,
		doctor.Name,
		doctor.OnCall,
		doctor.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *DoctorRepository) ListDoctorsOnCall(ctx context.Context) (domain.Doctors, error) {
	doctors := domain.Doctors{}

	rows, err := r.db.Query(ctx, `
		select id, name from doctors where on_call = true for update;
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		doctor := domain.Doctor{
			OnCall: true,
		}

		if err := rows.Scan(&doctor.ID, &doctor.Name); err != nil {
			return nil, err
		}

		doctors = append(doctors, &doctor)
	}

	return doctors, nil
}
