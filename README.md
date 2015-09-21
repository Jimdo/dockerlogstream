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
