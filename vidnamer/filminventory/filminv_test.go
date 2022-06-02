package filminventory

import (
	"testing"
)

func TestFilm_DesiredFilename(t *testing.T) {
	type fields struct {
		Signature  string
		Filename   string
		Title      string
		Author     string
		SourceSite string
		Keywords   []string
		Hh         int
		Room       string
		Test       string
		Duration   int
		FileExt    string
		Tags       map[string]string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "regular",
			fields: fields{
				Title:      "mytitle",
				SourceSite: "amz",
				Keywords:   []string{"key1", "key2"},
				Hh:         1,
				Room:       "main",
				Test:       "test02",
				Duration:   4,
			},
			want: "mytitle__amz__key1-key2__hh1-main-test02-d04.mp4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Film{
				Signature:  tt.fields.Signature,
				Filename:   tt.fields.Filename,
				Title:      tt.fields.Title,
				Author:     tt.fields.Author,
				SourceSite: tt.fields.SourceSite,
				Keywords:   tt.fields.Keywords,
				Hh:         tt.fields.Hh,
				Room:       tt.fields.Room,
				Test:       tt.fields.Test,
				Duration:   tt.fields.Duration,
				FileExt:    tt.fields.FileExt,
				Tags:       tt.fields.Tags,
			}
			if got := f.DesiredFilename(); got != tt.want {
				t.Errorf("Film.DesiredFilename() =\ngot  %q\nwant %q\n", got, tt.want)
			}
		})
	}
}
