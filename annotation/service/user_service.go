package service

import (
	"context"
	"encoding/json"
	"fmt"
	"go101/annotation/model"
)

type UserService interface {
	GetUserByID(ctx context.Context, userID int64) *model.User
}

type UserServiceImpl struct {
}

//@Cache
//@Log
func (u *UserServiceImpl) GetUserByID(ctx context.Context, userID int64) (aUser *model.User) {
	defer func() {
		fmt.Printf("GetUserByID. userID=[%d]\n", userID)
	}()
	user := &model.User{
		ID:   userID,
		Name: fmt.Sprintf("User-%d", userID),
	}
	marshal, _ := json.Marshal(user)
	fmt.Printf("GetUserByID. user=[%s]\n", marshal)
	return user
}

func NewUserService() UserService {
	return &UserServiceImpl{}
}
