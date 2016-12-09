package exe

type ExitError struct {
	error
	ExitCode int
}
