package service

import (
	"talky-space-be/daos"
)

type Service struct {
	daos *daos.PgxDao
}

func New(Dao *daos.PgxDao) *Service {
	return &Service{
		daos: Dao,
	}
}
