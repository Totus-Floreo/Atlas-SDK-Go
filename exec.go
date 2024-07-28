//go:build !test

package atlas_sdk

import (
	"fmt"
	"os/exec"
)

// exec executes a command using the specified entrypoint, command, and action.
// It writes the combined output of the command to the provided buffer and returns an error if any.
func (c *AtlasClient) exec() error {
	cmd := exec.Command(c.entrypoint, c.command, c.action)
	for _, arg := range c.args {
		cmd.Args = append(cmd.Args, arg.Flag)
		if arg.Value != "" {
			cmd.Args = append(cmd.Args, arg.Value)
		}
	}
	//fmt.Println(cmd.String())
	output, err := cmd.CombinedOutput()
	c.buf.Write(output)
	if err != nil {
		return fmt.Errorf("atlas error %w: %s", err, c.buf.String())
	}
	return nil
}
