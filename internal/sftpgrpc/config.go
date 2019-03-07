package sftpgrpc

import (
	"context"

	"github.com/navinds25/sftpmgmt/pkg/sftpconfig"
	"github.com/navinds25/sftpmgmt/pkg/sftpdata"
	pb "github.com/navinds25/sftpmgmt/pkg/sftpevent"
)

func pbtogoTransferConfig(pbtc *pb.SftpTransferConfig) *sftpconfig.TransferConfig {
	goTransferConf := &sftpconfig.TransferConfig{}
	goTransferConf.TransferID = pbtc.Transferid
	goTransferConf.Description = pbtc.Description
	goTransferConf.Type = pbtc.Type.String()
	goTransferConf.Source.Remote.RemotePath = pbtc.Remotepath
	goTransferConf.Source.Remote.Host = pbtc.Remotehost
	goTransferConf.Source.Remote.Port = int(pbtc.Remoteport)
	goTransferConf.Source.Remote.Auth.Username = pbtc.Remoteuser
	goTransferConf.Source.Remote.Auth.Password = pbtc.Remotepassword
	goTransferConf.Source.Remote.Auth.Key = pbtc.Remotekey
	goTransferConf.Destination.Local.DestFile = pbtc.Localfile
	goTransferConf.Destination.Local.DestPath = pbtc.Localpath
	return goTransferConf
}

// AddConfig adds a config to config db.
func (s *Server) AddConfig(ctx context.Context, tc *pb.SftpTransferConfig) (*pb.Ack, error) {
	transferConfig := pbtogoTransferConfig(tc)
	if err := sftpdata.Data.Config.AddSFTPEntry(transferConfig); err != nil {
		return &pb.Ack{
			Error: err.Error(),
		}, err
	}
	return &pb.Ack{
		Message: "Added config successfully.",
	}, nil
}
