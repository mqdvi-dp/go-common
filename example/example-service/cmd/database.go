package cmd

import (
	"github.com/mqdvi-dp/go-common/abstract"
	"github.com/mqdvi-dp/go-common/config/database"
	"github.com/mqdvi-dp/go-common/env"
)

func sqlDatabase() (abstract.SQLDatabase, abstract.SQLDatabase, abstract.SQLDatabase) {
	master, err := database.NewSqlxConnection(database.SetSqlDSN(env.GetString("DSN_MASTER", "host=pvg-stage.cnpic6ju5omf.ap-southeast-3.rds.amazonaws.com user=pvp_klikoo_transaction password=1*>VXywC40le4j!J dbname=pvp_klikoo_transaction port=5432 sslmode=require TimeZone=utc")))
	if err != nil {
		panic(err)
	}

	slave, err := database.NewSqlxConnection(database.SetSqlDSN(env.GetString("DSN_SLAVE", "host=pvg-stage.cnpic6ju5omf.ap-southeast-3.rds.amazonaws.com user=pvp_klikoo_transaction password=1*>VXywC40le4j!J dbname=pvp_klikoo_transaction port=5432 sslmode=require TimeZone=utc")))
	if err != nil {
		panic(err)
	}

	slave2, err := database.NewSqlxConnection(database.SetSqlDSN(env.GetString("DSN_SLAVE2", "host=pvg-stage.cnpic6ju5omf.ap-southeast-3.rds.amazonaws.com user=pvp_klikoo_transaction password=1*>VXywC40le4j!J dbname=pvp_klikoo_transaction port=5432 sslmode=require TimeZone=utc")))
	if err != nil {
		panic(err)
	}

	return master, slave, slave2
}

func redisDatabase() abstract.RedisDatabase {
	return database.NewRedisConnection(database.SetRedisAddress([]string{"127.0.0.1:6379"}))
}
