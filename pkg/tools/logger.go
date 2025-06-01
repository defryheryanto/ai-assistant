package tools

import "log"

func (r *registry) log(logString string) {
	if r.enableLog {
		log.Println(logString)
	}
}
