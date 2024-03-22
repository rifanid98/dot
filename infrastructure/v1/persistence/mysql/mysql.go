package mysql

import (
	"context"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"dot/config"
	"dot/infrastructure/v1/persistence/mysql/model"
	"dot/pkg/util"

	eLog "log"
)

var log = util.NewLogger()

func New(config *config.MySQLConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Name,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(
			eLog.New(os.Stdout, "\r\n", eLog.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second,   // Slow SQL threshold
				LogLevel:                  logger.Silent, // Log level
				IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
				ParameterizedQueries:      true,          // Don't include params in the SQL log
				Colorful:                  false,         // Disable color
			},
		),
	})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&model.Account{}, &model.Book{})
	if err != nil {
		log.Error(context.Background(), "failed to run auto migrate", err)
		panic("failed to run auto migrate")
	}

	return db
}
