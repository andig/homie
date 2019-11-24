package paho

import "log"

// Log implements paho.ErrorHandler to log mqtt errors
func Log(err error) {
	log.Printf("%v", err)
}
