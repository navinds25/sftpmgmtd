package app

import (
	"os"
	"path"
	"path/filepath"

	"github.com/urfave/cli"
)

// CliVal is an instance of struct for the main cli flags
var CliVal CliFlags

// CliFlags is a struct for the main cli flags
type CliFlags struct {
	SetupService bool
	DataDir      string
	Debug        bool
}

// GetCliFlags is a Factory function returning the struct for MainCli
func (cliflags *CliFlags) GetCliFlags() error {
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}
	cliflags.DataDir = path.Join(currentDir, "data")
	return nil
}

// Cli for all commandline arguments.
func Cli() *cli.App {
	app := cli.NewApp()
	app.Name = "sftpmgmt"
	app.Usage = "For external sftp connections"
	app.Commands = []cli.Command{}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "setup-service",
			Destination: &CliVal.SetupService,
		},
		cli.StringFlag{
			Name:        "datadir",
			Usage:       "provide the data dir",
			Destination: &CliVal.DataDir,
		},
		cli.BoolFlag{
			Name:        "debug",
			Destination: &CliVal.Debug,
		},
	}
	return app
}

// TO Be Handled in Cli:
/*
1. Add SFTP Transfer Config
2. Delete SFTP Transfer Config
3. List All Config
4. Delete files from Files DB
*/
