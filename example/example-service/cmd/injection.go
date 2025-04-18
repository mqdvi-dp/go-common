package cmd

import (
	"context"

	"github.com/mqdvi-dp/go-common/abstract"
	"github.com/mqdvi-dp/go-common/config"
)

func dependencies(cfg config.Config) (deps abstract.Dependency) {
	cfg.Injections(func(ctx context.Context) []abstract.Closer {
		// master, slave := sqlDatabase()
		// redis := redisDatabase()

		// put here if we've middleware

		deps = abstract.New(
		// abstract.SetSQLDatabase(abstract.Master, master),
		// abstract.SetSQLDatabase(abstract.Slave, slave),
		// abstract.SetRedisDatabase(redis),
		)

		return []abstract.Closer{
			// master, slave,
		}
	})

	return
}
