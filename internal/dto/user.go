package dto

import (
	"github.com/sprint-id/belimang/internal/entity"
	"github.com/sprint-id/belimang/pkg/auth"
)

type (
	ReqAdminRegister struct {
		Username string `json:"username" validate:"required,min=5,max=30"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=5,max=30"`
	}

	ReqUserRegister struct {
		Username string `json:"username" validate:"required,min=5,max=30"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=5,max=30"`
	}

	ReqLogin struct {
		Username string `json:"username" validate:"required,min=5,max=30"`
		Password string `json:"password" validate:"required,min=5,max=30"`
	}

	ResRegisterOrLogin struct {
		Token string `json:"token"`
	}
)

func (d *ReqAdminRegister) ToAdminEntity(cryptCost int) entity.User {
	return entity.User{Username: d.Username, Password: auth.HashPassword(d.Password, cryptCost), Email: d.Email, IsAdmin: true}
}

func (d *ReqUserRegister) ToUserEntity(cryptCost int) entity.User {
	return entity.User{Username: d.Username, Password: auth.HashPassword(d.Password, cryptCost), Email: d.Email, IsAdmin: false}
}
