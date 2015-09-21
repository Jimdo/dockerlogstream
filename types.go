package main

import (
	"strings"
	"time"

	"github.com/fsouza/go-dockerclient"
)

type message struct {
	Container docker.APIContainers
	Data      string
	Time      time.Time
}

type bufferedLogMessageWriter struct {
	buffer    []byte
	channel   chan *message
	container docker.APIContainers
	client    *docker.Client
}

func (b *bufferedLogMessageWriter) Write(p []byte) (int, error) {
	tmp := string(append(b.buffer, p...))

	linesEmitted := 0

	for strings.Contains(tmp, "\n") {
		pos := strings.Index(tmp, "\n")

		b.channel <- &message{
			Container: b.container,
			Data:      strings.TrimSpace(tmp[:pos]),
			Time:      time.Now(),
		}

		linesEmitted++

		if pos+1 >= len(tmp) {
			tmp = ""
		} else {
			tmp = tmp[pos+1:]
		}
	}

	b.buffer = []byte(tmp)

	return len(p), nil
}
