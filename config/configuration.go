package config

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/mqdvi-dp/go-common/abstract"
	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/logger"
	"github.com/mqdvi-dp/go-common/monitoring"
)

var (
	once sync.Once
)

type config struct {
	serviceName string
	// close all applications servers
	closers []abstract.Closer
}

// Config represents for implementing configuration applications
type Config interface {
	// GetServiceName returns the service name
	GetServiceName() string

	// Exit will close the application server
	Exit()

	// Injections will inject selected dependencies
	Injections(depsFunc func(context.Context) []abstract.Closer)
}

// New returns a instance of Config
func New(serviceName string) Config {
	envDriver := env.GetString("ENV_DRIVER")
	logger.Log.Printf(context.Background(), "env driver: %s", envDriver)
	// load configurations based on env
	// possible values for now is
	// - vault_file
	// - etcd
	//
	// if there's not yet set, we will
	// return panic
	switch strings.ToLower(envDriver) {
	case "vault_file":
		loadVaultFile(env.GetString("CONFIG_FILE"))
	case "etcd":
		loadEtcd(serviceName)
	default:
		envVal := env.GetString("ENV")
		// if env is not set or is local or dev, we will not set the env driver
		if envVal == "" || (strings.EqualFold(envVal, "local") || strings.EqualFold(envVal, "dev")) {
			break
		}
		logger.Log.Fatalf("env driver not yet set")
	}

	// start the db migrations if exists
	logger.BlueItalic("========= check db migrations =========")
	DbAutoMigrations()

	if !strings.EqualFold(env.GetString("OPENSEARCH_FLOW_SEND_DATA"), "queue") {
		// logger.OpenSearch()
		logger.Elasticsearch()
	} else {
		logger.Elasticsearch()

	}

	// start the prometheus server
	monitoring.New(serviceName)

	return &config{serviceName: serviceName}
}

func (c *config) GetServiceName() string {
	return c.serviceName
}

func (c *config) Injections(depsFunc func(context.Context) []abstract.Closer) {
	// set timeout for init configuration
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	result := make(chan error)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				result <- fmt.Errorf("failed init configuration :=> %v", r)
			}
			close(result)
		}()

		c.closers = depsFunc(ctx)
	}()

	// with timeout to init configuration
	select {
	case <-ctx.Done():
		logger.Log.Fatal(fmt.Errorf("timeout to load selected dependencies: %v", ctx.Err()))
	case err := <-result:
		if err != nil {
			logger.Log.Fatal(err)
		}
		return
	}
}

func (c *config) Exit() {
	// set timeout for close all configuration
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
	defer cancel()

	errCloseChan := make(chan error)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				errCloseChan <- fmt.Errorf("failed close connection :=> %v", r)
			}
			close(errCloseChan)
		}()

		for _, cl := range c.closers {
			_ = cl.Disconnect(ctx)
		}
	}()

	// for force exit
	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, os.Interrupt, syscall.SIGTERM)

	// with timeout to close all configuration
	select {
	case <-quitSignal:
		fmt.Println("\x1b[31;1mForce exit\x1b[0m")
	case <-ctx.Done():
		logger.Log.Fatal(fmt.Errorf("timeout to close all selected dependencies connection: %v", ctx.Err()))
	case err := <-errCloseChan:
		if err != nil {
			logger.Log.Fatal(err)
		}
		fmt.Printf("\x1b[32;1mSuccess close all config dependency at %s\x1b[0m\n", c.GetServiceName())
	}
}
