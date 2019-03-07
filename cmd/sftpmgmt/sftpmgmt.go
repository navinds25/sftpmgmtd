package main

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"

	pb "github.com/navinds25/sftpmgmt/pkg/sftpevent"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("0.0.0.0:8432", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	c := pb.NewSftpClient(conn)
	deadline := time.Now().Add(time.Duration(1000000) * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	transferTypeValue := pb.SftpTransferConfig_TransferType_value["Pull"]
	transferType := pb.SftpTransferConfig_TransferType(transferTypeValue)
	msg, err := c.AddConfig(ctx, &pb.SftpTransferConfig{
		Transferid: "test1",
		Type:       transferType,
	})
	//msg, err := c.GRPCTest(ctx, &pb.Ack{Message: "sending this message"})
	if err != nil {
		log.Fatal(err)
	}
	log.Info(msg)
}
