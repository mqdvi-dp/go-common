package helpers

import "fmt"

// PhoneNumber formatted phoneNumber with Indonesian Standard Phone Number
func PhoneNumber(val string) string {
	if val[:2] == "08" || val[:1] == "0" {
		return fmt.Sprintf("62%s", val[1:])
	}

	if val[:3] == "628" {
		return val
	}

	// send the same with request
	return val
}
