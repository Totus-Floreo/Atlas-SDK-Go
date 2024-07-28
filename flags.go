package atlas_sdk

import "fmt"

// flags consts
const (
	flagURL      = "--url"
	flagSchema   = "--schema"
	flagExclude  = "--exclude"
	flagFormat   = "--format"
	flagFromURL  = "--from"
	flagToURL    = "--to"
	flagDevURL   = "--dev-url"
	flagApproval = "--auto-approve"
	flagDryRun   = "--dry-run"
)

type flag struct {
	Flag  string
	Value string
}

func (f *flag) String() string {
	if f.Value == "" {
		return f.Flag
	}
	return fmt.Sprintf(`%s "%s"`, f.Flag, f.Value)
}
