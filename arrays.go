package go_utils


func StringInArray(value string, list []string) (bool) {
	for _,v := range list {
		if value == v {
			return true
		}
	}
	return false
}

