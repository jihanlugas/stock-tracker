package utils

import (
	"stock-tracker/constant"
	"time"
)

func ParseTime(t string) (time.Time, error) {
	return time.Parse(constant.FormatTimeLayout, t)
}
