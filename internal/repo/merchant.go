package repo

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/entity"
	"github.com/sprint-id/eniqilo-server/internal/ierr"
	timepkg "github.com/sprint-id/eniqilo-server/pkg/time"
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
