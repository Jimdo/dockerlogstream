var message = dockerlogstream.Message;
var container_name = message.Container.Names[0].substring(1);
var loggly_token = "YOUR_TOKEN";

dockerlogstream.SendLogLine(
  "<22>1 " +
  message.Time.UTC().Format("2006-01-02T15:04:05Z") + " " +
  dockerlogstream.Hostname + " " +
  container_name + " - - " +
  "[" + loggly_token + "@41058 tag=\"" + container_name + "\"] " +
  message.Data);
