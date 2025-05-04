package model

import (
	"strconv"
	"to-read/utils/logs"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	Hostname string `yaml:"hostname"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	SslMode  bool   `yaml:"sslMode"`
	TimeZone string `yaml:"timeZone"`
}

func Connect(d Database) error {
	const DbName = "toread"
	isSSL := "disable"
	if d.SslMode {
		isSSL = "enable"
	}
	dsn := "host=" + d.Hostname +
		" user=" + d.User +
		" password=" + d.Password +
		" dbname=" + DbName +
		" port=" + strconv.Itoa(d.Port) +
		" sslmode=" + isSSL +
		" TimeZone=" + d.TimeZone
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logs.Error("Failed to connect database "+DbName+".", zap.Error(err))
	}
	return err
}

// Almost used for creating tables
func AutoMigrateTable(dst ...interface{}) error {
	for _, d := range dst {
		err := db.Migrator().AutoMigrate(d)
		if err != nil {
			logs.Error("Auto migrate table failed.", zap.Any("model", d), zap.Error(err))
			return err
		}
		logs.Debug("Auto migrate table successfully. ", zap.Any("model", d))
	}
	return nil
}
