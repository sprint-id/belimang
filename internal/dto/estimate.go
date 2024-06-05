package dto

import "github.com/sprint-id/belimang/internal/entity"

// {
// 	"userLocation": {
// 	  "lat": 1, // not null | float
// 	  "long": 1  // not null | float
// 	},
// 	"orders": [
// 	  {
// 		"merchantId": "string", // not null
// 		"isStartingPoint" : true
// 		// ⬆️ not null | there's should be one isStartingPoint == true in orders array
// 		// | if none are true, or true > 1 items, it's not valid
// 		"items": [
// 		  {
// 			"itemId": "string", // not null
// 			"quantity": 1 // not null
// 		  }
// 		]
// 	  }
// 	]
//   }

type (
	ReqCreateEstimate struct {
		UserLocation ReqUserLocation `json:"userLocation" validate:"required"`
		Orders       []ReqOrder      `json:"orders" validate:"required,dive,required"`
	}

	ReqUserLocation struct {
		Lat  float64 `json:"lat" validate:"required"`
		Long float64 `json:"long" validate:"required"`
	}

	ReqOrder struct {
		MerchantID      string    `json:"merchantId" validate:"required"`
		IsStartingPoint *bool     `json:"isStartingPoint" validate:"required"`
		Items           []ReqItem `json:"items" validate:"required,dive,required"`
	}

	ReqItem struct {
		ItemID   string `json:"itemId" validate:"required"`
		Quantity int    `json:"quantity" validate:"required"`
	}

	// {
	// 	"totalPrice": 1,
	// 	"estimatedDeliveryTimeInMinutes": 1,
	// 	"calculatedEstimateId": "" // save the calculation in the system
	// }
	ResCreateEstimate struct {
		TotalPrice                     int    `json:"totalPrice"`
		EstimatedDeliveryTimeInMinutes int    `json:"estimatedDeliveryTimeInMinutes"`
		CalculatedEstimateID           string `json:"calculatedEstimateId"`
	}
)

func (d *ReqCreateEstimate) ToEstimateEntity(userId string, totalPrice int, estimatedTime int) entity.Estimate {
	return entity.Estimate{
		TotalPrice:   totalPrice,
		DeliveryTime: estimatedTime,
		UserLocation: entity.UserLocation{
			Lat:  d.UserLocation.Lat,
			Long: d.UserLocation.Long,
		},
		Orders: d.toOrderDetailEntity(),
		UserID: userId,
	}
}

func (d *ReqCreateEstimate) toOrderDetailEntity() []entity.OrderDetail {
	var orders []entity.OrderDetail
	for _, order := range d.Orders {
		orders = append(orders, entity.OrderDetail{
			MerchantID:      order.MerchantID,
			IsStartingPoint: *order.IsStartingPoint,
			Items:           order.toItemDetailEntity(),
		})
	}
	return orders
}

func (d *ReqOrder) toItemDetailEntity() []entity.ItemDetail {
	var items []entity.ItemDetail
	for _, item := range d.Items {
		items = append(items, entity.ItemDetail{
			ItemID:   item.ItemID,
			Quantity: item.Quantity,
		})
	}
	return items
}
