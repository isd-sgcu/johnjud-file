package database

import (
	"fmt"
	"strconv"

	"github.com/isd-sgcu/johnjud-file/cfgldr"
	"github.com/isd-sgcu/johnjud-file/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func InitPostgresDatabase(conf *cfgldr.Database, isDebug bool) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", conf.Host, strconv.Itoa(conf.Port), conf.Username, conf.Password, conf.Name, conf.SSL)

	gormConf := &gorm.Config{}

	if !isDebug {
		gormConf.Logger = gormLogger.Default.LogMode(gormLogger.Silent)
	}

	db, err = gorm.Open(postgres.Open(dsn), gormConf)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&model.Image{}, &model.Pet{})
	if err != nil {
		return nil, err
	}

	return
}
