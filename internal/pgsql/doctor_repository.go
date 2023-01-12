package pgsql

import (
	"context"

	"cleantx/internal/domain"

	"github.com/jackc/pgx/v5"
)

type DoctorRepository struct {
	db *pgx.Conn
}

func NewDoctorRepository(db *pgx.Conn) *DoctorRepository {
	return &DoctorRepository{
		db: db,
	}
}

func (r *DoctorRepository) Get(ctx context.Context, id int) (*domain.Doctor, error) {
	doctor := domain.Doctor{ID: id}

	row := r.db.QueryRow(ctx, `select name, is_on_shift from doctors where id = $1;`, id)

	if err := row.Scan(&doctor.Name, &doctor.IsOnSift); err != nil {
		return nil, err
	}

	return &doctor, nil
}

func (r *DoctorRepository) Update(ctx context.Context, doctor *domain.Doctor) error {
	_, err := r.db.Exec(ctx, `update doctors set name=$1, is_on_shift=$2 where id=$3`,
		doctor.Name,
		doctor.IsOnSift,
		doctor.ID,
	)

	if err != nil {
		return err
	}

	return nil
}
