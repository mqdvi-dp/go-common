package rpc

import (
	"fmt"

	"github.com/mqdvi-dp/go-common/env"
)

type option struct {
	tcpPort   string
	debugMode bool
}

type OptionFunc func(*option)

func getDefaultOption() option {
	return option{
		tcpPort:   fmt.Sprintf(":%d", env.GetInt("SERVICE_RPC_PORT", 6000)),
		debugMode: env.GetBool("DEBUG_MODE"),
	}
}

// SetTCPPort option func
func SetTCPPort(port int) OptionFunc {
	return func(o *option) {
		o.tcpPort = fmt.Sprintf(":%d", port)
	}
}

// SetDebugMode option func
func SetDebugMode(debugMode bool) OptionFunc {
	return func(o *option) {
		o.debugMode = debugMode
	}
}
