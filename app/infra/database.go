package infra

import (
	"fmt"
	"log"
	"os"
	"sample-service/app/model"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func DbConnect() error {

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DBHOST"),
		os.Getenv("DBPORT"),
		os.Getenv("DBUSER"),
		os.Getenv("DBPASS"),
		os.Getenv("DBNAME"),
	)

	//println(dsn)

	var err error

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
			Colorful:                  true,
		},
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: newLogger})

	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(&model.Package{}, &model.Area{}, &model.Esim{}, &model.EsimAssignedPlan{})
	DropUnusedColumns(&model.Package{}, &model.Area{}, &model.Esim{}, &model.EsimAssignedPlan{})

	return nil
}

func DropUnusedColumns(dsts ...interface{}) {

	for _, dst := range dsts {
		stmt := &gorm.Statement{DB: DB}
		stmt.Parse(dst)
		fields := stmt.Schema.Fields
		columns, _ := DB.Debug().Migrator().ColumnTypes(dst)

		for i := range columns {
			found := false
			for j := range fields {
				if columns[i].Name() == fields[j].DBName {
					found = true
					break
				}
			}
			if !found {
				fmt.Printf("Going to DROP %s column\n", columns[i].Name())
				DB.Migrator().DropColumn(dst, columns[i].Name())
			}
		}
	}
}
