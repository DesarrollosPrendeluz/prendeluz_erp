package repositories

//
// import (
// 	"log"
//
// 	"github.com/DATA-DOG/go-sqlmock"
// 	"gorm.io/driver/mysql"
// 	"gorm.io/driver/sqlite"
// 	"gorm.io/gorm"
// )
//
// func SetUpTestDB(models ...interface{}) (*gorm.DB, func()) {
// 	db, mock, err := sqlmock.New()
// 	gormDb, err := gorm.Open(mysql.New(mysql.Config{Conn: db, SkipInitializeWithVersion: true}), &gorm.Config{})
//
// 	// if err != nil {
// 	// 	log.Panic("Failed to connect Database")
// 	// }
// 	//
// 	// err = db.AutoMigrate(models...)
// 	// if err != nil {
// 	// 	log.Panic(err)
// 	// }
// 	// cleanup := func() {
// 	// 	for _, model := range models {
// 	// 		db.Migrator().DropTable(model)
// 	// 	}
// 	// }
//
// 	return db, cleanup
// }
