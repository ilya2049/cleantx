package cli

import (
	"cleantx/internal/app"
	"cleantx/internal/pgsql"

	"context"
	"flag"
	"log"
)

func ExecuteCommand(command func(context.Context, int, app.DoctorCommandExecutor) error) {
	var doctorID = flag.Int("id", 0, "doctor id")
	flag.Parse()

	dbConn, closeDBConn, err := pgsql.NewConnection()
	if err != nil {
		log.Println(err)

		return
	}

	defer closeDBConn()

	doctorRepository := pgsql.NewDoctorRepository(dbConn)
	unitOfWorkProvider := pgsql.NewUnitOfWorkProvider(dbConn)
	doctorCommandExecutor := app.NewDoctorCommandExecutor(unitOfWorkProvider, doctorRepository)

	if err := command(context.Background(), *doctorID, *doctorCommandExecutor); err != nil {
		log.Println(err)

		return
	}

	log.Println("OK")
}
