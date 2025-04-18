package rest

import (
	"context"
	"fmt"
	"time"

	"github.com/hellofresh/health-go/v5"
	"github.com/mqdvi-dp/go-common/abstract"
	"github.com/mqdvi-dp/go-common/factory"
)

func healthCheck(service factory.ServiceFactory) *health.Health {
	deps := service.GetDependencies()

	h, err := health.New(
		health.WithComponent(health.Component{
			Name: service.Name(),
		}),
		health.WithChecks(
			health.Config{
				Name:      "postgresql-master",
				SkipOnErr: false,
				Timeout:   5 * time.Second,
				Check: func(ctx context.Context) error {
					if deps == nil {
						return fmt.Errorf("dependencies not setup")
					}

					m := deps.GetSQLDatabase(abstract.Master)
					if m == nil {
						return nil
					}

					return m.Database().Ping()
				},
			},
			health.Config{
				Name:      "postgresql-slave",
				SkipOnErr: false,
				Timeout:   5 * time.Second,
				Check: func(ctx context.Context) error {
					if deps == nil {
						return fmt.Errorf("dependencies not setup")
					}

					s := deps.GetSQLDatabase(abstract.Slave)
					if s == nil {
						return nil
					}
					return s.Database().Ping()
				},
			},
			health.Config{
				Name:      "postgresql-slave2",
				SkipOnErr: false,
				Timeout:   5 * time.Second,
				Check: func(ctx context.Context) error {
					if deps == nil {
						return fmt.Errorf("dependencies not setup")
					}

					s := deps.GetSQLDatabase(abstract.Slave2)
					if s == nil {
						return nil
					}
					return s.Database().Ping()
				},
			},
			health.Config{
				Name:      "redis",
				SkipOnErr: true,
				Timeout:   5 * time.Second,
				Check: func(ctx context.Context) error {
					if deps == nil {
						return fmt.Errorf("dependencies not setup")
					}

					rds := deps.GetRedisDatabase()
					if rds == nil {
						return nil
					}

					return rds.Client().Ping(ctx)
				},
			},
		),
	)

	if err != nil {
		h, _ = health.New()
		return h
	}

	return h
}
