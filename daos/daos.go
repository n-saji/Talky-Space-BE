package daos

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxDao struct {
	pool *pgxpool.Pool
}

func NewPgxDao(pool *pgxpool.Pool) *PgxDao {
	return &PgxDao{pool: pool}
}
