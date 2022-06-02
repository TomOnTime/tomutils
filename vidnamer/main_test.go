package main

import (
	"testing"

	"github.com/TomOnTime/tomutils/vidnamer/filminventory"
)

func cleanFilename(old string) string {
	f := filminventory.ParseFilename(old)
	f.Signature = "d41d8cd98f00b204e9800998ecf8427e"

	return f.DesiredFilename()
}

func Test_cleanFilename(t *testing.T) {
	type args struct {
		old string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "00",
			args: args{`P G PMV - BIAS____fvoiceover-foo__hh1_side_d04.mp4`},
			want: `P G PMV - BIAS____fvoiceover-foo__hh1-side-d04.mp4`,
		},
		{
			name: "01",
			args: args{`Sound test__km__mf__side-d15.mpg`},
			want: `Sound test__km__mf__side-d15.mpg`,
		},
		// {
		// 	name: "02",
		// 	args: args{`Toys. Big tits Princess Pumpkins____implants-solo__side-d21..mp4`},
		// 	want: `Toys. Big tits Princess Pumpkins____implants-solo__side-d21.mp4`,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanFilename(tt.args.old); got != tt.want {
				t.Errorf("cleanFilename() = \ngot  %q\nwant %q\n", got, tt.want)
			}
		})
	}
}
