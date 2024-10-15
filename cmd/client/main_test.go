package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"testing"
	"thumb/proto"

	testData "thumb/assets"
	database "thumb/internal/db"
	rpc "thumb/internal/grpc"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

var lis *bufconn.Listener
var ThDB *database.ThumbnailDatabase
var savedTestImages []string

const bufSize = 1024 * 1024
const grpcClientDb = "grpcClient_test.db"

func init() {
	lis = bufconn.Listen(bufSize)
	ThDB = database.ThumbnailDatabaseBuilder(grpcClientDb)
	s := rpc.GrpcServerBuilder()
	proto.RegisterThumbnailServiceServer(s.Serv, &rpc.ThumbnailServiceServer{ThumbnailDatabase: *ThDB})
	go func() {
		if err := s.Serv.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}
func TestGetThumbnailImage(t *testing.T) {
	defer deleteDatabase(ThDB, grpcClientDb)
	conn, err := grpc.NewClient("127.0.0.1:8080", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := proto.NewThumbnailServiceClient(conn)
	for _, url := range testData.ValidURLs {
		getThumbnailImage(url, client)
		fileName := url[32:] + ".jpg"
		path := "./images/" + fileName
		_, err := os.Stat(path)
		exists := !errors.Is(err, os.ErrNotExist)
		require.True(t, exists)
		savedTestImages = append(savedTestImages, path)
	}
	for _, v := range savedTestImages {
		if err := os.Remove(v); err != nil {
			log.Fatal("failed to delete test image ", err)
		}
	}

}

// Closes connection and deletes test database
func deleteDatabase(dbStr *database.ThumbnailDatabase, dbFileName string) {
	if err := dbStr.Close(); err != nil {
		log.Fatalf("%v\nFailed to close connection to test database: %s\nPlease delete manually before future testing", err, dbFileName)
	}
	if err := os.Remove("../../internal/db/" + dbFileName); err != nil {
		log.Fatalf("%v\nFailed to delete test database: %s\nPlease delete manually before future testing", err, dbFileName)
	}
}
