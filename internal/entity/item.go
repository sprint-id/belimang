package entity

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

type Item struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	ProductCategory string `json:"productCategory"`
	Price           int    `json:"price"`
	ImageUrl        string `json:"imageUrl"`
	CreatedAt       string `json:"created_at"`

	UserID string `json:"user_id"`
}
