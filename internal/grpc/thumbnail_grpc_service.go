package grpcServer

import (
	"context"
	"log"

	"thumb/internal/thumbnail"
	"thumb/proto"

	database "thumb/internal/db"
)

type ThumbnailServiceServer struct {
	proto.UnimplementedThumbnailServiceServer
	database.ThumbnailDatabase
}

func (t *ThumbnailServiceServer) GetThumbnail(ctx context.Context, request *proto.ThumbnailRequest) (*proto.ThumbnailResponse, error) {
	videoURL := request.VideoUrl
	exists := t.ThumbnailDatabase.RecordExists(videoURL)
	if exists {
		log.Println("Sending cached result for", videoURL)
		thumb, err := t.ThumbnailDatabase.GetCached(videoURL)
		if err != nil {
			return nil, err
		}
		return &proto.ThumbnailResponse{ThumbnailData: thumb}, nil
	}
	thumb, err := thumbnail.GetThumbnail(videoURL)
	if err != nil {
		return nil, err
	}
	if err := t.ThumbnailDatabase.Save(videoURL, thumb); err != nil {
		log.Println("failed to cache thumbnail")
		return nil, err
	}
	log.Printf("Thumbnail for video %s successfuly cached.", videoURL)
	return &proto.ThumbnailResponse{ThumbnailData: thumb}, nil
}
