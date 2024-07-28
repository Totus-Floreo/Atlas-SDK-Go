package atlas_sdk

import "fmt"

// flags consts
const (
	flagURL     = "--url"
	flagSchema  = "--schema"
	flagExclude = "--exclude"
	flagFormat  = "--format"
)

type flag struct {
	Flag  string
	Value string
}

func (f *flag) String() string {
	return fmt.Sprintf(`%s "%s"`, f.Flag, f.Value)
}
