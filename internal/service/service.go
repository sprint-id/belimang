package service

import (
	"github.com/go-playground/validator/v10"

	"github.com/sprint-id/eniqilo-server/internal/cfg"
	"github.com/sprint-id/eniqilo-server/internal/repo"
)

type Service struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg

	User     *UserService
	Merchant *MerchantService
	Item     *ItemService
	Estimate *EstimateService
}

func NewService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *Service {
	service := Service{}
	service.repo = repo
	service.validator = validator
	service.cfg = cfg

	service.User = newUserService(repo, validator, cfg)
	service.Item = newItemService(repo, validator, cfg)
	service.Merchant = newMerchantService(repo, validator, cfg)
	service.Estimate = newEstimateService(repo, validator, cfg)

	return &service
}
