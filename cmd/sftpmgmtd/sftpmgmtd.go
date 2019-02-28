package main

import (
	"os"

	"github.com/jasonlvhit/gocron"
	"github.com/navinds25/sftpmgmt/internal/app"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Update Cli values
	appCli := app.Cli()
	if err := appCli.Run(os.Args); err != nil {
		log.Fatal(err)
	}

	//Default values for Cli
	if err := app.CliVal.GetCliFlags(); err != nil {
		log.Fatal(err)
	}

	// Setup service if flag is true
	if app.CliVal.SetupService {
		status, err := app.DaemonSetup()
		if err != nil {
			log.Fatal(err)
		}
		log.Info(status)
		os.Exit(0)
	}

	// Setup DBs
	if err := app.DBSetup(); err != nil {
		log.Fatal(err)
	}

	// Cli database updates

	// Run Tasks
	s := gocron.NewScheduler()
	s.Every(8).Seconds().Do(app.RunJobs)
	<-s.Start()
}
