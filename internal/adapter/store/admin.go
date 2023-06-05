package store

import (
	"Farashop/internal/adapter/store/model"
	"Farashop/internal/dto"
	"Farashop/internal/entity"
	"context"
)

func (conn DbConn) ShowMembers(ctx context.Context) ([]entity.User, error) {
	var users []model.User

	//get all id,email,username,password
	cheek := conn.Db.WithContext(ctx).Select("id", "email", "username", "password").Find(&users)
	if cheek.Error != nil {
		return nil, cheek.Error
	}
	usersEntities := make([]entity.User, len(users))

	for i := range users {
		usersEntities[i] = model.MapToUserEntity(users[i])
	}
	//return
	return usersEntities, nil
}

func (conn DbConn) DeleteMember(ctx context.Context, userID uint) error {
	var user model.User

	//find user for delete
	cheekFind := conn.Db.WithContext(ctx).Where("id = ?", userID).First(&user)
	if cheekFind.Error != nil {
		return cheekFind.Error
	}

	//delete username by id
	cheekDelete := conn.Db.WithContext(ctx).Delete(&user)
	if cheekDelete.Error != nil {
		return cheekDelete.Error
	}
	//return
	return nil
}

func (conn DbConn) ShowInfoMember(ctx context.Context, userID uint) (dto.ShowInfoMember, error) {
	var info dto.ShowInfoMember

	//get "users.email", "users.username", "users.access", "users.is_verified", "products.name", "orders.number", "orders.status", "orders.buy_time" by userID
	cheek := conn.Db.WithContext(ctx).Table("users").Select("users.email", "users.username", "users.access", "users.is_verified", "wallets.balance").Joins("JOIN wallets ON users.id = wallets.user_id").Where("users.id", userID).Find(&info)
	if cheek.Error != nil {
		return dto.ShowInfoMember{}, cheek.Error
	}
	//return
	return info, nil
}
