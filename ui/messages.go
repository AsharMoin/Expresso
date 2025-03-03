package ui

// Response is the message returned when a command is generated
type Response struct {
	command     string
	description string
}

// Exiting is the message sent when the application is quitting
type Exiting struct {
	success string
	output  string
}
