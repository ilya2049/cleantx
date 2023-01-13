package domain

import "context"

type Doctor struct {
	ID     int
	Name   string
	OnCall bool
}

func (d *Doctor) FinishShift() {
	d.OnCall = false
}

func (d *Doctor) TakeShift() {
	d.OnCall = true
}

type DoctorRepository interface {
	Get(ctx context.Context, id int) (*Doctor, error)
	Update(ctx context.Context, doctor *Doctor) error
}
