package main

import (
	"cleantx/internal/app"
	"cleantx/internal/cli"

	"context"
)

func main() {
	cli.ExecuteCommand(func(
		ctx context.Context,
		doctorID int,
		doctorCommandExecutor app.DoctorCommandExecutor,
	) error {
		return doctorCommandExecutor.FinishShift(ctx, doctorID)
	})
}
