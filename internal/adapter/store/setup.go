package store

import (
	"Farashop/internal/adapter/store/model"
	"Farashop/internal/config"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//connect to the postgresql

type DbConn struct {
	Db *gorm.DB
}

func New() DbConn {
	var (
		err      error
		database *gorm.DB
	)

	configSys := config.GetConfig()

	//connection strings
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		configSys.DB.Host,
		configSys.DB.Username,
		configSys.DB.Password,
		configSys.DB.Dbname,
	)

	//connect database
	database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	//migrate model
	err = database.AutoMigrate(&model.User{}, &model.Product{}, &model.Order{}, &model.Wallet{})
	if err != nil {
		panic("Failed to auto migrate database!")
	}

	//create information default
	/*-----------------------------------------------------*/
	err = InsertSeedAdmin(database)
	if err != nil {
		panic("Failed to Insert Admin")
	}
	/*-----------------------------------------------------*/

	//return connection database
	return DbConn{Db: database}
}
