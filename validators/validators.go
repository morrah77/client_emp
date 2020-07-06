package validators

import "strconv"

func IsValidNumberId(id string) bool {
	_, err := strconv.Atoi(id)
	// in more general case regexp could be used - the Atoi looks just more performant for short inputs
	//res, err := regexp.MatchString("^[0-9]*$", id)
	if err != nil {
		return false
	}
	return true
}
