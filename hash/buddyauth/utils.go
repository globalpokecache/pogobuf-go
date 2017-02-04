package buddyauth

import (
	"fmt"
)

func debug(format string, a ...interface{}) {
	if Debug {
		fmt.Printf(fmt.Sprintf("(BuddyAuth) %s\n", format), a...)
	}
}
