package atlas_sdk

import (
	"bytes"
	"net/url"
	"testing"
)

func TestAtlasClient_SchemaApply(t *testing.T) {
	tests := []struct {
		name    string
		opts    SchemaApplyOptions
		want    string
		wantErr bool
	}{
		{
			name: "NilCurrentURL",
			opts: SchemaApplyOptions{
				CurrentURL: nil,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "EmptyDesiredURLs",
			opts: SchemaApplyOptions{
				CurrentURL:  &url.URL{},
				DesiredURLs: []*url.URL{},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "NilDevURL",
			opts: SchemaApplyOptions{
				CurrentURL:  &url.URL{},
				DesiredURLs: []*url.URL{{}},
				DevURL:      nil,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "For checking command string",
			opts: SchemaApplyOptions{
				CurrentURL:  &url.URL{Path: "current"},
				DesiredURLs: []*url.URL{{Path: "desired"}},
				DevURL:      &url.URL{Path: "dev"},
				Schemas:     []string{"Schema1", "Schema2"},
				Exclude:     []string{"Exclude1", "Exclude2"},
				Format:      nil,
				Approval:    false,
				DryRun:      false,
			},
			want:    "atlas schema apply --url \"current\" --to \"desired\" --dev-url \"dev\" --schema \"Schema1\" --schema \"Schema2\" --exclude \"Exclude1\" --exclude \"Exclude2\"",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			c := NewClient(buf)
			if err := c.SchemaApply(tt.opts); (err != nil) != tt.wantErr {
				t.Errorf("SchemaApply() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.want != buf.String() {
				t.Errorf("SchemaApply() got = %v, want %v", buf.String(), tt.want)
			}
		})
	}
}
