package storage

import (
	"errors"
	"fmt"
	"urlshortener/internal/urlshortener/encoder"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type pgStorage struct {
	enc encoder.Encoder
	db  *gorm.DB
}

type DbConfig struct {
	Hostname string
	Port     uint
	User     string
	Password string
	DbName   string
}

type urls struct {
	Id  uint64 `gorm:"index;primaryKey;autoIncrement:true"`
	Url string
}

func NewPgStorage(enc encoder.Encoder, dbCfg DbConfig) (*pgStorage, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbCfg.Hostname, dbCfg.Port, dbCfg.User, dbCfg.Password, dbCfg.DbName)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, DatabaseError{err.Error()}
	}
	db.AutoMigrate(&urls{})

	return &pgStorage{enc, db}, nil
}

func (s *pgStorage) Close() {
	sqlDb, err := s.db.DB()
	if err != nil {
		return
	}
	sqlDb.Close()
}

func (s *pgStorage) Create(url string) (string, error) {
	var row urls
	result := s.db.Model(&urls{}).Create(&row)
	if result.Error != nil {
		return "", DatabaseError{result.Error.Error()}
	}

	encoded, err := s.enc.Encode(row.Id)
	if err != nil {
		s.db.Model(&urls{}).Delete(&urls{}, row.Id)
		return "", err
	}
	row.Url = url

	result = s.db.Save(&row)
	if result.Error != nil {
		panic(result.Error.Error())
	}

	return encoded, nil
}

func (s *pgStorage) Get(url string) (string, error) {
	id, err := s.enc.Decode(url)
	if err != nil {
		return "", err
	}

	var decodedUrl string
	result := s.db.Model(urls{}).Select("url").Where("Id = ?", id).Take(&decodedUrl)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", UrlNotFoundError{result.Error.Error()}
		} else {
			return "", DatabaseError{result.Error.Error()}
		}
	}

	return decodedUrl, nil
}
