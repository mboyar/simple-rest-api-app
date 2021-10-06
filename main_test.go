package main

import (
	"reflect"
	"testing"
)

var strTest1 string = "Startup duration of the system: Startup finished in 5.225s (kernel) + 2min 15.289s (userspace) = 2min 20.515s graphical.target reached after 2min 15.220s in userspace"

var g_buf []byte = []byte(strTest1)

func Test_responser(t *testing.T) {

	type args struct {
		duration []byte
		version  string
		msg      string
	}
	tests := []struct {
		name string
		args args
	}{
		{"test1", args{g_buf, "v2.0", "Hello golang rest api"}},
		{"test2", args{g_buf, "v0.2", "Mello golang fest zapi"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responser(tt.args.duration, tt.args.version, tt.args.msg)
		})
	}
}

func Test_getDuration(t *testing.T) {
	tests := []struct {
		name    string
		want    []byte
		wantErr bool
	}{
		{"test1", g_buf, true},
		{"test2", g_buf, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDuration()
			if (err != nil) != tt.wantErr {
				t.Errorf("getDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseSystemdAnalyze(t *testing.T) {
	type args struct {
		cmdStdout []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseSystemdAnalyze(tt.args.cmdStdout); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSystemdAnalyze() = %v, want %v", got, tt.want)
			}
		})
	}
}
