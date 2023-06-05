package store

//Account Store intractor

import (
	"Farashop/internal/adapter/store/model"

	"Farashop/internal/entity"
	"context"
)

func (conn DbConn) Register(ctx context.Context, user entity.User) (bool, error) {
	u := model.MapFromUserEntity(user)

	//cheek username and email
	Cheek := conn.Db.WithContext(ctx).Select("id").Where("username = ? OR email = ?", u.Username, u.Email).First(&u)
	if Cheek.Error != nil && Cheek.RowsAffected != 0 && u.ID != 0 {
		return false, Cheek.Error
	}
	//create user
	Create := conn.Db.WithContext(ctx).Create(&u)
	if Create.Error != nil {
		return false, Create.Error
	}
	//return
	return true, nil
}

func (conn DbConn) Login(ctx context.Context, user entity.User) (entity.User, error) {
	u := model.MapFromUserEntity(user)

	//get "id", "email", "password", "access", "username" by username
	Cheek := conn.Db.WithContext(ctx).Select("id", "email", "password", "access", "username").Where("username = ?", u.Username).First(&u)
	if Cheek.Error != nil {
		return entity.User{}, Cheek.Error
	}
	//return
	return model.MapToUserEntity(u), nil
}

func (conn DbConn) MemberValidation(ctx context.Context, user entity.User) (bool, error) {
	u := model.MapFromUserEntity(user)

	//update verify code
	Cheek := conn.Db.WithContext(ctx).Model(&u).Where("username = ? AND verification_code = ?", u.Username, u.Verification_code).Update("is_verified", "active")
	if Cheek.RowsAffected == 0 || Cheek.Error != nil {
		return false, Cheek.Error
	}
	//return
	return true, nil
}
