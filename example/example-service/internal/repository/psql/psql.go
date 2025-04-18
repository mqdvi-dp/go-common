package psql

import "github.com/mqdvi-dp/go-common/config/database/dbc"

type psqlRepository struct {
	master, slave, slave2 dbc.SqlDbc
}

func New(m, s, s2 dbc.SqlDbc) *psqlRepository {
	return &psqlRepository{master: m, slave: s, slave2: s2}
}
