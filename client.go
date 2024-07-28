package atlas_sdk

import (
	"bytes"
	"fmt"
	"net/url"
)

// entrypoint to Atlas cli
const (
	entrypoint = "atlas"
)

// AtlasClient is a type that represents a client for interacting with the Atlas CLI.
type AtlasClient struct {
	buf        *bytes.Buffer
	entrypoint string
	command    string
	action     string
	args       []flag
}

// NewClient creates a new instance of AtlasClient with the specified entrypoint.
// The returned AtlasClient can be used to perform various operations.
//
// Returns:
//
//	*AtlasClient: A new instance of AtlasClient with the specified entrypoint.
func NewClient(buffer *bytes.Buffer) *AtlasClient {
	return &AtlasClient{
		buf:        buffer,
		entrypoint: entrypoint,
	}
}

// SchemaInspectOptions contains the optional and required parameters needed
// to inspect the schema of the target database.
//
// URL: URL of the database to be inspected - mandatory.
//
// Schemas: (optional, may be supplied multiple times) - List of schemas to be inspected within the target database.
//
// Exclude: (optional) - Filter out resources matching the given glob pattern.
//
// Format: (optional) - Go template used to format the output.
//
// !Unsupported! Web: (optional) - visualize the schema as an ERD on Atlas Cloud.
type SchemaInspectOptions struct {
	URL     *url.URL
	Schemas []string
	Exclude string
	Format  Format
	//Web     bool
}

// SchemaInspect inspects the schema within the target database using the specified options.
//
// It writes the output to the provided buffer and returns an error, if any.
func (c *AtlasClient) SchemaInspect(opts SchemaInspectOptions) error {
	// Error check for required parameters
	if opts.URL == nil {
		return fmt.Errorf("URL in SchemaInspectOptions must be defined")
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

// SchemaDiffOptions contains the optional and required parameters needed
// for comparing the schema within the target database.
//
// CurrentURLs[FromURLs]: A list of URLs to the current state - it can be a database URL, an HCL or SQL schema, or a migration directory - mandatory.
//
// DesiredURLs[ToURLs]: A list of URLs to the desired state - it can be a database URL, an HCL or SQL schema, or a migration directory - mandatory.
//
// DevURL: A URL to the development database.
//
// Schemas: (optional, may be supplied multiple times) - List of schemas to inspect within the target database.
//
// Exclude: (optional, may be supplied multiple times) - Filter out resources matching the given glob pattern.
//
// Format: (optional) - Go template used to format the output.
//
// !Unsupported! Web: (optional) - visualize the schema diff as an ERD on Atlas Cloud.
type SchemaDiffOptions struct {
	CurrentURLs []*url.URL
	DesiredURLs []*url.URL
	DevURL      *url.URL
	Schemas     []string
	Exclude     []string
	Format      Format
	//Web      bool
}

// SchemaDiff compares the schema within the target database using the specified options.
//
// It writes the output to the provided buffer and returns any errors encountered.
func (c *AtlasClient) SchemaDiff(opts SchemaDiffOptions) error {
	// Error check for required parameters
	if len(opts.CurrentURLs) == 0 || len(opts.DesiredURLs) == 0 || opts.DevURL == nil {
		return fmt.Errorf("FromURLs, ToURLs and DevURL in SchemaDiffOptions must be defined")
	}

	// sets command and action
	c.command = commandSchema
	c.action = actionDiff

	// create args slice
	cmdArgs := make([]flag, 0, 2+len(opts.CurrentURLs)+len(opts.DesiredURLs)+len(opts.Schemas)+len(opts.Exclude))

	// args.FromURLs
	for _, fromURL := range opts.CurrentURLs {
		cmdArgs = append(cmdArgs, flag{flagFromURL, fromURL.String()})
	}

	// args.ToURLs
	for _, toURL := range opts.DesiredURLs {
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

// SchemaApplyOptions contains the optional and required parameters needed
// to apply the desired schema changes to the target database.
//
// - URL: URL of the database to be inspected - mandatory.
//
// - ToURLs: A list of URLs to the desired state. It can be a database URL, an HCL or SQL schema,
// or a migration directory - mandatory.
//
// DevURL: A URL to the Dev-Database.
//
// Schemas: (optional, may be supplied multiple times) - Schemas to inspect within the target database.
//
// Exclude: (optional, may be supplied multiple times) - This helps to filter out resources matching the given glob pattern.
//
// Format: (optional) - This is a Go template used to format the output.
//
// Approval: (optional) - If users wish to automatically approve, they may run the schema apply command with the --auto-approve flag.
//
// DryRun: (optional) - To skip the execution of the SQL queries against the target database, users may provide the --dry-run flag.
type SchemaApplyOptions struct {
	CurrentURL  *url.URL
	DesiredURLs []*url.URL
	DevURL      *url.URL
	Schemas     []string
	Exclude     []string
	Format      Format
	Approval    bool
	DryRun      bool
}

// SchemaApply applies the desired schema changes to the target database using the specified options.
//
// It writes the output to the provided buffer and returns an error, if any.
func (c *AtlasClient) SchemaApply(opts SchemaApplyOptions) error {
	// Error check for required parameters
	if opts.CurrentURL == nil || len(opts.DesiredURLs) == 0 || opts.DevURL == nil {
		return fmt.Errorf("URL, ToURLs and DevURL in SchemaApplyOptions must be defined")
	}

	// sets command and action
	c.command = commandSchema
	c.action = actionApply

	// create args slice
	cmdArgs := make([]flag, 0, 5+len(opts.DesiredURLs)+len(opts.Schemas)+len(opts.Exclude))

	// args.URL
	cmdArgs = append(cmdArgs, flag{flagURL, opts.CurrentURL.String()})

	// args.ToURLs
	for _, toURL := range opts.DesiredURLs {
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
