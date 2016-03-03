[![License: Apache v2.0](https://badge.luzifer.io/v1/badge?color=5d79b5&title=license&text=Apache+v2.0)](http://www.apache.org/licenses/LICENSE-2.0)
[![GoBuilder Download](https://badge.luzifer.io/v1/badge?color=5d79b5&title=Download&text=on GoBuilder)](https://gobuilder.me/github.com/Jimdo/dockerlogstream)
[![Build Status](https://travis-ci.org/Jimdo/dockerlogstream.svg?branch=master)](https://travis-ci.org/Jimdo/dockerlogstream)


# Jimdo / dockerlogstream

This project makes use of the fluentd log-driver implemented in Docker v1.8. Instead of fetching the logs using the `docker.sock` or having to run a ruby container with fluentd you can run this as a single binary on your system. The log format is 100% compatible to the format used by fluentd.

## Usage

### As a systemd unit

- Deploy the binary to `/usr/local/bin/dockerlogstream`
- Create following unit file

```
[Unit]
Description=DockerLogStream
After=docker.service
[Service]
Type=simple
TimeoutStartSec=0
TimeoutStopSec=0
Restart=always
RestartSec=30
SyslogIdentifier=dockerlogstream
ExecStart=/usr/local/bin/dockerlogstream --endpoint=...
```

### As a docker container

```bash
# docker run -ti -v /var/run/docker.sock:/var/run/docker.sock -p 127.0.0.1:24224:24224 Jimdo/dockerlogstream --endpoint=...
```

### Configuring your containers to log into the dockerlogstream

You need to add the option `--log-driver=fluentd` to you `docker run` command:

```bash
# docker run --log-driver=fluentd --rm -ti alpine echo "Hello World!"
```

## JavaScript line formatter

The log line formatting (and filtering) are done by simple JavaScript files. One example for a more complex solution can be found in the `lineconverter.js` file inside this repository.

```javascript
// We will get some variables set from the Go program:
// dockerlogstream    Object{}    Interaction interface for the program

// Inside dockerlogstream there are two functions and two properties available:
// Message                Object{}  Message as struct as specified in types.go
// Hostname               string    Hostname of the machine the program is running on
// SendLogLine(string)              Pass your processed log line to this function
// SkipLogLine()                    If you do filtering in here you can skip lines with this function
```

Given this you can do a very simple line formatter using this JavaScript code:

```javascript
var message = dockerlogstream.Message;
dockerlogstream.SendLogLine(
  message.Time.Format("Jan 2 15:04:05") + " " +
  dockerlogstream.Hostname + " " +
  message.Container.Names[0].substring(1) + ": " +
  message.Data);
```

For more examples see the `example` folder in this repository.
