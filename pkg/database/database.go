package database

import (
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"nenvoy.com/pkg/constants"
)

// NewSession - Return the db object to create transactions on the database
func NewSession() (db *gorm.DB, err error) {
	// Connect and open the database
	db, err = gorm.Open(sqlite.Open(constants.DBPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return db, errors.Wrap(err, "failed to connect database")
	}

	return db, nil

}
