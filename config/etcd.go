package config

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/kvtools/etcdv3"
	"github.com/mqdvi-dp/go-common/env"
)

func loadEtcd(key string) {
	once.Do(
		func() {
			var kv = make(map[string]interface{})
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
			defer cancel()

			// get host from env variable
			host := env.GetString("ETCD_ADDR")
			if host == "" {
				return
			}
			// make it array
			hosts := strings.Split(strings.TrimSpace(host), ";")
			// init etcd connection
			store, err := etcdv3.New(ctx, hosts, nil)
			if err != nil {
				panic(err)
			}
			// get key-values from etcd based on keys
			keyValues, err := store.Get(ctx, key, nil)
			if err != nil {
				panic(err)
			}
			// unmarshal the key values data
			err = json.Unmarshal(keyValues.Value, &kv)
			if err != nil {
				panic(err)
			}

			// set to environment
			for key, value := range kv {
				val, err := reflection(value)
				if err != nil {
					panic(err)
				}
				// logger.Log.Debugf(context.Background(), "key: %s || value: %s", key, val)
				// set to environment
				// set key with uppercase
				err = os.Setenv(strings.ToUpper(key), val)
				if err != nil {
					panic(err)
				}
			}
		},
	)
}
