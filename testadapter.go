package main

import "log"

type TestAdapter struct {
}

func (a *TestAdapter) Stream(logstream chan *message) {
	for msg := range logstream {
		logline, skipLogLine, err := formatLogLine(msg)
		if err != nil {
			log.Printf("ERROR: %s", err)
		}

		if skipLogLine {
			continue
		}

		log.Printf(logline)
	}
}
