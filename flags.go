package atlas_sdk

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
