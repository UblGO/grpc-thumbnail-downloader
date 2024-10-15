package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	VideoURL  string `gorm:"unique;type:varchar(255);not null"`
	Thumbnail []byte `gorm:"type:bytes;not null"`
}

type ThumbnailDatabase struct {
	db *gorm.DB
}

// Creates new ThumbnailDatabase struct, creates new db file if necsesary and opens connection to it.
func ThumbnailDatabaseBuilder(dbFileName string) *ThumbnailDatabase {
	db, err := gorm.Open(sqlite.Open("../../internal/db/"+dbFileName), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	//create table if it doesn't exists
	if !db.Migrator().HasTable(&Video{}) {
		if err := db.AutoMigrate(&Video{}); err != nil {
			log.Fatal(err)
		}
	}
	return &ThumbnailDatabase{db: db}
}

// Check if row exists in database.
func (d *ThumbnailDatabase) RecordExists(videoUrl string) bool {
	var exists bool
	err := d.db.Model(&Video{}).
		Select("count(*) > 0").
		Where("video_url = ?", videoUrl).
		Find(&exists).
		Error
	if err != nil {
		return false
	}
	return exists
}

// Get a thumbnail from database.
func (d *ThumbnailDatabase) GetCached(videoUrl string) ([]byte, error) {
	var video Video
	tx := d.db.Table("videos").Select("thumbnail").Where("video_url=?", videoUrl).Take(&video)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return video.Thumbnail, nil
}

// Save a thumbnail to database.
func (d *ThumbnailDatabase) Save(videoUrl string, thumbnailData []byte) error {
	video := Video{VideoURL: videoUrl, Thumbnail: thumbnailData}
	tx := d.db.Create(&video)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// Close DB connection
func (d *ThumbnailDatabase) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Close()
	return err
}
