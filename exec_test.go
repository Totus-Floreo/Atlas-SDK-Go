//go:build test

package atlas_sdk

import (
	"os/exec"
)

func (c *AtlasClient) exec() error {
	cmd := exec.Command(c.entrypoint, c.command, c.action)
	for _, arg := range c.args {
		cmd.Args = append(cmd.Args, arg.Flag)
		if arg.Value != "" {
			cmd.Args = append(cmd.Args, arg.Value)
		}
	}
	c.buf.WriteString(cmd.String())
	return nil
}
