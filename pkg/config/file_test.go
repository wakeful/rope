package config

import (
	"reflect"
	"rope/pkg/helpers"
	"testing"
)

var inputFile = []byte(`services:
  nginx: 3
  ubuntu: 1
`)

func TestLoadFileError(t *testing.T) {
	tests := []struct {
		name       string
		workingDir string
		want       []byte
		wantErr    bool
	}{
		{
			workingDir: "../../",
			want:       inputFile,
			wantErr:    false,
		},
		{
			workingDir: "/non_existing/path",
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if "" != tt.workingDir {
				helpers.ProjectDir = tt.workingDir
			}
			got, err := LoadFile()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    File
		wantErr bool
	}{
		{
			input:   []byte(`fizzBuzz`),
			wantErr: true,
		},
		{
			input: inputFile,
			want: File{
				Services: map[string]int{
					"nginx":  3,
					"ubuntu": 1,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseConfig(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}
