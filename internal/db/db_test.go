package database

import (
	"errors"
	"log"
	"os"
	"testing"

	testData "thumb/assets"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// db files for each testFunc, benefits:
// 1.testing with -race flag
// 2.Deleting each test database after test is done
// 3.Excluding unqiue constraint violation error
const (
	newThumbDb     = "NewThumbnailDatabase_test.db"
	saveThumbDb    = "SaveThumbnail_test.db"
	cachedThumbDb  = "GetCachedThumbnail_test.db"
	recordExistsDb = "RecordExists_test.db"
)

// Place holder for thumbnail byte data,
// each element equals it index in string format for easy tracing
var thumbsPlaceHolder = [][]byte{{48}, {49}, {50}, {51}}

var videos = []Video{
	{VideoURL: "https://www.youtube.com/watch?v=y0sF5xhGreA", Thumbnail: thumbsPlaceHolder[0]},
	{VideoURL: "https://www.youtube.com/watch?v=MlDtL2hIj-Q", Thumbnail: thumbsPlaceHolder[1]},
	{VideoURL: "https://www.youtube.com/watch?v=DixdzXAFS18", Thumbnail: thumbsPlaceHolder[2]},
	{VideoURL: "https://www.youtube.com/watch?v=KvE92fCMbmc", Thumbnail: thumbsPlaceHolder[3]},
}

func TestNewThumbnailDatabase(t *testing.T) {
	testDb := ThumbnailDatabaseBuilder(newThumbDb)
	defer deleteDatabase(testDb, newThumbDb)
	//test if db file created
	_, err := os.Stat(newThumbDb)
	exists := !errors.Is(err, os.ErrNotExist)
	require.True(t, exists)
	//test if table created
	tableVideosCreated := testDb.db.Migrator().HasTable(&Video{})
	require.True(t, tableVideosCreated)
}

func TestSaveThumbnail(t *testing.T) {
	testDb := ThumbnailDatabaseBuilder(saveThumbDb)
	defer deleteDatabase(testDb, saveThumbDb)
	sqlDb, err := testDb.db.DB()
	if err != nil {
		log.Fatal("Failed to recieve a *sql.DB")
	}

	for _, video := range videos {
		err := testDb.Save(video.VideoURL, video.Thumbnail)
		assert.NoError(t, err)
	}

	for _, video := range videos {
		var tempVideo Video
		//Executing raw sql,
		res := sqlDb.QueryRow("SELECT video_url, thumbnail FROM videos WHERE video_url=?;", video.VideoURL)
		err = res.Scan(&tempVideo.VideoURL, &tempVideo.Thumbnail)
		assert.NoError(t, err)
		assert.Equal(t, video.VideoURL, tempVideo.VideoURL)
		assert.Equal(t, video.Thumbnail, tempVideo.Thumbnail)
	}
}

func TestGetCachedThumbnail(t *testing.T) {
	testDb := ThumbnailDatabaseBuilder(cachedThumbDb)
	defer deleteDatabase(testDb, cachedThumbDb)

	for _, video := range videos {
		err := testDb.Save(video.VideoURL, video.Thumbnail)
		assert.NoError(t, err)
	}

	for _, video := range videos {
		thumb, err := testDb.GetCached(video.VideoURL)
		assert.NoError(t, err)
		assert.Equal(t, video.Thumbnail, thumb)
	}
}

func TestRecordExists(t *testing.T) {
	testDb := ThumbnailDatabaseBuilder(recordExistsDb)
	defer deleteDatabase(testDb, recordExistsDb)

	for _, video := range videos {
		err := testDb.Save(video.VideoURL, video.Thumbnail)
		assert.NoError(t, err)
	}

	for _, video := range videos {
		exists := testDb.RecordExists(video.VideoURL)
		assert.True(t, exists)
	}

	for _, url := range testData.InvalidURLs {
		exists := testDb.RecordExists(url)
		assert.False(t, exists)
	}
}

// Closes connection and deletes test database
func deleteDatabase(dbStr *ThumbnailDatabase, dbFileName string) {
	if err := dbStr.Close(); err != nil {
		log.Fatalf("%v\nFailed to close connection to test database: %s\nPlease delete manually before future testing", err, dbFileName)
	}
	if err := os.Remove("../db/" + dbFileName); err != nil {
		log.Fatalf("%v\nFailed to delete test database: %s\nPlease delete manually before future testing", err, dbFileName)
	}
}
