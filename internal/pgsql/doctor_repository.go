package pgsql

import (
	"context"
	"errors"
	"log"

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

func (*DoctorRepository) updateWithTx(ctx context.Context, tx pgx.Tx, doctor *domain.Doctor) error {
	_, err := tx.Exec(ctx, `update doctors set name=$1, on_call=$2 where id=$3`,
		doctor.Name,
		doctor.OnCall,
		doctor.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (*DoctorRepository) listDoctorsOnCall(ctx context.Context, tx pgx.Tx) (domain.Doctors, error) {
	doctors := domain.Doctors{}

	rows, err := tx.Query(ctx, `
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

func (r *DoctorRepository) FinishShift(ctx context.Context, doctorID int) (err error) {
	tx, err := r.db.Begin(ctx)
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

	doctorsOnCall, err := r.listDoctorsOnCall(ctx, tx)
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

	if err := r.updateWithTx(ctx, tx, doctor); err != nil {
		return err
	}

	return nil
}
