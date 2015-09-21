// We will get some variables set from the Go program:
// message    Object{}    Message as struct as specified in types.go
// hostname   string      Hostname of the machine the program is running on

// Service / Component
var container_name = message.Container.Names[0].substring(1)
var parsed_container = container_name.match(/-[0-9]+-([^_]+)--([^_]+)-[a-z0-9]+$/)

if (parsed_container !== null) {
  service = parsed_container[1]
  component = parsed_container[2]
} else {
  service = hostname
  component = container_name
}

// Compile line
result = "<22> " + message.Time.Format("Jan 2 15:04:05") + " " + service + " " + component + ": " + message.Data;
