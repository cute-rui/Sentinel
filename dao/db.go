package dao

import (
	"Sentinel/dao/models"
	"Sentinel/utils/config"
	utils "Sentinel/utils/string"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
)

type SingletonMySQL struct {
	Database *gorm.DB
}

// global call, need to deprecate
var Instance *SingletonMySQL
var lock sync.Once

func setMySQLInstance(config string) *SingletonMySQL {
	lock.Do(func() {
		// Init MySQL Instance
		db, err := gorm.Open(mysql.Open(config), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: false,
			Logger:                                   logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			panic(err)
		}

		Instance = &SingletonMySQL{
			Database: db,
		}
	})
	return Instance
}

func getMysqlDialConfigString() string {
	return utils.StringBuilder(
		config.Conf.GetString(`Database.User`), `:`,
		config.Conf.GetString(`Database.Pass`), `@(`,
		config.Conf.GetString(`Database.Host`), `:`,
		config.Conf.GetString(`Database.Port`), `)/`,
		config.Conf.GetString(`Database.Name`), `?parseTime=true`)
}

func InitDatabase() {
	dialString := getMysqlDialConfigString()

	db := setMySQLInstance(dialString).Database
	err := db.Set(`gorm:table_options`, `ENGINE=InnoDB`).AutoMigrate(
		&models.User{},
		&models.Verify{},
	)

	if err != nil {
		panic(err)
	}

}
