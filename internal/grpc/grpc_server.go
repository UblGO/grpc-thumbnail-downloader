package grpcServer

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"thumb/proto"

	database "thumb/internal/db"

	"google.golang.org/grpc"
)

const dbFileName = "thumbnails.db"

type GrpcServer struct {
	Serv *grpc.Server
}

// New gRPC server
func GrpcServerBuilder() *GrpcServer {
	return &GrpcServer{
		Serv: grpc.NewServer(),
	}
}

// Run server
func RunServer() {
	thDB := database.ThumbnailDatabaseBuilder(dbFileName)
	server := GrpcServerBuilder()
	srv := &ThumbnailServiceServer{ThumbnailDatabase: *thDB}
	proto.RegisterThumbnailServiceServer(server.Serv, srv)
	errs := make(chan error, 2)
	go func() {
		errs <- server.ServerAndListen()
	}()
	//graceful stop realization
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		errs <- server.ServerGracefulStop(signalChan)
	}()
	err := <-errs
	if err != nil {
		log.Fatal("Failed to start server", err)
	}
}

// Creates new listner
func (s *GrpcServer) ServerAndListen() error {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}
	log.Println("Listening on port :8080")
	return s.Serv.Serve(listener)
}

// Graceful stop implementation
func (s *GrpcServer) ServerGracefulStop(signalChan <-chan os.Signal) error {
	sig := <-signalChan
	log.Println(sig.String(), "\nGracefully stopping...")
	s.Serv.GracefulStop()
	return nil
}
