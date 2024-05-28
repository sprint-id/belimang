package service

import (
	"context"
	"fmt"
	"math"

	"github.com/go-playground/validator/v10"
	"github.com/sprint-id/eniqilo-server/internal/cfg"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
	"github.com/sprint-id/eniqilo-server/internal/repo"
)

type EstimateService struct {
	repo      *repo.Repo
	validator *validator.Validate
	cfg       *cfg.Cfg
}

type pos struct {
	φ float64 // latitude, radians
	ψ float64 // longitude, radians
}

func newEstimateService(repo *repo.Repo, validator *validator.Validate, cfg *cfg.Cfg) *EstimateService {
	return &EstimateService{repo, validator, cfg}
}

func (u *EstimateService) CreateEstimate(ctx context.Context, body dto.ReqCreateEstimate, sub string) error {
	err := u.validator.Struct(body)
	if err != nil {
		fmt.Printf("error CreateEstimate: %v\n", err)
		return ierr.ErrBadRequest
	}

	// calculate total price from items price in orders
	// var totalPrice int
	// for _, order := range body.Orders {
	// 	// get item from order items
	// 	for _, item := range order.Items {
	// 		// get item from item id
	// 		itemEntity, err := u.repo.Item.GetItemByID(ctx, item.ItemID)
	// 		if err != nil {
	// 			return err
	// 		}

	// 		totalPrice += itemEntity[0].Price * item.Quantity
	// 	}
	// }
	// Collect all item IDs
	itemIDs := make([]string, 0, len(body.Orders))
	for _, order := range body.Orders {
		for _, item := range order.Items {
			itemIDs = append(itemIDs, item.ItemID)
		}
	}

	// Fetch all items in a single query
	items, err := u.repo.Item.GetItemsByIDs(ctx, itemIDs)
	if err != nil {
		return err
	}

	// Create a map for item prices
	itemPrices := make(map[string]int)
	for _, item := range items {
		itemPrices[item.ID] = item.Price
	}

	// Calculate total price
	var totalPrice int
	for _, order := range body.Orders {
		for _, item := range order.Items {
			totalPrice += itemPrices[item.ItemID] * item.Quantity
		}
	}

	// calculate distance and estimated time
	var estimatedTime float64
	// user lat and lon
	userLat := body.UserLocation.Lat
	userLong := body.UserLocation.Long
	// merchant start point lat and lon
	var merchantLat, merchantLong float64
	for _, order := range body.Orders {
		if order.IsStartingPoint {
			// get merchant from merchant id
			merchant, err := u.repo.Merchant.GetMerchantByID(ctx, order.MerchantID)
			if err != nil {
				return err
			}

			merchantLat = merchant.Location.Lat
			merchantLong = merchant.Location.Long
		}
	}

	// calculate estimated time in second with velocity 40 km/h
	estimatedTime = 3600 * (hsDist(degPos(userLat, userLong), degPos(merchantLat, merchantLong)) / 40)

	estimate := body.ToEstimateEntity(sub, totalPrice, int(estimatedTime))
	err = u.repo.Estimate.CreateEstimate(ctx, sub, estimate)
	if err != nil {
		if err == ierr.ErrDuplicate {
			return ierr.ErrBadRequest
		}
		return err
	}

	return nil
}

// Reference: https://rosettacode.org/wiki/Haversine_formula#Go
func haversine(θ float64) float64 {
	return .5 * (1 - math.Cos(θ))
}

func degPos(lat, lon float64) pos {
	return pos{lat * math.Pi / 180, lon * math.Pi / 180}
}

const rEarth = 6372.8 // km

func hsDist(p1, p2 pos) float64 {
	return 2 * rEarth * math.Asin(math.Sqrt(haversine(p2.φ-p1.φ)+
		math.Cos(p1.φ)*math.Cos(p2.φ)*haversine(p2.ψ-p1.ψ)))
}
