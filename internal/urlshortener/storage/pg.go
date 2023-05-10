package storage

import (
	"context"
	"errors"
	"fmt"
	"urlshortener/internal/urlshortener/encoder"
	"urlshortener/pkg/urlshortener/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type pgStorage struct {
	enc encoder.Encoder
	db  *gorm.DB
}

type urls struct {
	Id  uint64 `gorm:"index;primaryKey;autoIncrement:true"`
	Url string
}

// Returns `(nil, DatabaseError)` if fails to connect to the db.
func NewPgStorage(enc encoder.Encoder, dbCfg config.Config) (*pgStorage, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbCfg.DbHost, dbCfg.DbPort, dbCfg.DbUser, dbCfg.DbPassword, dbCfg.DbName)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, DatabaseError{err.Error()}
	}
	db.AutoMigrate(&urls{})

	return &pgStorage{enc, db}, nil
}

// Returns `("", DatabaseError)` if fails to create a record or saves changes in the db,
// `("", EncodingOverflowError)` if encoding overflow has occured.
func (s *pgStorage) Shorten(ctx context.Context, url string) (string, error) {
	var row urls
	result := s.db.WithContext(ctx).Model(&urls{}).Create(&row)
	if result.Error != nil {
		return "", DatabaseError{result.Error.Error()}
	}

	encoded, err := s.enc.Encode(row.Id)
	if err != nil {
		s.db.Model(&urls{}).Delete(&urls{}, row.Id)
		return "", err
	}
	row.Url = url

	result = s.db.WithContext(ctx).Save(&row)
	if result.Error != nil {
		return "", DatabaseError{result.Error.Error()}
	}

	return encoded, nil
}

// Returns `"", DecodingError` if `url` has invalid encoding,
// `"", UrlNotFoundError` if decoded value was not found in the db,
// `"", DatabaseError` if fails to query the db.
func (s *pgStorage) Unshorten(ctx context.Context, url string) (string, error) {
	id, err := s.enc.Decode(url)
	if err != nil {
		return "", err
	}

	var decodedUrl string
	result := s.db.WithContext(ctx).Model(urls{}).Select("url").Where("Id = ?", id).Take(&decodedUrl)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", UrlNotFoundError{result.Error.Error()}
		} else {
			return "", DatabaseError{result.Error.Error()}
		}
	}

	return decodedUrl, nil
}
