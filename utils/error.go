package utils

import "github.com/go-logr/logr"

// CheckErr checks error with a message
func CheckErr(err error, msg string, log logr.Logger) {
	if err != nil {
		log.Error(err, msg)
	}
}
