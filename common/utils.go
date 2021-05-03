package common

import (
	"fmt"
	"time"
)

func MeasureTime(f func()) {
	startTime := time.Now()

	f()

	endTime := time.Now()
	diff := endTime.Sub(startTime)
	fmt.Println("🏁 Took ===============> ", diff.Seconds(), "seconds")
}
