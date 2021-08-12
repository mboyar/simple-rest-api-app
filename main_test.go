package main

import (
	"reflect"
	"testing"
)

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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
