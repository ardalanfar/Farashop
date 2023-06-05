package store

import (
	"Farashop/internal/adapter/store/model"
	"Farashop/pkg/encrypt"

	"gorm.io/gorm"
)

func InsertSeedAdmin(Db *gorm.DB) error {
	var (
		err  error
		user model.User
	)

	//check admin
	res := Db.Where("username = ?", "admin").Find(&user)
	if res.RowsAffected != 0 {
		return nil
	}

	//create admin system
	pass, _ := encrypt.HashPassword("123456")
	user = model.User{Username: "admin", Email: "admin@yahoo.com", Password: pass, Access: 1, Is_verified: "active"}

	err = Db.Create(&user).Error
	if err != nil {
		return err
	}

	//return
	return nil
}
