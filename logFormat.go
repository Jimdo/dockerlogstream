package main

type jsVMController struct {
	logLine        string
	logLineSkipped bool
	Message        message
	Hostname       string
}

func (j *jsVMController) SendLogLine(line string) {
	j.logLine = line
}

func (j *jsVMController) SkipLogLine() {
	j.logLineSkipped = true
}

func formatLogLine(m *message) (string, bool, error) {
	controller := &jsVMController{
		Message:  *m,
		Hostname: hostname,
	}

	jsVM.Set("dockerlogstream", controller)

	_, err := jsVM.Run(jsLineConverter)
	if err != nil {
		return "", false, err
	}

	return controller.logLine, controller.logLineSkipped, nil
}
