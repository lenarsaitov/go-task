package users

import (
	"context"
	"database/sql"
	"errors"
)

type UserService struct {
	storage *UserStorage
}

func NewUserService(storage *UserStorage) *UserService {
	return &UserService{storage: storage}
}

func (service *UserService) GetUser(c context.Context, userID int) (*UserInfo, error) {
	return service.storage.FindOne(c, userID)
}

func (service *UserService) GetListUsers(c context.Context, params *FilterParams) (*Pagination, error) {
	return service.storage.FindMany(c, params)
}

func (service *UserService) AddUser(c context.Context, params *AddUserRequestParams) (*int64, error) {
	return service.storage.AddUserItem(c, params)
}

func (service *UserService) UpdateUser(c context.Context, params *UpdateUserRequestParams) (bool, error) {
	err := service.storage.UpdateUserItem(c, params)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (service *UserService) DeleteUser(c context.Context, userID int) (bool, error) {
	isExist, err := service.storage.isExistCards(c, userID)
	if err != nil {
		return false, err
	}
	if isExist {
		return false, errors.New("Cant delete user, because exist his cards")
	}

	err = service.storage.DeleteUserItem(c, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
