package main

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func Test_executeFetchAllLogs(t *testing.T) {
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
			if err := executeFetchAllLogs(tt.args.cmd, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("executeFetchAllLogs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_newFetchAllCmd(t *testing.T) {
	tests := []struct {
		name string
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newFetchAllCmd(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newFetchAllCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}
