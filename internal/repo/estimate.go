package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sprint-id/eniqilo-server/internal/dto"
	"github.com/sprint-id/eniqilo-server/internal/entity"
)

type estimateRepo struct {
	conn *pgxpool.Pool
}

func newEstimateRepo(conn *pgxpool.Pool) *estimateRepo {
	return &estimateRepo{conn}
}

func (cr *estimateRepo) CreateEstimate(ctx context.Context, sub string, estimate entity.Estimate) (dto.ResCreateEstimate, error) {
	// add estimate
	q := `INSERT INTO estimates (id, user_id, total_price, delivery_time, created_at)
	VALUES ( gen_random_uuid(), $1, $2, $3, EXTRACT(EPOCH FROM now())::bigint) RETURNING id`

	var id string
	err := cr.conn.QueryRow(ctx, q, sub, estimate.TotalPrice, estimate.DeliveryTime).Scan(&id)
	if err != nil {
		fmt.Printf("error query: %v\n", err)
		return dto.ResCreateEstimate{}, err
	}

	return dto.ResCreateEstimate{
		CalculatedEstimateID:           id,
		TotalPrice:                     estimate.TotalPrice,
		EstimatedDeliveryTimeInMinutes: estimate.DeliveryTime}, nil
}
