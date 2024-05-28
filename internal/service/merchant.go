package service

import (
	"context"
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/sprint-id/eniqilo-server/internal/cfg"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
	"github.com/sprint-id/eniqilo-server/internal/repo"
)

type MerchantService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newMerchantService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *MerchantService {
	return &MerchantService{repo, validator, cfg}
}

func (u *MerchantService) CreateMerchant(ctx context.Context, body dto.ReqCreateMerchant, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		fmt.Printf("Error CreateMerchant: %v\n", err)
		return ierr.ErrBadRequest
	}

	// check Image URL if invalid or not complete URL
	if !isValidURL(body.ImageUrl) {
		return ierr.ErrBadRequest
	}

	merchant := body.ToMerchantEntity(sub)
	err = u.repo.Merchant.CreateMerchant(ctx, sub, merchant)
	if err != nil {
		if err == ierr.ErrDuplicate {
			return ierr.ErrDuplicate
		}
		return err
	}

	return nil
}

func (u *MerchantService) GetMerchant(ctx context.Context, param dto.ParamGetMerchant, sub string) ([]dto.ResGetMerchant, error) {
	res, err := u.repo.Merchant.GetMerchant(ctx, param, sub)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func isValidURL(urlString string) bool {
	// url validation using regex
	fmt.Printf("urlString: %s\n", urlString)
	regex := regexp.MustCompile(`^(https?|ftp)://[^/\s]+\.[^/\s]+(?:/.*)?(?:\.[^/\s]+)?$`)
	return regex.MatchString(urlString)
}
