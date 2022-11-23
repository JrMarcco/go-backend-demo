package models

import (
	"github.com/jrmarcco/go-backend-demo/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

type Model struct {
	ID        uint64 `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

var db *gorm.DB

func Setup() {
	config, err := util.LoadConfig("..")
	if err != nil {
		log.Fatal("can not load config:", err)
	}

	db, err = gorm.Open(mysql.Open(config.Db.Source), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)
}
