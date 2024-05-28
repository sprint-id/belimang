package repo

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	conn *pgxpool.Pool

	User     *userRepo
	Item     *itemRepo
	Merchant *merchantRepo
	Estimate *estimateRepo
}

func NewRepo(conn *pgxpool.Pool) *Repo {
	repo := Repo{}
	repo.conn = conn

	repo.User = newUserRepo(conn)
	repo.Item = newItemRepo(conn)
	repo.Merchant = newMerchantRepo(conn)
	repo.Estimate = newEstimateRepo(conn)

	return &repo
}
