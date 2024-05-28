package dto

import (
	"github.com/sprint-id/eniqilo-server/internal/entity"
)

type (
	ReqCreateMerchant struct {
		Name             string   `json:"name" validate:"required,min=2,max=30"`
		MerchantCategory string   `json:"merchantCategory" validate:"required,oneof=SmallRestaurant MediumRestaurant LargeRestaurant MerchandiseRestaurant BoothKiosk ConvenienceStore"`
		ImageUrl         string   `json:"imageUrl" validate:"required,url"`
		Location         Location `json:"location" validate:"required"`
	}

	Location struct {
		Lat  float64 `json:"lat"`
		Long float64 `json:"long"`
	}

	ParamGetMerchant struct {
		MerchantId       string `json:"merchantId"`
		Limit            int    `json:"limit"`
		Offset           int    `json:"offset"`
		Name             string `json:"name"`
		MerchantCategory string `json:"merchantCategory"`
		CreatedAt        string `json:"createdAt"`
	}

	ResCreateMerchant struct {
		MerchantId string `json:"merchantId"`
	}

	ResGetMerchant struct {
		MerchantId       string   `json:"merchantId"`
		Name             string   `json:"name"`
		MerchantCategory string   `json:"merchantCategory"`
		ImageUrl         string   `json:"imageUrl"`
		Location         Location `json:"location"`
		CreatedAt        string   `json:"createdAt"`
	}
)

// ToEntity to convert dto to entity
func (d *ReqCreateMerchant) ToMerchantEntity(userId string) entity.Merchant {
	return entity.Merchant{
		Name:             d.Name,
		MerchantCategory: d.MerchantCategory,
		ImageUrl:         d.ImageUrl,
		Location: entity.Location{
			Lat:  d.Location.Lat,
			Long: d.Location.Long,
		},

		UserID: userId,
	}
}
