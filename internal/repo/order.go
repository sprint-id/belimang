package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sprint-id/belimang/internal/dto"
	"github.com/sprint-id/belimang/internal/entity"
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
