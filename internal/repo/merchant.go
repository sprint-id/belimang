package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sprint-id/belimang/internal/dto"
	"github.com/sprint-id/belimang/internal/entity"
	"github.com/sprint-id/belimang/internal/ierr"
	timepkg "github.com/sprint-id/belimang/pkg/time"
)

type merchantRepo struct {
	conn *pgxpool.Pool
}

func newMerchantRepo(conn *pgxpool.Pool) *merchantRepo {
	return &merchantRepo{conn}
}

func (mr *merchantRepo) CreateMerchant(ctx context.Context, sub string, merchant entity.Merchant) (dto.ResCreateMerchant, error) {
	// Start a transaction with serializable isolation level
	tx, err := mr.conn.Begin(ctx)
	if err != nil {
		return dto.ResCreateMerchant{}, err
	}
	defer tx.Rollback(ctx)

	q := `INSERT INTO merchants (id, user_id, name, merchant_category, image_url, location_lat, location_long, created_at)
	VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, EXTRACT(EPOCH FROM now())::bigint) RETURNING id`
	fmt.Printf("sub: %s\n", sub)

	var id string
	err = tx.QueryRow(ctx, q, sub,
		merchant.Name,
		merchant.MerchantCategory,
		merchant.ImageUrl,
		merchant.Location.Lat,
		merchant.Location.Long).Scan(&id)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				return dto.ResCreateMerchant{}, ierr.ErrDuplicate
			}
		}
		return dto.ResCreateMerchant{}, err
	}

	// Commit the transaction
	err = tx.Commit(ctx)
	if err != nil {
		return dto.ResCreateMerchant{}, err
	}

	return dto.ResCreateMerchant{MerchantId: id}, nil
}

func (mr *merchantRepo) GetMerchant(ctx context.Context, param dto.ParamGetMerchant, sub string) ([]dto.ResGetMerchant, error) {
	var query strings.Builder

	query.WriteString("SELECT id, name, merchant_category, image_url, location_lat, location_long, created_at FROM merchants WHERE 1=1 ")

	if param.Name != "" {
		query.WriteString(fmt.Sprintf("AND LOWER(name) LIKE LOWER('%s') ", fmt.Sprintf("%%%s%%", param.Name)))
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

	rows, err := mr.conn.Query(ctx, query.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var merchants []dto.ResGetMerchant
	for rows.Next() {
		var merchant dto.ResGetMerchant
		var createdAt int64
		err = rows.Scan(&merchant.MerchantId, &merchant.Name, &merchant.MerchantCategory, &merchant.ImageUrl, &merchant.Location.Lat, &merchant.Location.Long, &createdAt)
		if err != nil {
			return nil, err
		}

		merchant.CreatedAt = timepkg.TimeToISO8601(time.Unix(createdAt, 0))
		merchants = append(merchants, merchant)
	}

	return merchants, nil
}

func (mr *merchantRepo) GetMerchantByID(ctx context.Context, id string) (entity.Merchant, error) {
	q := `SELECT id, name, merchant_category, image_url, location_lat, location_long, created_at FROM merchants WHERE id = $1`

	var merchant entity.Merchant
	var createdAt int64
	err := mr.conn.QueryRow(ctx, q, id).Scan(&merchant.ID, &merchant.Name, &merchant.MerchantCategory, &merchant.ImageUrl, &merchant.Location.Lat, &merchant.Location.Long, &createdAt)
	if err != nil {
		return entity.Merchant{}, err
	}

	merchant.CreatedAt = timepkg.TimeToISO8601(time.Unix(createdAt, 0))

	return merchant, nil
}

func (mr *merchantRepo) GetNearbyMerchant(ctx context.Context, param dto.ParamGetNearbyMerchant, sub string, lat, long float64) ([]dto.ResGetNearbyMerchant, error) {
	var query strings.Builder

	query.WriteString("SELECT id, name, merchant_category, image_url, location_lat, location_long, " +
		"ST_Distance(" +
		"ST_SetSRID(ST_MakePoint(location_long, location_lat), 4326)::geography, " +
		"ST_SetSRID(ST_MakePoint($1, $2), 4326)::geography " +
		") AS distance, " +
		"created_at FROM merchants WHERE 1=1 ")

	if param.Name != "" {
		query.WriteString(fmt.Sprintf("AND LOWER(name) LIKE LOWER('%s') ", fmt.Sprintf("%%%s%%", param.Name)))
	}

	// merchant category
	if param.MerchantCategory != "" {
		query.WriteString(fmt.Sprintf("AND LOWER(merchant_category) LIKE LOWER('%s') ", fmt.Sprintf("%%%s%%", param.MerchantCategory)))
	}

	// order by distance
	query.WriteString("ORDER BY distance")

	// limit and offset
	if param.Limit == 0 {
		param.Limit = 5
	}

	query.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", param.Limit, param.Offset))

	// show query
	fmt.Println(query.String())

	rows, err := mr.conn.Query(ctx, query.String(), long, lat)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var merchants []dto.ResGetNearbyMerchant
	for rows.Next() {
		var merchant dto.ResGetNearbyMerchant
		var createdAt int64
		err = rows.Scan(&merchant.Merchant.MerchantId, &merchant.Merchant.Name, &merchant.Merchant.MerchantCategory, &merchant.Merchant.ImageUrl, &merchant.Merchant.Location.Lat, &merchant.Merchant.Location.Long, &merchant.Merchant.Distance, &createdAt)
		if err != nil {
			return nil, err
		}

		// get merchant items
		items, err := mr.GetMerchantItems(ctx, merchant.Merchant.MerchantId)
		if err != nil {
			return nil, err
		}

		merchant.Items = items
		merchant.Merchant.CreatedAt = timepkg.TimeToISO8601(time.Unix(createdAt, 0))
		merchants = append(merchants, merchant)
	}

	return merchants, nil
}

func (mr *merchantRepo) GetMerchantItems(ctx context.Context, merchantID string) ([]dto.ResGetItem, error) {
	q := `SELECT id, name, product_category, price, image_url, created_at FROM items WHERE merchant_id = $1`

	rows, err := mr.conn.Query(ctx, q, merchantID)
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

		result.ImageUrl = imageUrl.String
		result.CreatedAt = timepkg.TimeToISO8601(time.Unix(createdAt, 0))
		results = append(results, result)
	}

	// show results
	fmt.Printf("results: %v\n", results)

	return results, nil
}
