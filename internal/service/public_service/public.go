package public_service

import (
	"Farashop/internal/contract"
	"Farashop/internal/dto"
	"Farashop/internal/entity"
	"Farashop/pkg/encrypt"
	"context"
	"math/rand"
)

type Interactor struct {
	store contract.PublicStore
}

func NewPublic(store contract.PublicStore) Interactor {
	return Interactor{store: store}
}

func (i Interactor) Register(ctx context.Context, req dto.RegisterUserRequest) (dto.RegisterUserResponse, error) {
	var (
		err error
	)
	user := entity.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	//create hash password
	pass, err := encrypt.HashPassword(user.Password)
	if err != nil {
		return dto.RegisterUserResponse{Result: false}, err
	}
	user.Password = pass

	//create verification code
	min := 10000
	max := 99999
	randCode := rand.Intn(max-min) + min
	user.Verification_code = uint(randCode)

	//create user
	create, err := i.store.Register(ctx, user)
	if err != nil || create == false {
		return dto.RegisterUserResponse{Result: false}, err
	}
	//return
	return dto.RegisterUserResponse{Result: true}, nil
}

func (i Interactor) Login(ctx context.Context, req dto.LoginUserRequest) (dto.LoginUserResponse, error) {
	var (
		err error
	)
	user := entity.User{
		Username: req.Username,
		Password: req.Password,
	}

	//get information user by username
	info, err := i.store.Login(ctx, user)
	if err != nil {
		return dto.LoginUserResponse{}, err
	}

	//check password with username
	checkpass := encrypt.CheckPasswordHash(user.Password, info.Password)
	if checkpass != nil {
		return dto.LoginUserResponse{Result: false, User: info}, checkpass
	}
	//return
	return dto.LoginUserResponse{Result: true, User: info}, nil
}

func (i Interactor) MemberValidation(ctx context.Context, req dto.MemberValidationRequest) (dto.MemberValidationResponse, error) {
	var (
		err error
	)
	user := entity.User{
		Username:          req.Username,
		Verification_code: req.Code,
	}

	//member validation
	update, err := i.store.MemberValidation(ctx, user)
	if err != nil || update == false {
		return dto.MemberValidationResponse{Result: false}, err
	}
	//return true
	return dto.MemberValidationResponse{Result: true}, nil
}
