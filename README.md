[![License: Apache v2.0](https://badge.luzifer.io/v1/badge?color=5d79b5&title=license&text=Apache+v2.0)](http://www.apache.org/licenses/LICENSE-2.0)
[![GoBuilder Download](https://badge.luzifer.io/v1/badge?color=5d79b5&title=Download&text=on GoBuilder)](https://gobuilder.me/github.com/Jimdo/dockerlogstream)


# Jimdo / dockerlogstream

This project is a spike to replace the flaky logspout container. The daemon can be compiled statically and then run on every ECS cluster machine to fetch logs from every running Docker container. The logs afterwards get streamed to PaperTrail using TCP connections.

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
ExecStart=/usr/local/bin/dockerlogstream --papertrail-endpoint=...
```

### As a docker container

```bash
# docker run -ti -v /var/run/docker.sock:/var/run/docker.sock Jimdo/dockerlogstream --papertrail-endpoint=...
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
