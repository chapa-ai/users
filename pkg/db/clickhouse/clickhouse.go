package clickhouse

import (
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"log"
	"time"
)

type Logs struct {
	Name      string
	Age       string
	Timestamp time.Time
}

func Clickhouse() (*gorm.DB, error) {
	dsn := "tcp://localhost:9000?database=default&username=default"

	db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Error connecting to click: %v", err)
		return nil, err
	}

	err = db.AutoMigrate(&Logs{})
	if err != nil {
		log.Printf("Error migrating to click: %v", err)
		return nil, err
	}

	db.Set("gorm:table_options", "ENGINE=Distributed(cluster, default, hits)").AutoMigrate(&Logs{})

	return db, nil
}
