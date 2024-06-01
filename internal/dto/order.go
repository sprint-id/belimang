package dto

import "github.com/sprint-id/belimang/internal/entity"

type (
	// "calculatedEstimateId": "", // not null
	ReqCreateOrder struct {
		CalculatedEstimateID string `json:"calculatedEstimateId" validate:"required"`
	}

	ParamGetOrderHistory struct {
		MerchantId       string `json:"merchantId"`
		Limit            int    `json:"limit"`
		Offset           int    `json:"offset"`
		Name             string `json:"name"`
		MerchantCategory string `json:"merchantCategory"`
	}

	ResCreateOrder struct {
		OrderID string `json:"orderId"`
	}

	// 	[
	//   {
	//     "orderId": "string",
	// 	  "orders": [
	// 		  {
	// 		    "merchant": {
	// 			    "merchantId":"",
	// 					"name":"",
	// 					"merchantCategory": "",
	// 					"imageUrl": "",
	// 				  "location": {
	// 				    "lat": 1,
	// 				    "long": 1
	// 				  },
	// 				  "createdAt": ""  // should in ISO 8601 format with nanoseconds
	// 				},
	// 		    "items": [
	// 					{
	// 						"itemId":"",
	// 				    "name": "string",
	// 				    "productCategory": ""
	// 				    "price": 1,
	// 	          "quantity": 1,
	// 						"imageUrl": "",
	// 					  "createdAt": ""  // should in ISO 8601 format with nanoseconds
	// 				  }
	// 				]
	// 		  }
	// 	  ]
	//   }
	// ]

	ResGetOrderHistory struct {
		OrderID string     `json:"orderId"`
		Orders  []ResOrder `json:"orders"`
	}

	ResOrder struct {
		Merchant entity.Merchant `json:"merchant"`
		Items    []entity.Item   `json:"items"`
	}
)

func (d *ReqCreateOrder) ToOrderEntity(userId string) entity.Order {
	return entity.Order{
		EstimateID: d.CalculatedEstimateID,
		UserID:     userId,
	}
}
