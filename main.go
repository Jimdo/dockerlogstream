package main

import (
	"log"
	"sync"
	"time"

	"github.com/Luzifer/rconfig"
	"github.com/fsouza/go-dockerclient"
	"github.com/robertkrimen/otto"
)

var (
	version = "dev"

	logstream chan *message
	cfg       config
	client    *docker.Client
	err       error

	receiveChannels     = map[string]*bufferedLogMessageWriter{}
	receiveChannelsLock = &sync.RWMutex{}

	jsVM            = otto.New()
	jsLineConverter *otto.Script
)

type config struct {
	PapertrailEndpoint string `flag:"papertrail-endpoint" description:"Logging target in PaperTrail (TCP, Plain)"`
	DockerAPI          string `flag:"docker-endpoint" default:"/var/run/docker.sock" description:"Address of docker endpoint to use"`
	Testing            bool   `flag:"testing" default:"false" description:"Do not stream but write to STDOUT"`
	LineConverter      string `flag:"line-converter" default:"lineconverter.js" description:"Sets the JavaScript to compile the log line to be sent"`
}

func init() {
	rconfig.Parse(&cfg)
	logstream = make(chan *message, 5000)

	jsLineConverter, err = jsVM.Compile(cfg.LineConverter, nil)
	if err != nil {
		log.Fatalf("Unable to parse line converter script: %s", err)
	}
}

func main() {
	// Connect to Docker socket
	client, err = docker.NewClient("unix://" + cfg.DockerAPI)
	if err != nil {
		log.Fatalf("Unable to connect to Docker daemon: %s", err)
	}

	if cfg.Testing {
		// Log to STDOUT instead of streaming
		ta := TestAdapter{}
		go ta.Stream(logstream)
	} else {
		// Create sending part of the logger
		sl, err := NewSyslogAdapter(cfg.PapertrailEndpoint)
		if err != nil {
			log.Fatalf("Unable to create logging adapter: %s", err)
		}
		go sl.Stream(logstream)
	}

	for {
		containers, err := listContainers()
		if err != nil {
			log.Fatalf("Unable to list containers: %s", err)
		}

		touched := []string{}
		for _, container := range containers {
			touched = append(touched, container.ID)

			// If we don't have a channel for this container create it
			if _, ok := receiveChannels[container.ID]; !ok {
				log.Printf("INFO: Found container %s, attaching...", container.ID)
				receiveChannelsLock.Lock()
				receiveChannels[container.ID] = &bufferedLogMessageWriter{
					channel:   logstream,
					container: container,
					buffer:    []byte{},
				}
				receiveChannelsLock.Unlock()

				go func(containerID string) {
					receiveChannelsLock.RLock()
					stream := receiveChannels[containerID]
					receiveChannelsLock.RUnlock()

					err := client.AttachToContainer(docker.AttachToContainerOptions{
						Container:    containerID,
						OutputStream: stream,
						ErrorStream:  stream,
						Logs:         false,
						Stream:       true,
						Stdout:       true,
						Stderr:       true,
					})

					log.Printf("ERROR: Unable to attach or attach ended for container %s: %s", containerID, err)
					receiveChannelsLock.Lock()
					delete(receiveChannels, containerID)
					receiveChannelsLock.Unlock()
				}(container.ID)
			}

		}

		for k := range receiveChannels {
			if !inSlice(touched, k) {
				log.Printf("INFO: Missing container %s, removing logger...", k)
				receiveChannelsLock.Lock()
				delete(receiveChannels, k)
				receiveChannelsLock.Unlock()
			}
		}
		time.Sleep(time.Second)
	}
}

func listContainers() ([]docker.APIContainers, error) {
	containers, err := client.ListContainers(docker.ListContainersOptions{
		All: false, // Do not list dead containers
	})

	if err != nil {
		return []docker.APIContainers{}, err
	}

	return containers, nil
}

func inSlice(s []string, k string) bool {
	for _, v := range s {
		if v == k {
			return true
		}
	}

	return false
}
