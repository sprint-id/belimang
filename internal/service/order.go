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

type OrderService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newOrderService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *OrderService {
	return &OrderService{repo, validator, cfg}
}

func (u *OrderService) CreateOrder(ctx context.Context, body dto.ReqCreateOrder, sub string) (dto.ResCreateOrder, error) {
	var res dto.ResCreateOrder
	err := u.validator.Struct(body)
	if err != nil {
		fmt.Printf("error CreateOrder: %v\n", err)
		return dto.ResCreateOrder{}, ierr.ErrBadRequest
	}

	// check if user or admin, this is for user only
	isAdmin, err := u.repo.User.IsAdmin(ctx, sub)
	if err != nil {
		return dto.ResCreateOrder{}, ierr.ErrInternal
	}

	if isAdmin {
		return dto.ResCreateOrder{}, ierr.ErrForbidden
	}

	order := body.ToOrderEntity(sub)
	res, err = u.repo.Order.CreateOrder(ctx, sub, order)
	if err != nil {
		if err == ierr.ErrDuplicate {
			return dto.ResCreateOrder{}, ierr.ErrDuplicate
		}
		return dto.ResCreateOrder{}, ierr.ErrInternal
	}

	return res, nil
}

func (u *OrderService) GetOrderHistory(ctx context.Context, param dto.ParamGetOrderHistory, sub string) ([]dto.ResGetOrderHistory, error) {
	var res []dto.ResGetOrderHistory

	// check if user or admin, this is for user only
	isAdmin, err := u.repo.User.IsAdmin(ctx, sub)
	if err != nil {
		return nil, ierr.ErrInternal
	}

	if isAdmin {
		return nil, ierr.ErrForbidden
	}

	res, err = u.repo.Order.GetOrderHistory(ctx, param, sub)
	if err != nil {
		return nil, ierr.ErrInternal
	}

	return res, nil
}
