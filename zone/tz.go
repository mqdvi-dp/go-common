package zone

import "time"

const (
	tzJakarta = "Asia/Jakarta"
)

func TzJakarta() *time.Location {
	tz, _ := time.LoadLocation(tzJakarta)
	return tz
}
