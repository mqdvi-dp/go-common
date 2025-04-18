package abstract

import (
	"github.com/mqdvi-dp/go-common/config/database/dbc"
	"github.com/mqdvi-dp/go-common/config/database/rdc"
)

type Instance string

const (
	Master Instance = "master"
	Slave  Instance = "slave"
	Slave2 Instance = "slave2"
)

type SQLDatabase interface {
	Database() dbc.SqlDbc

	Closer
}

type RedisDatabase interface {
	Client() rdc.Rdc

	Closer
}
