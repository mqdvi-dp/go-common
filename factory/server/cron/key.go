package cron

import "encoding/json"

type schedulerKey struct {
	JobName  string `json:"jobName"`
	Interval string `json:"interval"`
}

// String implement stringer
func (sk schedulerKey) String() string {
	s, _ := json.Marshal(sk)

	return string(s)
}

// CreateSchedulerKey helper
//
// interval value allowed:
//
// * Cron Expression: e.g.: 1/* * * * *
//
// * Duration: e.g.: 2s, 10m, 1h
//
// * Custom start time and repeat duration, e.g:
//   - 07:00@daily, will start at 07:00 UTC+7 and repeat every day
//   - 07:00@weekly, will start at 07:00 UTC+7 and repeat every week
//   - 07:00@monthly, will start at 07:00 UTC+7 and repeat every month
//   - 07:00@10s, will start at 07:00 UTC+7 and next repeat every 10 seconds
//   - 07:00@1m, will start at 07:00 UTC+7 and next repeat every 1 minute
func CreateSchedulerKey(jobName, interval string) string {
	return schedulerKey{JobName: jobName, Interval: interval}.String()
}

// ParseSchedulerKey helpers
func ParseSchedulerKey(val string) (string, string) {
	var sk schedulerKey
	err := json.Unmarshal([]byte(val), &sk)
	if err != nil {
		return "", ""
	}

	return sk.JobName, sk.Interval
}
