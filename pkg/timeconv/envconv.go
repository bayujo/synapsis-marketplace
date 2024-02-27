package envconv

import(
	"time"
)

func ParseDuration(timestr string) time.Duration {
	duration, _ := time.ParseDuration(timestr)
	
	return duration
}