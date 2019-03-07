package main

import (
	"net"
	"os"

	"github.com/jasonlvhit/gocron"
	"github.com/navinds25/sftpmgmt/internal/app"
	"github.com/navinds25/sftpmgmt/internal/sftpgrpc"
	"github.com/navinds25/sftpmgmt/pkg/sftpdata"
	pb "github.com/navinds25/sftpmgmt/pkg/sftpevent"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func runGRPCServer(lis net.Listener, s *grpc.Server) error {
	pb.RegisterSftpServer(s, &sftpgrpc.Server{})
	log.Info("Started GRPC Server")
	if err := s.Serve(lis); err != nil {
		return err
	}
	defer s.GracefulStop()
	defer lis.Close()
	return nil
}

// Version for inserting version via ldflags
var Version string

func main() {
	// Update Cli values
	appCli := app.Cli()
	appCli.Version = Version
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
	defer sftpdata.Data.Config.CloseConfigDB()
	defer sftpdata.Data.Files.CloseFilesDB()

	// GRPC Server
	grpcAddress := "0.0.0.0:8432"
	grpclis, err := net.Listen("tcp", grpcAddress)
	defer grpclis.Close()
	if err != nil {
		log.Fatal(err)
	}

	grpcSftpI := grpc.NewServer()
	go func() {
		log.Info("Starting grpc server on address", grpcAddress)
		if err := runGRPCServer(grpclis, grpcSftpI); err != nil {
			log.Fatal(err)
		}
	}()

	// Run Tasks
	scheduleI := gocron.NewScheduler()
	scheduleI.Every(8).Seconds().Do(app.RunJobs)
	<-scheduleI.Start()
}
