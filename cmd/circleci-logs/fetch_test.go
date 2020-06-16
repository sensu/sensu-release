package main

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func Test_executeFetchLogs(t *testing.T) {
	type args struct {
		cmd  *cobra.Command
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := executeFetchLogs(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("executeFetchLogs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newFetchCmd(t *testing.T) {
	tests := []struct {
		name string
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newFetchCmd(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newFetchCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}
