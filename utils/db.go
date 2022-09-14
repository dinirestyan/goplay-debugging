package utils

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type dbUtils struct {
	db *gorm.DB
}

var dbInstance *dbUtils
var dbOnce sync.Once

func GetDBConnection() *gorm.DB {
	dbOnce.Do(func() {
		log.Println("Initialize db connection...")
		connection := "host=" + os.Getenv("DATABASE_HOST") + " port=" + os.Getenv("DATABASE_PORT") + " user=" + os.Getenv("USERNAME_DB") + " dbname=" + os.Getenv("DATABASE_NAME") +
			" password=" + os.Getenv("PASSWORD_DB") + " sslmode=" + os.Getenv("DATABASE_SSL")
		log.Println(connection)
		db, err := gorm.Open(os.Getenv("DATABASE_TYPE"), connection)

		if err != nil {
			log.Println(err)
			return
		}

		db.DB().SetConnMaxLifetime(time.Second * 60)
		db.SingularTable(true)
		db.LogMode(true)

		if err != nil {
			log.Println(err)
		}

		dbInstance = &dbUtils{
			db: db,
		}
	})

	return dbInstance.db
}
