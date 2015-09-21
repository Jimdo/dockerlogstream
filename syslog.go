package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/cenkalti/backoff"
)

var (
	hostname string
)

const (
	readWriteTimeout = 1 * time.Second
	tcpDialTimeout   = 5 * time.Second
)

func init() {
	hostname, _ = os.Hostname()
}

func NewSyslogAdapter(address string) (*SyslogAdapter, error) {
	return &SyslogAdapter{
		address: address,
		dialer: &net.Dialer{
			Timeout: tcpDialTimeout,
		},
	}, nil
}

type SyslogAdapter struct {
	address string
	dialer  *net.Dialer
}

func (a *SyslogAdapter) Stream(logstream chan *message) {
	backoff.Retry(func() error {
		conn, err := a.dialer.Dial("tcp", a.address)
		if err != nil {
			fmt.Printf("syslog: Unable to dial to remote address")
			return fmt.Errorf("Catch me if you can.")
		}
		defer conn.Close()

		b := new(bytes.Buffer)
		for msg := range logstream {
			b.Reset()

			logline, err := formatLogLine(msg)
			if err != nil {
				fmt.Printf("syslog: Unable to compile log line: %s\n", err)
				return fmt.Errorf("Catch me if you can.")
			}
			fmt.Fprintf(b, "%s\n", logline)

			if err := conn.SetDeadline(time.Now().Add(readWriteTimeout)); err != nil {
				fmt.Printf("syslog: Unable to set deadline: %s\n", err)
				return fmt.Errorf("Catch me if you can.")
			}

			logLine := b.Bytes()
			written, err := io.Copy(conn, b)

			if err != nil {
				if written > 0 {
					fmt.Printf("syslog: (%d/%d) %s\n", written, len(logLine), err)
				} else {
					fmt.Printf("syslog: %s\n", err)
				}
				return fmt.Errorf("syslog: %s", err)
			}
		}

		fmt.Printf("syslog: I got out of the channel watch. This should never happen.")
		return fmt.Errorf("Wat? Why am I here?")

	}, &backoff.ZeroBackOff{})
}
