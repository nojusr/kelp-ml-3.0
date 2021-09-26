package kelp

import (
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitializeDatabase() {
	db, err := gorm.Open(sqlite.Open(os.Getenv("KELP_DB_PATH")), &gorm.Config{})

	if err != nil {
		panic("failed to connect to db")
	}

	db.AutoMigrate(&KelpUser{}, &KelpInvite{}, &KelpPaste{}, &KelpFile{})

	Db = db
}
