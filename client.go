package atlas_sdk

import (
	"fmt"
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
	// Error check for required parameters
	if opts.URL == nil {
		return nil, fmt.Errorf("URL in SchemaInspectOptions must be defined")
	}

	// sets command and action
	c.command = commandSchema
	c.action = actionInspect

	// create args slice
	cmdArgs := make([]flag, 0, 3+len(opts.Schemas))

	// args.Url
	cmdArgs = append(cmdArgs, flag{flagURL, opts.URL.String()})

	// args.Schemas
	if len(opts.Schemas) > 0 {
		for _, schema := range opts.Schemas {
			cmdArgs = append(cmdArgs, flag{flagSchema, schema})
		}
	}

	// args.Exclude
	if opts.Exclude != "" {
		cmdArgs = append(cmdArgs, flag{flagExclude, opts.Exclude})
	}

	// args.Format
	if opts.Format != nil {
		cmdArgs = append(cmdArgs, flag{flagFormat, opts.Format.GoFormat()})
	}

	// set args to client
	c.args = cmdArgs

	// Exec
	return c.exec()
}

type SchemaDiffOptions struct {
	FromURLs []*url.URL // a list of URLs to the current state: can be a database URL, an HCL or SQL schema, or a migration directory.
	ToURLs   []*url.URL // a list of URLs to the desired state: can be a database URL, an HCL or SQL schema, or a migration directory.
	DevURL   *url.URL   // (optional) a URL to the Dev-Database.
	Schemas  []string   // (optional, may be supplied multiple times) - schemas to inspect within the target database.
	Exclude  []string   // (optional, may be supplied multiple times) - filter out resources matching the given glob pattern.
	Format   Format     // (optional) - Go template to use to format the output.
	//Web      bool       // !Unsupported! (-w accepted as well) - visualize the schema diff as an ERD on Atlas Cloud.
}

// SchemaDiff compares the schema within the target database using the specified options.
//
// It returns the output as a byte slice and any errors encountered.
func (c *AtlasClient) SchemaDiff(opts SchemaDiffOptions) ([]byte, error) {
	// Error check for required parameters
	if len(opts.FromURLs) == 0 || len(opts.ToURLs) == 0 {
		return nil, fmt.Errorf("FromURLs, ToURLs in SchemaDiffOptions must be defined")
	}

	// sets command and action
	c.command = commandSchema
	c.action = actionDiff

	// create args slice
	cmdArgs := make([]flag, 0, 2+len(opts.FromURLs)+len(opts.ToURLs)+len(opts.Schemas)+len(opts.Exclude))

	// args.FromURLs
	for _, fromURL := range opts.FromURLs {
		cmdArgs = append(cmdArgs, flag{flagFromURL, fromURL.String()})
	}

	// args.ToURLs
	for _, toURL := range opts.ToURLs {
		cmdArgs = append(cmdArgs, flag{flagToURL, toURL.String()})
	}

	// args.DevURL
	if opts.DevURL != nil {
		cmdArgs = append(cmdArgs, flag{flagDevURL, opts.DevURL.String()})
	}

	// args.Schemas
	if len(opts.Schemas) > 0 {
		for _, schema := range opts.Schemas {
			cmdArgs = append(cmdArgs, flag{flagSchema, schema})
		}
	}

	// args.Exclude
	if len(opts.Exclude) > 0 {
		for _, ex := range opts.Exclude {
			cmdArgs = append(cmdArgs, flag{flagExclude, ex})
		}
	}

	// args.Format
	if opts.Format != nil {
		cmdArgs = append(cmdArgs, flag{flagFormat, opts.Format.GoFormat()})
	}

	// Add command arguments to client
	c.args = cmdArgs

	// Execution of command
	return c.exec()
}

type SchemaApplyOptions struct {
	URL      *url.URL   // URL of the database to be inspected.
	ToURLs   []*url.URL // a list of URLs to the desired state: can be a database URL, an HCL or SQL schema, or a migration directory.
	DevURL   *url.URL   // (optional) a URL to the Dev-Database.
	Schemas  []string   // (optional, may be supplied multiple times) - schemas to inspect within the target database.
	Exclude  []string   // (optional, may be supplied multiple times) - filter out resources matching the given glob pattern.
	Format   Format     // (optional) - Go template to use to format the output.
	Approval bool       // (optional) Users that wish to automatically approve may run the schema apply command with the --auto-approve flag.
	DryRun   bool       // (optional) In order to skip the execution of the SQL queries against the target database, users may provide the --dry-run flag.
}

// SchemaApply applies the desired schema changes to the target database using the specified options.
//
// It returns the output as a byte slice and an error, if any.
func (c *AtlasClient) SchemaApply(opts SchemaApplyOptions) ([]byte, error) {
	// Error check for required parameters
	if opts.URL == nil || len(opts.ToURLs) == 0 {
		return nil, fmt.Errorf("URL and ToURLs in SchemaApplyOptions must be defined")
	}

	// sets command and action
	c.command = commandSchema
	c.action = actionApply

	// create args slice
	cmdArgs := make([]flag, 0, 5+len(opts.ToURLs)+len(opts.Schemas)+len(opts.Exclude))

	// args.URL
	cmdArgs = append(cmdArgs, flag{flagURL, opts.URL.String()})

	// args.ToURLs
	for _, toURL := range opts.ToURLs {
		cmdArgs = append(cmdArgs, flag{flagToURL, toURL.String()})
	}

	// args.DevURL
	if opts.DevURL != nil {
		cmdArgs = append(cmdArgs, flag{flagDevURL, opts.DevURL.String()})
	}

	// args.Schemas
	if len(opts.Schemas) > 0 {
		for _, schema := range opts.Schemas {
			cmdArgs = append(cmdArgs, flag{flagSchema, schema})
		}
	}

	// args.Exclude
	if len(opts.Exclude) > 0 {
		for _, ex := range opts.Exclude {
			cmdArgs = append(cmdArgs, flag{flagExclude, ex})
		}
	}

	// args.Format
	if opts.Format != nil {
		cmdArgs = append(cmdArgs, flag{flagFormat, opts.Format.GoFormat()})
	}

	// args.Approval
	if opts.Approval {
		cmdArgs = append(cmdArgs, flag{Flag: flagApproval})
	}

	// args.DryRun
	if opts.DryRun {
		cmdArgs = append(cmdArgs, flag{Flag: flagDryRun})
	}

	// Add command arguments to client
	c.args = cmdArgs

	// Execution of command
	return c.exec()
}

// exec executes a command using the specified entrypoint, command, and action.
// It returns the combined output of the command as a byte slice and an error if any.
func (c *AtlasClient) exec() ([]byte, error) {
	cmd := exec.Command(c.entrypoint, c.command, c.action)
	for _, arg := range c.args {
		cmd.Args = append(cmd.Args, arg.Flag)
		if arg.Value != "" {
			cmd.Args = append(cmd.Args, arg.Value)
		}
	}
	return cmd.CombinedOutput()
}
