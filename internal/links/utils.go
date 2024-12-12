package links

import (
	"fmt"
	"regexp"
)

var utmValuePattern = regexp.MustCompile(`^[a-zA-Z0-9\-_\.]+$`)

func validateUTMValue(val string) error {
	if !utmValuePattern.MatchString(val) {
		return fmt.Errorf("invalid utm value: %s", val)
	}

	return nil
}
