package imgcontent

import (
	"reflect"
	"testing"
)

func TestNewFilename(t *testing.T) {
	type args struct {
		raw string
	}
	tests := []struct {
		name      string
		args      args
		wantFname *Filename
		wantErr   bool
	}{
		{
			name:      "filename will be trimmed",
			args:      args{raw: " filename "},
			wantFname: &Filename{Value: "filename", Valid: true},
			wantErr:   false,
		},
		{
			name:      "white space will be replaced with '_'",
			args:      args{raw: "file name"},
			wantFname: &Filename{Value: "file_name", Valid: true},
			wantErr:   false,
		},
		{
			name:      "sharp will be replaced with '_'",
			args:      args{raw: "file#name#"},
			wantFname: &Filename{Value: "file_name_", Valid: true},
			wantErr:   false,
		},
		{
			name:    "err must return if file name contains CR",
			args:    args{raw: "filenamewith\r"},
			wantErr: true,
		},
		{
			name:    "err must return if file name contains LF",
			args:    args{raw: "filenamewith\n"},
			wantErr: true,
		},
		{
			name:    "err must return if file name is empty",
			args:    args{raw: ""},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFname, err := NewFilename(tt.args.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewFilename() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFname, tt.wantFname) {
				t.Errorf("NewFilename() = %v, want %v", gotFname, tt.wantFname)
			}
		})
	}
}
