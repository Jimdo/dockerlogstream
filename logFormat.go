package main

func formatLogLine(m *message) (string, error) {
	jsVM.Set("hostname", hostname)
	jsVM.Set("message", *m)

	retVal, err := jsVM.Run(jsLineConverter)
	if err != nil {
		return "", err
	}

	logLine, err := retVal.ToString()
	if err != nil {
		return "", err
	}

	return logLine, nil
}
