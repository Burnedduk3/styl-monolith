package aws

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	gormLogger "gorm.io/gorm/logger"
	"styl-monolith/pkg/logger"
	"sync"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Singleton struct for database instance
type Database struct {
	conn *gorm.DB
}

var (
	instance *Database
	once     sync.Once
)

// GetDatabaseInstance returns the single instance of Database
func GetDatabaseInstance(logs *logrus.Logger) *Database {
	logs.Debug(logger.OpeningDatabaseConnection)
	once.Do(func() {
		host := viper.GetString("DB_HOST")
		user := viper.GetString("DB_USER")
		password := viper.GetString("DB_PASSWORD")
		dbname := viper.GetString("DB_NAME")
		port := viper.GetString("DB_PORT")
		sslmode := viper.GetString("DB_SSLMODE")
		log.Debug(fmt.Sprintf(logger.ConnectionVariables, user, host, port, dbname, sslmode))
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
			host, user, password, dbname, port, sslmode)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: gormLogger.Default.LogMode(gormLogger.Silent),
		})
		if err != nil {
			log.Fatalf(fmt.Sprintf(logger.FatalErrorConnectingToDatabase, err))
		}
		instance = &Database{conn: db}
	})
	return instance
}

// GetDB returns the gorm.DB instance
func (d *Database) GetDB() *gorm.DB {
	return d.conn
}
