package service

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/sprint-id/belimang/internal/cfg"
	"github.com/sprint-id/belimang/internal/dto"
	"github.com/sprint-id/belimang/internal/ierr"
	"github.com/sprint-id/belimang/internal/repo"
)

type ItemService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newItemService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *ItemService {
	return &ItemService{repo, validator, cfg}
}

// {
// 	"identityNumber": 123123, // not null, should be 16 digit
// 	"symptoms": "", // not null, minLength 1, maxLength 2000,
// 	"medications" : "" // not null, minLength 1, maxLength 2000
// }

func (u *ItemService) AddItem(ctx context.Context, body dto.ReqAddItem, sub, merchantId string) (dto.ResAddItem, error) {
	var res dto.ResAddItem
	err := u.validator.Struct(body)
	if err != nil {
		fmt.Printf("error AddItem: %v\n", err)
		return dto.ResAddItem{}, ierr.ErrBadRequest
	}

	// check if user or admin, this is for admin only
	isAdmin, err := u.repo.User.IsAdmin(ctx, sub)
	if err != nil {
		return dto.ResAddItem{}, ierr.ErrInternal
	}

	if !isAdmin {
		return dto.ResAddItem{}, ierr.ErrForbidden
	}

	// validate image url
	// check Image URL if invalid or not complete URL
	if !isValidURL(body.ImageUrl) {
		return dto.ResAddItem{}, ierr.ErrInvalidURL
	}

	item := body.ToItemEntity(sub)
	res, err = u.repo.Item.AddItem(ctx, sub, merchantId, item)
	if err != nil {
		if err == ierr.ErrDuplicate {
			return dto.ResAddItem{}, ierr.ErrDuplicate
		}
		return dto.ResAddItem{}, ierr.ErrInternal
	}

	return res, nil
}

func (u *ItemService) GetItem(ctx context.Context, param dto.ParamGetItem, sub string) ([]dto.ResGetItem, error) {

	err := u.validator.Struct(param)
	if err != nil {
		fmt.Printf("error GetRecord: %v\n", err)
		return nil, ierr.ErrBadRequest
	}

	// check if user or admin, this is for admin only
	isAdmin, err := u.repo.User.IsAdmin(ctx, sub)
	if err != nil {
		return nil, ierr.ErrInternal
	}

	if !isAdmin {
		return nil, ierr.ErrForbidden
	}

	res, err := u.repo.Item.GetItem(ctx, param, sub)
	if err != nil {
		return nil, err
	}

	return res, nil
}
