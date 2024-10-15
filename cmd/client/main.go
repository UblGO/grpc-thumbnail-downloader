package main

import (
	"bytes"
	"context"
	"flag"
	"image"
	"image/jpeg"
	"log"
	"os"
	"sync"

	"thumb/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var async bool         // async flag
var videoUrls []string //slice of user provided URL's

func init() {
	//setting up  async flag
	flag.BoolVar(&async, "async", false, "async thumbnail download")
	//checking if images dir exists, creating if not
	if _, err := os.Stat("./images"); err != nil {
		if err := os.Mkdir("images", 0777); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	RunClient()
}

func RunClient() {
	flag.Parse()
	//getting videos URL's and checking if there any
	videoUrls = flag.Args()
	if len(videoUrls) < 1 {
		log.Fatal("no video url provided")
	}
	//setting up gRPC client
	conn, err := grpc.NewClient("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
		return
	}
	client := proto.NewThumbnailServiceClient(conn)
	log.Printf("downloading %d thumbnails", len(videoUrls))
	switch async {
	case true: //async downloading
		wg := sync.WaitGroup{}
		for _, url := range videoUrls {
			wg.Add(1)
			go func() {
				defer wg.Done()
				getThumbnailImage(url, client)
			}()
		}
		wg.Wait()
	case false: //sequential downloading
		for _, url := range videoUrls {
			getThumbnailImage(url, client)

		}
	}
}

func getThumbnailImage(url string, client proto.ThumbnailServiceClient) {
	//youtube video link in format https://www.youtube.com/watch?v=/*unique_video_code*/
	//is 32 symbols long + unique video code which length is always 11
	if len(url) != 43 {
		log.Println("wrong link", url)
		return
	}
	response, err := client.GetThumbnail(context.Background(),
		&proto.ThumbnailRequest{VideoUrl: url})
	log.Println("getting thumbnail for", url)

	if err != nil {
		log.Printf("no response: %v", err)
		return
	}
	//using unique video code as file name
	fileName := url[32:] + ".jpg"
	if err = saveThumbnailImage(response.ThumbnailData, fileName); err != nil {
		log.Printf("failed to save thumbnail for %s: %v", url, err)
		return
	}
	log.Printf("Thumbnail for %s saved as %s", url, fileName)
}

func saveThumbnailImage(thumbnailData []byte, fileName string) error {
	img, _, _ := image.Decode(bytes.NewReader(thumbnailData))
	out, err := os.Create("./images/" + fileName)
	if err != nil {
		return err
	}
	defer out.Close()
	//image quality 0-100
	opt := jpeg.Options{Quality: 100}
	err = jpeg.Encode(out, img, &opt)
	if err != nil {
		return err
	}

	return nil
}
