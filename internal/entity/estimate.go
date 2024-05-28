package entity

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

// {
// 	"totalPrice": 1,
// 	"estimatedDeliveryTimeInMinutes": 1,
// 	"calculatedEstimateId": "" // save the calculation in the system
// }

type (
	Estimate struct {
		ID           string `json:"id"`
		TotalPrice   int    `json:"total_price"`
		DeliveryTime int    `json:"delivery_time"`
		CreatedAt    string `json:"created_at"`

		UserID string `json:"user_id"`
	}
)
