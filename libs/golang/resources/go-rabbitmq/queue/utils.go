package queue

import "log"

func failOnError(err error, msg string) bool {
	if err != nil {
		log.Printf("%s: %s", msg, err)
		return true
	}
	return false
}
