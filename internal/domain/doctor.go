package domain

import (
	"context"
	"errors"
)

var (
	ErrDoctorNotFound         = errors.New("doctor not found")
	ErrAtLeastOneDoctorOnCall = errors.New("failed to finish shift: at least one doctor must be on call")
)

type Doctor struct {
	ID     int
	Name   string
	OnCall bool
}

func (d *Doctor) FinishShift(doctorsOnCallCount int) error {
	if doctorsOnCallCount == 1 {
		return ErrAtLeastOneDoctorOnCall
	}

	d.OnCall = false

	return nil
}

func (d *Doctor) TakeShift() {
	d.OnCall = true
}

type Doctors []*Doctor

func (doctors Doctors) Get(id int) *Doctor {
	for _, doctor := range doctors {
		if doctor.ID == id {
			return doctor
		}
	}

	return nil
}

type DoctorRepository interface {
	Get(ctx context.Context, id int) (*Doctor, error)
	Update(ctx context.Context, doctor *Doctor) error
	FinishShift(ctx context.Context, doctorID int) error
}
