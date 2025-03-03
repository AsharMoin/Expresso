package ui

type Output struct {
	stdout string
	result string
}

func NewOutput() *Output {
	return &Output{stdout: "", result: ""}
}

// SetStdout sets the stdout content
func (o *Output) SetStdout(content string) {
	o.stdout = content
}

// AppendOutput appends content to stdout
func (o *Output) AppendOutput(content string) {
	if o.stdout == "" {
		o.stdout = content
	} else {
		o.stdout += content
	}
}

// GetStdout returns the current stdout content
func (o *Output) GetStdout() string {
	return o.stdout
}
