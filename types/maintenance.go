package types

import "time"

type Maintenance struct {
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	WhitelistedUsers []string  `json:"whitelisted_users"`
	Channels         []string  `json:"channels"`
}
