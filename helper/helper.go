package helper

import "log"

//FailOnError logs to stdout and os.Exit(1)
func FailOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
