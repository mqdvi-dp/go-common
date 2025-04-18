package cronexpr

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mqdvi-dp/go-common/zone"
)

const (
	// oneMinute const
	oneMinute = 1 * time.Minute
	// oneDay const
	oneDay = 24 * time.Hour
	// oneWeek const
	oneWeek = 7 * oneDay
	// oneMonth const
	oneMonth = 30 * oneDay
	// oneYear const
	oneYear = 12 * oneMonth

	// daily const
	daily = "daily"
	// weekly const
	weekly = "weekly"
	// monthly const
	monthly = "monthly"
	// yearly const
	yearly = "yearly"
)

const (
	layoutTime = "15:04"
)

func ParseDuration(t string) (duration, nextDuration time.Duration, err error) {
	interval, err := time.ParseDuration(t)
	if err == nil {
		return interval, 0, nil
	}

	delimiter := strings.Split(t, "@")

	ts, err := time.Parse(layoutTime, delimiter[0])
	if err != nil {
		return 0, 0, errors.New("time format error. must be HH:mm")
	}

	repeat := oneMinute
	if len(delimiter) > 1 {
		switch strings.ToLower(delimiter[1]) {
		case daily:
			repeat = oneDay
		case weekly:
			repeat = oneWeek
		case monthly:
			repeat = oneMonth
		case yearly:
			repeat = oneYear
		default:
			repeat, err = time.ParseDuration(delimiter[1])
			if err != nil {
				return 0, 0, fmt.Errorf(
					`invalid descriptor "%s". Must One of ("daily", "weekly", "monthly", "yearly") or duration string`,
					delimiter[1],
				)
			}
		}
	}

	now := time.Now().In(zone.TzJakarta())
	atTime := time.Date(now.Year(), now.Month(), now.Day(), ts.Hour(), ts.Minute(), 0, 0, now.Location())
	if now.Before(atTime) {
		duration = atTime.Sub(now)
	} else {
		duration = oneDay - now.Sub(atTime)
	}

	if duration < 0 {
		duration *= -1
	}

	nextDuration = repeat

	return
}
