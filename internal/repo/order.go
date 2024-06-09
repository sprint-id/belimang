package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sprint-id/belimang/internal/dto"
	"github.com/sprint-id/belimang/internal/entity"
	timepkg "github.com/sprint-id/belimang/pkg/time"
)

type orderRepo struct {
	conn *pgxpool.Pool
}

func newOrderRepo(conn *pgxpool.Pool) *orderRepo {
	return &orderRepo{conn}
}

func (cr *orderRepo) CreateOrder(ctx context.Context, sub string, order entity.Order) (dto.ResCreateOrder, error) {
	// add order
	q := `INSERT INTO orders (id, user_id, estimate_id, created_at)
	VALUES ( gen_random_uuid(), $1, $2, EXTRACT(EPOCH FROM now())::bigint) RETURNING id`

	var id string
	err := cr.conn.QueryRow(ctx, q, sub, order.EstimateID).Scan(&id)
	if err != nil {
		fmt.Printf("error query: %v\n", err)
		return dto.ResCreateOrder{}, err
	}

	return dto.ResCreateOrder{OrderID: id}, nil
}

// 	[
//   {
//     "orderId": "string",
// 	   "orders": [
// 		  {
// 		    "merchant": {
// 			    "merchantId":"",
// 				"name":"",
// 				"merchantCategory": "",
// 				"imageUrl": "",
// 				"location": {
// 				    "lat": 1,
// 				    "long": 1
// 				},
// 				"createdAt": ""  // should in ISO 8601 format with nanoseconds
// 			},
// 		    "items": [
// 				{
// 					"itemId":"",
// 				    "name": "string",
// 				    "productCategory": ""
// 				    "price": 1,
// 	          		"quantity": 1,
// 					"imageUrl": "",
// 					"createdAt": ""  // should in ISO 8601 format with nanoseconds
// 				 }
// 			]
// 		  }
// 	  ]
//   }
// ]

func (cr *orderRepo) GetOrderHistory(ctx context.Context, param dto.ParamGetOrderHistory, sub string) ([]dto.ResGetOrderHistory, error) {
	q := `
	SELECT 
	o.id,
	m.id, m.name, m.merchant_category, m.image_url, m.location_lat, m.location_long, m.created_at,
	i.id, i.name, i.product_category, i.price, eoi.quantity, i.image_url, i.created_at
	FROM orders o 
	LEFT JOIN estimates e ON o.estimate_id = e.id
	LEFT JOIN estimate_orders eo ON e.id = eo.estimate_id
	LEFT JOIN merchants m ON eo.merchant_id = m.id
	LEFT JOIN estimate_order_items eoi ON eo.id = eoi.estimate_order_id
	LEFT JOIN items i ON eoi.item_id = i.id
	WHERE o.user_id = $1
	`

	rows, err := cr.conn.Query(ctx, q, sub)
	if err != nil {
		fmt.Printf("error query: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var orders []dto.ResGetOrderHistory
	var currentOrderID string
	var currentMerchantID string
	var currentOrder dto.ResGetOrderHistory
	var currentMerchant dto.ResOrder

	for rows.Next() {
		var orderID string
		var merchantID string
		var merchantCreatedAt int64
		var itemCreatedAt int64

		// var order dto.ResGetOrderHistory
		var merchant entity.Merchant
		var item dto.ResItem

		err := rows.Scan(
			&orderID,
			&merchantID,
			&merchant.Name,
			&merchant.MerchantCategory,
			&merchant.ImageUrl,
			&merchant.Location.Lat,
			&merchant.Location.Long,
			&merchantCreatedAt,
			&item.ItemID,
			&item.Name,
			&item.ProductCategory,
			&item.Price,
			&item.Quantity,
			&item.ImageURL,
			&itemCreatedAt,
		)
		if err != nil {
			fmt.Printf("error scanning rows: %v\n", err)
			return nil, err
		}

		if currentOrderID != orderID {
			if currentOrderID != "" {
				orders = append(orders, currentOrder)
			}

			currentOrderID = orderID
			currentOrder = dto.ResGetOrderHistory{
				OrderID: orderID,
				Orders:  []dto.ResOrder{},
			}
		}

		if currentMerchantID != merchantID {
			if currentMerchantID != "" {
				currentOrder.Orders = append(currentOrder.Orders, currentMerchant)
			}

			currentMerchantID = merchantID
			merchant.ID = currentMerchantID
			merchant.CreatedAt = timepkg.TimeToISO8601(time.Unix(merchantCreatedAt, 0))
			currentMerchant = dto.ResOrder{
				Merchant: merchant,
				Items:    []dto.ResItem{},
			}
		}

		item.CreatedAt = timepkg.TimeToISO8601(time.Unix(itemCreatedAt, 0))
		currentMerchant.Items = append(currentMerchant.Items, item)
	}

	if currentMerchantID != "" {
		currentOrder.Orders = append(currentOrder.Orders, currentMerchant)
	}

	if currentOrderID != "" {
		orders = append(orders, currentOrder)
	}

	if err := rows.Err(); err != nil {
		fmt.Printf("error iterating over rows: %v\n", err)
		return nil, err
	}

	return orders, nil
}
