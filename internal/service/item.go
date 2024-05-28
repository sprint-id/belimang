package service

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/sprint-id/eniqilo-server/internal/cfg"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
	"github.com/sprint-id/eniqilo-server/internal/repo"
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

func (u *ItemService) AddItem(ctx context.Context, body dto.ReqAddItem, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		fmt.Printf("error AddItem: %v\n", err)
		return ierr.ErrBadRequest
	}

	item := body.ToItemEntity(sub)
	err = u.repo.Item.AddItem(ctx, sub, item)
	if err != nil {
		if err == ierr.ErrDuplicate {
			return ierr.ErrBadRequest
		}
		return err
	}

	return nil
}

func (u *ItemService) GetItem(ctx context.Context, param dto.ParamGetItem, sub string) ([]dto.ResGetItem, error) {

	err := u.validator.Struct(param)
	if err != nil {
		fmt.Printf("error GetRecord: %v\n", err)
		return nil, ierr.ErrBadRequest
	}

	res, err := u.repo.Item.GetItem(ctx, param, sub)
	if err != nil {
		return nil, err
	}

	return res, nil
}
