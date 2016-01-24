var message = dockerlogstream.Message;
dockerlogstream.SendLogLine(
  message.Time.Format("Jan 2 15:04:05") + " " +
  dockerlogstream.Hostname + " " +
  message.Container.Names[0].substring(1) + ": " +
  message.Data);
