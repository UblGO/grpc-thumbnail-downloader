package thumbnail

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// Gets thumbnail by url
func GetThumbnail(videoURL string) ([]byte, error) {
	log.Println("Getting thumbnail for ", videoURL)
	htmlString, err := getHTMLString(videoURL)
	if err != nil {
		return nil, err
	}
	thumbnailURL, err := getThumbnailURL(htmlString)
	if err != nil {
		return nil, err
	}
	thumbnailData, err := downloadThumbnail(thumbnailURL)
	if err != nil {
		return nil, err
	}

	return thumbnailData, nil

}

// Gets raw HTML by url
func getHTMLString(videoURL string) (string, error) {
	response, err := http.Get(videoURL)
	if err != nil {
		return "", err
	}
	htmlString, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	return string(htmlString), nil
}

// Gets thumbnail url from raw HTML
func getThumbnailURL(htmlString string) (string, error) {
	var thumbnailURL string
	startTag := "<meta property=\"og:image\" content=\""
	endTag := "\">"
	startIndex := strings.Index(htmlString, startTag)
	if startIndex == -1 {
		return "", fmt.Errorf("thumbnail tag not found")
	}
	htmlString = htmlString[startIndex+len(startTag):]
	endIndex := strings.Index(htmlString, endTag)
	if endIndex == -1 {
		return "", fmt.Errorf("thumbnail tag not found")
	}
	thumbnailURL = htmlString[:endIndex]
	return thumbnailURL, nil
}

// Gets thumbnail in []byte format
func downloadThumbnail(thumbnailURL string) ([]byte, error) {
	response, err := http.Get(thumbnailURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code is not 200\trecieved code: %d", response.StatusCode)
	}
	thumbnailData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return thumbnailData, err
}
