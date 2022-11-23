package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jrmarcco/go-backend-demo/util"
	gormMysql "gorm.io/driver/mysql"
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
	cfg, err := util.LoadConfig("..")
	if err != nil {
		log.Fatal("can not load config:", err)
	}

	sqlDB, err := sql.Open(cfg.Db.Driver, cfg.Db.Source)
	if err != nil {
		log.Fatal(err)
	}
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(20)

	db, err = gorm.Open(gormMysql.New(gormMysql.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
}
