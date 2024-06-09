package service

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/sprint-id/belimang/internal/cfg"
	"github.com/sprint-id/belimang/internal/dto"
	"github.com/sprint-id/belimang/internal/ierr"
	"github.com/sprint-id/belimang/internal/repo"
	"github.com/sprint-id/belimang/pkg/geo"
	"github.com/sprint-id/belimang/pkg/ids"
)

type EstimateService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

func newEstimateService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *EstimateService {
	return &EstimateService{repo, validator, cfg}
}

func (u *EstimateService) CreateEstimate(ctx context.Context, body dto.ReqCreateEstimate, sub string) (dto.ResCreateEstimate, error) {
	var res dto.ResCreateEstimate
	err := u.validator.Struct(body)
	if err != nil {
		fmt.Printf("error CreateEstimate: %v\n", err)
		return dto.ResCreateEstimate{}, ierr.ErrBadRequest
	}

	// check if user or admin, this is for user only
	isAdmin, err := u.repo.User.IsAdmin(ctx, sub)
	if err != nil {
		return dto.ResCreateEstimate{}, ierr.ErrInternal
	}

	if isAdmin {
		return dto.ResCreateEstimate{}, ierr.ErrForbidden
	}

	// check just one starting point
	var startingPointCount int
	for _, order := range body.Orders {
		if *order.IsStartingPoint {
			startingPointCount++
		}
	}

	if startingPointCount != 1 {
		return dto.ResCreateEstimate{}, ierr.ErrBadRequest
	}

	// check merchant id and order id is valid form as uuid
	for _, order := range body.Orders {
		if !ids.ValidateUUID(order.MerchantID) {
			return dto.ResCreateEstimate{}, ierr.ErrNotFound
		}
		// check merchant id is valid in database
		_, err := u.repo.Merchant.GetMerchantByID(ctx, order.MerchantID)
		if err != nil {
			return dto.ResCreateEstimate{}, ierr.ErrNotFound
		}
	}

	// check item id is valid
	for _, order := range body.Orders {
		for _, item := range order.Items {
			if !ids.ValidateUUID(item.ItemID) {
				return dto.ResCreateEstimate{}, ierr.ErrNotFound
			}
			// check item id is valid in database
			_, err := u.repo.Item.GetItemByID(ctx, item.ItemID)
			if err != nil {
				return dto.ResCreateEstimate{}, ierr.ErrNotFound
			}
		}
	}

	// calculate total price from items price in orders
	var totalPrice int
	for _, order := range body.Orders {
		// get item from order items
		for _, item := range order.Items {
			// get item from item id
			itemEntity, err := u.repo.Item.GetItemByID(ctx, item.ItemID)
			if err != nil {
				return dto.ResCreateEstimate{}, ierr.ErrInternal
			}

			totalPrice += itemEntity.Price * item.Quantity
		}
	}

	fmt.Printf("totalPrice: %v\n", totalPrice)
	// // Collect all item IDs
	// itemIDs := make([]string, 0, len(body.Orders))
	// for _, order := range body.Orders {
	// 	for _, item := range order.Items {
	// 		itemIDs = append(itemIDs, item.ItemID)
	// 	}
	// }

	// // Fetch all items in a single query
	// items, err := u.repo.Item.GetItemsByIDs(ctx, itemIDs)
	// if err != nil {
	// 	return err
	// }

	// // Create a map for item prices
	// itemPrices := make(map[string]int)
	// for _, item := range items {
	// 	itemPrices[item.ID] = item.Price
	// }

	// // Calculate total price
	// var totalPrice int
	// for _, order := range body.Orders {
	// 	for _, item := range order.Items {
	// 		totalPrice += itemPrices[item.ItemID] * item.Quantity
	// 	}
	// }

	// calculate distance and estimated time
	var estimatedTime float64
	// user lat and lon
	userLat := body.UserLocation.Lat
	userLong := body.UserLocation.Long
	// merchant start point lat and lon
	var merchantLat, merchantLong float64
	for _, order := range body.Orders {
		if *order.IsStartingPoint {
			// get merchant from merchant id
			merchant, err := u.repo.Merchant.GetMerchantByID(ctx, order.MerchantID)
			if err != nil {
				return dto.ResCreateEstimate{}, ierr.ErrInternal
			}

			merchantLat = merchant.Location.Lat
			merchantLong = merchant.Location.Long
		}
	}

	// calculate distance between user and merchant
	distance := geo.CalculateDistance(userLat, userLong, merchantLat, merchantLong)
	fmt.Printf("original distance: %v\n", distance)
	// print user and merchant location
	fmt.Printf("userLat: %v, userLong: %v\n", userLat, userLong)
	fmt.Printf("merchantLat: %v, merchantLong: %v\n", merchantLat, merchantLong)
	// if distance is > 3 km, bad request
	if distance > 3 {
		return dto.ResCreateEstimate{}, ierr.ErrBadRequest
	}
	// calculate estimated time in minutes with velocity 40 km/h
	estimatedTime = 60 * geo.CalculateDistance(userLat, userLong, merchantLat, merchantLong) / 40
	fmt.Printf("original estimatedTimeInMinutes: %v\n", estimatedTime)

	estimate := body.ToEstimateEntity(sub, totalPrice, int(estimatedTime))
	res, err = u.repo.Estimate.CreateEstimate(ctx, sub, estimate)
	if err != nil {
		if err == ierr.ErrDuplicate {
			return dto.ResCreateEstimate{}, ierr.ErrDuplicate
		}
		return dto.ResCreateEstimate{}, ierr.ErrInternal
	}

	return res, nil
}
