package zone

import "time"

func WIB(t time.Time, locs ...*time.Location) Zone {
	loc := TzJakarta() // default location
	if len(locs) > 0 {
		loc = locs[0]
	}

	t = t.In(loc)
	return &zone{value: t, tz: t.Location()}
}

func FromString(val string, layout string) Zone {
	t, err := time.Parse(layout, val)
	if err != nil {
		return nil
	}

	return &zone{value: t, tz: t.Location()}
}
