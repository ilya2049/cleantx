package domain

import "context"

type Doctor struct {
	ID       int
	Name     string
	IsOnSift bool
}

func (d *Doctor) FinishShift() {
	d.IsOnSift = false
}

func (d *Doctor) TakeShift() {
	d.IsOnSift = true
}

type DoctorRepository interface {
	Get(ctx context.Context, id int) (*Doctor, error)
	Update(ctx context.Context, doctor *Doctor) error
}
