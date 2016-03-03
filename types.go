package main

import (
	"time"

	"github.com/fsouza/go-dockerclient"
)

type message struct {
	Container *docker.Container
	Data      string
	Time      time.Time
}
