package ui

type Output struct {
	stdout string
	result string
}

func NewOutput() *Output {
	return &Output{stdout: "", result: ""}
}

func (o *Output) AppendOutput(output string) {
	o.stdout = o.stdout + output
}

func (o *Output) GetStdout() string {
	return o.stdout
}
