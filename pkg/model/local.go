package model

import "time"

var (
	nowFunc func() time.Time = func() time.Time { return time.Now() } //nolint:deadcode,unused,varcheck
)
