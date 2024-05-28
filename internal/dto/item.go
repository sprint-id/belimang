package dto

import (
	"github.com/sprint-id/eniqilo-server/internal/entity"
)

// {
// 	"name": "string", // not null | minLength 2 | maxLength 30
// 	"productCategory": "" /** enum of:
// 	- `Beverage`
// 	  - `Food`
// 	  - `Snack`
// 	  - `Condiments`
// 	  - `Additions`
// 	  */
// 	"price": 1, // not null | min 1
// 	  "imageUrl": "" // not null | should be image url
//   }

type (
	ReqAddItem struct {
		Name            string `json:"name" validate:"required,min=2,max=30"`
		ProductCategory string `json:"productCategory" validate:"required,oneof=Beverage Food Snack Condiments Additions"`
		Price           int    `json:"price" validate:"required,min=1"`
		ImageUrl        string `json:"imageUrl" validate:"required,url"`
	}

	ParamGetItem struct {
		ItemId          string `json:"itemId"`
		Limit           int    `json:"limit"`
		Offset          int    `json:"offset"`
		Name            string `json:"name"`
		ProductCategory string `json:"productCategory"`
		CreatedAt       string `json:"createdAt"`
	}

	ResAddItem struct {
		ItemId string `json:"itemId"`
	}

	ResGetItem struct {
		ItemId          string `json:"itemId"`
		Name            string `json:"name"`
		ProductCategory string `json:"productCategory"`
		Price           int    `json:"price"`
		ImageUrl        string `json:"imageUrl"`
		CreatedAt       string `json:"createdAt"`
	}
)

func (d *ReqAddItem) ToItemEntity(userId string) entity.Item {
	return entity.Item{
		Name:            d.Name,
		ProductCategory: d.ProductCategory,
		Price:           d.Price,
		ImageUrl:        d.ImageUrl,
		UserID:          userId,
	}
}
