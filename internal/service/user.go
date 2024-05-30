package service

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/sprint-id/belimang/internal/cfg"
	"github.com/sprint-id/belimang/internal/dto"
	"github.com/sprint-id/belimang/internal/ierr"
	"github.com/sprint-id/belimang/internal/repo"
	"github.com/sprint-id/belimang/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newUserService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *UserService {
	return &UserService{repo, validator, cfg}
}

func (u *UserService) RegisterAdmin(ctx context.Context, body dto.ReqAdminRegister) (dto.ResRegisterOrLogin, error) {
	res := dto.ResRegisterOrLogin{}

	err := u.validator.Struct(body)
	if err != nil {
		fmt.Printf("error RegisterAdmin: %v\n", err)
		return res, ierr.ErrBadRequest
	}

	admin := body.ToAdminEntity(u.cfg.BCryptSalt)
	userID, err := u.repo.User.Insert(ctx, admin)
	if err != nil {
		return res, err
	}

	token, _, err := auth.GenerateToken(u.cfg.JWTSecret, 8, auth.JwtPayload{Sub: userID})
	if err != nil {
		return res, err
	}

	res.Token = token

	return res, nil
}

func (u *UserService) RegisterUser(ctx context.Context, body dto.ReqUserRegister) (dto.ResRegisterOrLogin, error) {
	res := dto.ResRegisterOrLogin{}

	err := u.validator.Struct(body)
	if err != nil {
		fmt.Printf("error RegisterNurse: %v\n", err)
		return res, ierr.ErrBadRequest
	}

	user := body.ToUserEntity(u.cfg.BCryptSalt)
	userID, err := u.repo.User.Insert(ctx, user)
	if err != nil {
		return res, err
	}

	token, _, err := auth.GenerateToken(u.cfg.JWTSecret, 8, auth.JwtPayload{Sub: userID})
	if err != nil {
		return res, err
	}

	res.Token = token

	return res, nil
}

func (u *UserService) LoginAdmin(ctx context.Context, body dto.ReqLogin) (dto.ResRegisterOrLogin, error) {
	res := dto.ResRegisterOrLogin{}

	err := u.validator.Struct(body)
	if err != nil {
		fmt.Printf("error LoginAdmin: %v\n", err)
		return res, ierr.ErrBadRequest
	}

	user, err := u.repo.User.GetByUsername(ctx, body.Username)
	if err != nil {
		return res, err
	}

	fmt.Printf("is admin: %v\n", user.IsAdmin)

	// check if user is admin
	if !user.IsAdmin {
		return res, ierr.ErrBadRequest
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return res, ierr.ErrBadRequest
		}
		return res, err
	}

	token, _, err := auth.GenerateToken(u.cfg.JWTSecret, 8, auth.JwtPayload{Sub: user.ID})
	if err != nil {
		return res, err
	}

	res.Token = token

	return res, nil
}

func (u *UserService) LoginUser(ctx context.Context, body dto.ReqLogin) (dto.ResRegisterOrLogin, error) {
	res := dto.ResRegisterOrLogin{}

	err := u.validator.Struct(body)
	if err != nil {
		fmt.Printf("error LoginUser: %v\n", err)
		return res, ierr.ErrBadRequest
	}

	user, err := u.repo.User.GetByUsername(ctx, body.Username)
	if err != nil {
		return res, err
	}

	// check if user is user not admin
	if user.IsAdmin {
		return res, ierr.ErrBadRequest
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return res, ierr.ErrBadRequest
		}
		return res, err
	}

	token, _, err := auth.GenerateToken(u.cfg.JWTSecret, 8, auth.JwtPayload{Sub: user.ID})
	if err != nil {
		return res, err
	}

	res.Token = token

	return res, nil
}
