package app

import (
	"os"
	"path"

	"github.com/dgraph-io/badger"
	"github.com/navinds25/sftpmgmt/pkg/sftpdata"
	log "github.com/sirupsen/logrus"
	"github.com/takama/daemon"
)

// DaemonSetup sets up the daemon
func DaemonSetup() (string, error) {
	service, err := daemon.New("sftpmgmtd", "SFTP Management Daemon")
	if err != nil {
		return "", err
	}
	status, err := service.Install()
	if err != nil {
		return status, err
	}
	return status, err
}

func createDataDir(datadir string) error {
	log.Debug("datadir:", datadir)
	if err := os.MkdirAll(datadir, 0755); err != nil {
		return err
	}
	log.Info("created data directory: ", datadir)
	return nil
}

// DBSetup sets up and intializes the DB.
func DBSetup() error {
	log.Info("Starting DB Setup")
	configDBDir := path.Join(CliVal.DataDir, "config")
	filesDBDir := path.Join(CliVal.DataDir, "files")
	_, Err := os.Stat(CliVal.DataDir)
	if Err != nil {
		if err := createDataDir(CliVal.DataDir); err != nil {
			return err
		}
	}
	configDBopts := badger.DefaultOptions
	configDBopts.Dir = configDBDir
	configDBopts.ValueDir = path.Join(configDBDir, "value")
	//configDBopts.Logger = log.StandardLogger()
	configDB, err := badger.Open(configDBopts)
	if err != nil {
		return err
	}
	sftpdata.InitConfigDB(sftpdata.BadgerDB{
		ConfigDB: configDB,
	})
	filesDBopts := badger.DefaultOptions
	filesDBopts.Dir = filesDBDir
	filesDBopts.ValueDir = path.Join(filesDBDir, "value")
	//filesDBopts.Logger = log.StandardLogger()
	filesDB, err := badger.Open(filesDBopts)
	if err != nil {
		return err
	}
	sftpdata.InitFilesDB(sftpdata.BadgerDB{
		FilesDB: filesDB,
	})
	log.Info("Databases initialized.")
	return nil
}

// RunJobs running the configured jobs via gocron.
func RunJobs() error {
	allConfigs, err := sftpdata.Data.Config.GetAll()
	if err != nil {
		return err
	}
	log.Debug("All configs", allConfigs)
	for i, config := range allConfigs {
		log.Println("processing config: ", i, config)
	}
	return nil
}
