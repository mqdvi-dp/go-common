package featureflag

import "github.com/mqdvi-dp/go-common/env"

func IsUsePaymentSchema(username string) bool {
	return implementWhitelistedUser(username, env.GetListString("FF_WHITELISTED_PAYMENT_SCHEMA_USERS"), env.GetBool("IS_USE_PAYMENT_SCHEMA"))
}

func IsUseOrderTable(username string) bool {
	return implementWhitelistedUser(username, env.GetListString("FF_WHITELISTED_ORDER_TABLE_USERS"), env.GetBool("FF_IS_USE_ORDER_TABLE"))
}

func implementWhitelistedUser(username string, envWhitelisted []string, envIgnoreWhitelisted bool) bool {
	var (
		isWhitelisted = false
		whitelists    = envWhitelisted
	)
	for _, whitelist := range whitelists {
		if whitelist == username {
			isWhitelisted = true
			break
		}
	}

	return envIgnoreWhitelisted || isWhitelisted
}
