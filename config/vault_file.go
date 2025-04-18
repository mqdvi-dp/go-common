package config

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/mqdvi-dp/go-common/env"
	"github.com/mqdvi-dp/go-common/logger"
)

// loadVaultFile is a function that loads a vault file and sets the environment variables
func loadVaultFile(pathFile string) {
	once.Do(func() {
		logger.Log.Printf(context.Background(), "pathFile: %s", pathFile)
		d, err := os.ReadFile(pathFile)
		if err != nil {
			log.Fatalf("Error reading file: %s", err)
		}

		var kv = make(map[string]interface{})
		err = json.Unmarshal(d, &kv)
		if err != nil {
			log.Fatalf("Error unmarshalling file: %s", err)
		}

		if strings.EqualFold(env.GetString("ENV"), "staging") {
			logger.Log.Printf(context.Background(), "kv: %v", kv)
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
	})
}
