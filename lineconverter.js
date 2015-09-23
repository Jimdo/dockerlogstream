// Service / Component
var message = dockerlogstream.Message;
var container_name = message.Container.Names[0].substring(1);
var parsed_container = container_name.match(/-[0-9]+-([^_]+)--([^_]+)-[a-z0-9]+$/);

if (parsed_container !== null) {
  service = parsed_container[1];
  component = parsed_container[2];
} else {
  service = dockerlogstream.Hostname;
  component = container_name;
}

// Compile line
dockerlogstream.SendLogLine("<22> " +
  message.Time.Format("Jan 2 15:04:05") + " " +
  service + " " + component + ": " +
  message.Data);
