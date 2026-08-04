package db

import "gorm.io/gorm"

var PostgresConn *gorm.DB

func GetPostgresConnection() *gorm.DB {
	return PostgresConn
}

func InitPostgres() CloseConn {
	var err error
	db, err := NewDatabase()

	if err != nil {
		panic(err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}

	PostgresConn = db

	return closeConn(sqlDb)
}
