package time

import (
	"time"
)

func GetDayStr() string {
	return time.Now().Format("20060102")
}
