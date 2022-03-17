package battery

import (
	"fmt"
	"os/exec"
)

type DefaultDeviceNotFoundError struct{}

func (*DefaultDeviceNotFoundError) Error() string {
	return "BAT#"
}

type InvalidDeviceError struct {
	Name string
}

func (e *InvalidDeviceError) Error() string {
	return fmt.Sprintf("invalid device: %q", e.Name)
}

type CommandError struct {
	Cmd *exec.Cmd
	Err error
}

func (e *CommandError) Error() string {
	return fmt.Sprintf(
		"there was an error running the command %#q: %s",
		e.Cmd.String(),
		e.Err)
}

func (e *CommandError) Unwrap() error {
	return e.Err
}
