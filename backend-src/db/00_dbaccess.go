package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type dbAccess struct {
	DB *gorm.DB
}

func DBaccess(dsn string) (*dbAccess, error) {
	// Establish the database connection using gorm.Open
	gormDB, err := gorm.Open(postgres.Open(dsn), nil)
	if err != nil {
		return nil, err
	}

	// If you need to add any custom configurations to gormDB, you can do it here.

	// Create and return the dbAccess instance with the gormDB connection
	return &dbAccess{DB: gormDB}, nil
}
