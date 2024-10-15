package grpcServer

import (
	"context"
	"log"
	"net"
	"os"
	"testing"

	testData "thumb/assets"
	database "thumb/internal/db"
	proto "thumb/proto"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024
const grpcGetThumbDb = "grpcGetThumbDb_test.db"

var lis *bufconn.Listener
var ThDB *database.ThumbnailDatabase

func init() {
	lis = bufconn.Listen(bufSize)
	ThDB = database.ThumbnailDatabaseBuilder(grpcGetThumbDb)
	s := GrpcServerBuilder()
	proto.RegisterThumbnailServiceServer(s.Serv, &ThumbnailServiceServer{ThumbnailDatabase: *ThDB})
	go func() {
		if err := s.Serv.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestGetThumbnail(t *testing.T) {
	defer deleteDatabase(ThDB, grpcGetThumbDb)
	conn, err := grpc.NewClient("127.0.0.1:8080", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := proto.NewThumbnailServiceClient(conn)
	for _, url := range testData.ValidURLs {
		resp, err := client.GetThumbnail(context.Background(), &proto.ThumbnailRequest{VideoUrl: url})
		require.NoError(t, err)
		require.NotEmpty(t, resp.ThumbnailData)
	}
	for _, url := range testData.InvalidURLs {
		resp, err := client.GetThumbnail(context.Background(), &proto.ThumbnailRequest{VideoUrl: url})
		require.Error(t, err)
		require.Empty(t, resp)
	}
	for _, url := range testData.BrokenLinks {
		resp, err := client.GetThumbnail(context.Background(), &proto.ThumbnailRequest{VideoUrl: url})
		require.Error(t, err)
		require.Empty(t, resp)
	}
}

// Closes connection and deletes test database
func deleteDatabase(dbStr *database.ThumbnailDatabase, dbFileName string) {
	if err := dbStr.Close(); err != nil {
		log.Fatalf("%v\nFailed to close connection to test database: %s\nPlease delete manually before future testing", err, dbFileName)
	}
	if err := os.Remove("../db/" + dbFileName); err != nil {
		log.Fatalf("%v\nFailed to delete test database: %s\nPlease delete manually before future testing", err, dbFileName)
	}
}
