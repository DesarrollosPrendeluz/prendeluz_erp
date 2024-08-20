package repositories

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetUpTestDB(models ...interface{}) (*gorm.DB, func()) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	if err != nil {
		log.Panic("Failed to connect Database")
	}

	err = db.AutoMigrate(models...)
	if err != nil {
		log.Panic(err)
	}
	cleanup := func() {
		for _, model := range models {
			db.Migrator().DropTable(model)
		}
	}

	return db, cleanup
}
