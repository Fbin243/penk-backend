package time

import (
	"fmt"
	"time"
)

func GetTheCurrentTimeInFormat(format string) {
	fmt.Println("Current time:", time.Now().Format(format))
}
