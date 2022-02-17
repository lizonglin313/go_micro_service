package service

import (
	"context"
	"errors"
	"security/model"
)

var (
	ErrUserNotExist = errors.New("username is not exist")
	ErrPassword = errors.New("invalid password")
)

type UserDetailService interface {
	GetUserDetailsByUsername(ctx context.Context, username, password string) (*model.UserDetails, error)
}

type InMemoryUserDetailsService struct {
	userDetailsDict map[string]*model.UserDetails
}

func (service *InMemoryUserDetailsService) GetUserDetailsByUsername(
	ctx context.Context, username, password string) (*model.UserDetails, error) {
	userDetails, ok := service.userDetailsDict[username]
	if ok {
		if userDetails.Password == password {
			return userDetails, nil
		} else {
			return nil, ErrPassword
		}
	} else {
		return nil, ErrUserNotExist
	}
}

// NewInMemoryUserDetailsService 构造方法.
func NewInMemoryUserDetailsService(userDetailsList []*model.UserDetails) *InMemoryUserDetailsService {
	userDetailsDict := make(map[string]*model.UserDetails)

	if userDetailsList != nil {
		for _, value := range userDetailsList {
			userDetailsDict[value.Username] = value
		}
	}

	return &InMemoryUserDetailsService{
		userDetailsDict: userDetailsDict,
	}
}

