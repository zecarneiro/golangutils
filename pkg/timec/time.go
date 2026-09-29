package timec

import "time"

func Sleep(second int) {
	time.Sleep(time.Second * time.Duration(second))
}
