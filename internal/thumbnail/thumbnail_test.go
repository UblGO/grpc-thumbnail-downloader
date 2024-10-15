package thumbnail

import (
	"net/http"
	"testing"
	testData "thumb/assets"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHTMLString(t *testing.T) {
	for _, url := range testData.ValidURLs {
		html, err := getHTMLString(url)
		assert.NoError(t, err)
		assert.NotEmpty(t, html)
	}
	for _, url := range testData.BrokenLinks {
		html, err := getHTMLString(url)
		assert.Error(t, err)
		assert.Empty(t, html)
	}

}

func TestGetThumbnailURL(t *testing.T) {
	for _, url := range testData.ValidURLs {
		html, err := getHTMLString(url)
		require.NoError(t, err)
		thumbURL, err := getThumbnailURL(html)
		assert.NoError(t, err)
		require.NotEmpty(t, thumbURL)
		resp, err := http.Get(thumbURL)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
}

func TestDownloadThumbnail(t *testing.T) {
	for _, url := range testData.ValidURLs {
		html, err := getHTMLString(url)
		require.NoError(t, err)
		thumbURL, err := getThumbnailURL(html)
		require.NoError(t, err)
		thumb, err := downloadThumbnail(thumbURL)
		assert.NoError(t, err)
		assert.NotEmpty(t, thumb)
	}
}

func TestGetThumbnail(t *testing.T) {

	for _, url := range testData.ValidURLs {
		bytes, err := GetThumbnail(url)
		assert.NoError(t, err)
		assert.NotEmpty(t, bytes)
	}
	for _, url := range testData.InvalidURLs {
		bytes, err := GetThumbnail(url)
		assert.Error(t, err)
		assert.Empty(t, bytes)
	}
}
