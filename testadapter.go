package main

import "log"

type TestAdapter struct {
}

func (a *TestAdapter) Stream(logstream chan *message) {
	for msg := range logstream {
		logline, err := formatLogLine(msg)
		if err != nil {
			log.Printf("ERROR: %s", err)
		}
		log.Printf(logline)
	}
}
