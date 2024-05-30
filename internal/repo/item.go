package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sprint-id/belimang/internal/dto"
	"github.com/sprint-id/belimang/internal/entity"
	timepkg "github.com/sprint-id/belimang/pkg/time"
)

type itemRepo struct {
	conn *pgxpool.Pool
}

func newItemRepo(conn *pgxpool.Pool) *itemRepo {
	return &itemRepo{conn}
}

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

func (cr *itemRepo) AddItem(ctx context.Context, sub, merchantId string, item entity.Item) (dto.ResAddItem, error) {
	// add item
	q := `INSERT INTO items (id, user_id, merchant_id, name, product_category, price, image_url, created_at)
	VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, EXTRACT(EPOCH FROM now())::bigint) RETURNING id`

	var id string
	err := cr.conn.QueryRow(ctx, q, sub, merchantId, item.Name, item.ProductCategory, item.Price, item.ImageUrl).Scan(&id)
	if err != nil {
		fmt.Printf("error query: %v\n", err)
		return dto.ResAddItem{}, err
	}

	return dto.ResAddItem{ItemId: id}, nil
}

func (cr *itemRepo) GetItem(ctx context.Context, param dto.ParamGetItem, sub string) ([]dto.ResGetItem, error) {
	var query strings.Builder

	query.WriteString("SELECT id, name, product_category, price, image_url, created_at FROM items WHERE 1=1 ")

	// param id
	if param.ItemId != "" {
		id, err := strconv.Atoi(param.ItemId)
		if err != nil {
			return nil, err
		}
		query.WriteString(fmt.Sprintf("AND id = %d ", id))
	}

	// it should search by wildcard (ex: if search by name=een then user with name kayleen should appear)
	if param.Name != "" {
		query.WriteString(fmt.Sprintf("AND nip LIKE '%s' ", fmt.Sprintf("%%%s%%", param.Name)))
	}

	// productCategory filter based on category
	if param.ProductCategory != "" {
		query.WriteString(fmt.Sprintf("AND product_category = '%s' ", param.ProductCategory))
	}

	// param createdAt sort by created time asc or desc, if value is wrong, just ignore the param
	if param.CreatedAt == "asc" && param.Offset == 0 {
		query.WriteString("ORDER BY created_at ASC ")
	} else if param.CreatedAt == "desc" && param.Offset == 0 {
		query.WriteString("ORDER BY created_at DESC ")
	} else if param.Offset == 0 {
		query.WriteString("ORDER BY created_at DESC ")
	}

	// limit and offset
	if param.Limit == 0 {
		param.Limit = 5
	}

	query.WriteString(fmt.Sprintf("LIMIT %d OFFSET %d", param.Limit, param.Offset))

	// show query
	fmt.Println(query.String())

	rows, err := cr.conn.Query(ctx, query.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]dto.ResGetItem, 0, 10)
	for rows.Next() {
		var imageUrl sql.NullString
		var createdAt int64

		result := dto.ResGetItem{}
		err := rows.Scan(
			&result.ItemId,
			&result.Name,
			&result.ProductCategory,
			&result.Price,
			&imageUrl,
			&createdAt)
		if err != nil {
			return nil, err
		}

		result.CreatedAt = timepkg.TimeToISO8601(time.Unix(createdAt, 0))
		results = append(results, result)
	}

	return results, nil
}

func (cr *itemRepo) GetItemByID(ctx context.Context, id string) (entity.Item, error) {
	q := `SELECT id, user_id, name, product_category, price, image_url, created_at FROM items WHERE id = $1`

	var result entity.Item
	var createdAt int64
	err := cr.conn.QueryRow(ctx, q, id).Scan(
		&result.ID,
		&result.UserID,
		&result.Name,
		&result.ProductCategory,
		&result.Price,
		&result.ImageUrl,
		&createdAt)
	if err != nil {
		return entity.Item{}, err
	}

	result.CreatedAt = timepkg.TimeToISO8601(time.Unix(createdAt, 0))

	return result, nil
}

func (cr *itemRepo) GetItemsByIDs(ctx context.Context, ids []string) ([]entity.Item, error) {
	fmt.Printf("ids: %v\n", ids)
	// array to string like this (id_1, id_2)
	strIds := arrayToString(ids)

	q := `SELECT id, user_id, name, product_category, price, image_url, created_at FROM items WHERE id IN ($1)`

	rows, err := cr.conn.Query(ctx, q, strIds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]entity.Item, 0, 10)
	for rows.Next() {
		result := entity.Item{}
		var createdAt int64
		err := rows.Scan(
			&result.ID,
			&result.UserID,
			&result.Name,
			&result.ProductCategory,
			&result.Price,
			&result.ImageUrl,
			&createdAt)
		if err != nil {
			return nil, err
		}

		result.CreatedAt = timepkg.TimeToISO8601(time.Unix(createdAt, 0))
		results = append(results, result)
	}

	return results, nil
}

func arrayToString(arr []string) string {
	// Create a slice of strings to hold the string representations of the integers
	strArr := make([]string, len(arr))
	copy(strArr, arr) // Fix: Use copy instead of loop
	// Join the string slice with commas and wrap with parentheses
	return fmt.Sprintf("(%s)", strings.Join(strArr, ","))
}
