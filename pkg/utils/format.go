package utils

import "time"

const (
	DateLayout       = "2006-01-02"
	ISO8601JSTLayout = "2006-01-02T15:04:00.000+09:00"
)

func ParseDateFormat(v string) (time.Time, error) {
	return time.Parse(DateLayout, v)
}
