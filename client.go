package atlas_sdk

import (
	"net/url"
	"os/exec"
)

// entrypoint to Atlas cli
const (
	entrypoint = "atlas"
)

// AtlasClient is a type that represents a client for interacting with the Atlas CLI.
type AtlasClient struct {
	entrypoint string
	command    string
	action     string
	args       []flag
}

// NewClient creates a new instance of AtlasClient with the specified entrypoint.
// The returned AtlasClient can be used to perform various operations.
// Returns:
//
//	*AtlasClient: A new instance of AtlasClient with the specified entrypoint.
func NewClient() *AtlasClient {
	return &AtlasClient{
		entrypoint: entrypoint,
	}
}

type SchemaInspectOptions struct {
	URL     *url.URL
	Schemas []string // (optional, may be supplied multiple times) - schemas to inspect within the target database.
	Exclude string   // (optional) - filter out resources matching the given glob pattern.
	Format  Format   // (optional) - Go template to use to format the output.
	//Web     bool     // !Unsupported! (optional) - visualize the schema as an ERD on Atlas Cloud.
}

// SchemaInspect inspects the schema within the target database using the specified options.
//
// It returns the output as a byte slice and an error, if any.
func (c *AtlasClient) SchemaInspect(opts SchemaInspectOptions) ([]byte, error) {
	// command
	c.command = commandSchema

	// action
	c.action = actionInspect

	// create args slice
	cmdArgs := make([]flag, 3+len(opts.Schemas))

	// args.Url
	cmdArgs = append(cmdArgs, flag{flagURL, opts.URL.String()})

	// args.Schemas
	for _, schema := range opts.Schemas {
		cmdArgs = append(cmdArgs, flag{flagSchema, schema})
	}

	// args.Exclude
	cmdArgs = append(cmdArgs, flag{flagExclude, opts.Exclude})

	// args.Format
	cmdArgs = append(cmdArgs, flag{flagFormat, opts.Format.GoFormat()})

	// set args to client
	c.args = cmdArgs

	// Exec
	return c.exec()
}

// exec executes a command using the specified entrypoint, command, and action.
// It returns the combined output of the command as a byte slice and an error if any.
func (c *AtlasClient) exec() ([]byte, error) {
	cmd := exec.Command(c.entrypoint, c.command, c.action)
	return cmd.CombinedOutput()
}
