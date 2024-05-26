package logger

import "log"


func LogAndExit(e error) {
	if e!= nil {
		log.Println(e)
		panic(e)
	}
}