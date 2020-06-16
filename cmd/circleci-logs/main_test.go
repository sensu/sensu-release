package main

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
)

func Test_initCobra(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initCobra()
		})
	}
}

func Test_postInitCommands(t *testing.T) {
	type args struct {
		commands []*cobra.Command
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			postInitCommands(tt.args.commands)
		})
	}
}

func Test_presetRequiredFlags(t *testing.T) {
	type args struct {
		cmd *cobra.Command
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			presetRequiredFlags(tt.args.cmd)
		})
	}
}

func TestFetcher_Slug(t *testing.T) {
	type fields struct {
		apiToken        string
		workflowID      string
		projectType     string
		projectUsername string
		projectRepo     string
		jobName         string
		outputPath      string
		outputDir       string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Fetcher{
				apiToken:        tt.fields.apiToken,
				workflowID:      tt.fields.workflowID,
				projectType:     tt.fields.projectType,
				projectUsername: tt.fields.projectUsername,
				projectRepo:     tt.fields.projectRepo,
				jobName:         tt.fields.jobName,
				outputPath:      tt.fields.outputPath,
				outputDir:       tt.fields.outputDir,
			}
			if got := f.Slug(); got != tt.want {
				t.Errorf("Fetcher.Slug() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFetcher_fetchURL(t *testing.T) {
	type fields struct {
		apiToken        string
		workflowID      string
		projectType     string
		projectUsername string
		projectRepo     string
		jobName         string
		outputPath      string
		outputDir       string
	}
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Fetcher{
				apiToken:        tt.fields.apiToken,
				workflowID:      tt.fields.workflowID,
				projectType:     tt.fields.projectType,
				projectUsername: tt.fields.projectUsername,
				projectRepo:     tt.fields.projectRepo,
				jobName:         tt.fields.jobName,
				outputPath:      tt.fields.outputPath,
				outputDir:       tt.fields.outputDir,
			}
			got, err := f.fetchURL(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetcher.fetchURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fetcher.fetchURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFetcher_FetchWorkflowJobs(t *testing.T) {
	type fields struct {
		apiToken        string
		workflowID      string
		projectType     string
		projectUsername string
		projectRepo     string
		jobName         string
		outputPath      string
		outputDir       string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *WorkflowJobs
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Fetcher{
				apiToken:        tt.fields.apiToken,
				workflowID:      tt.fields.workflowID,
				projectType:     tt.fields.projectType,
				projectUsername: tt.fields.projectUsername,
				projectRepo:     tt.fields.projectRepo,
				jobName:         tt.fields.jobName,
				outputPath:      tt.fields.outputPath,
				outputDir:       tt.fields.outputDir,
			}
			got, err := f.FetchWorkflowJobs()
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetcher.FetchWorkflowJobs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fetcher.FetchWorkflowJobs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFetcher_FetchJob(t *testing.T) {
	type fields struct {
		apiToken        string
		workflowID      string
		projectType     string
		projectUsername string
		projectRepo     string
		jobName         string
		outputPath      string
		outputDir       string
	}
	type args struct {
		jobNumber int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Job
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Fetcher{
				apiToken:        tt.fields.apiToken,
				workflowID:      tt.fields.workflowID,
				projectType:     tt.fields.projectType,
				projectUsername: tt.fields.projectUsername,
				projectRepo:     tt.fields.projectRepo,
				jobName:         tt.fields.jobName,
				outputPath:      tt.fields.outputPath,
				outputDir:       tt.fields.outputDir,
			}
			got, err := f.FetchJob(tt.args.jobNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetcher.FetchJob() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fetcher.FetchJob() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFetcher_FetchStepOutputs(t *testing.T) {
	type fields struct {
		apiToken        string
		workflowID      string
		projectType     string
		projectUsername string
		projectRepo     string
		jobName         string
		outputPath      string
		outputDir       string
	}
	type args struct {
		jobNumber  int64
		stepNumber int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []JobStepOutput
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := Fetcher{
				apiToken:        tt.fields.apiToken,
				workflowID:      tt.fields.workflowID,
				projectType:     tt.fields.projectType,
				projectUsername: tt.fields.projectUsername,
				projectRepo:     tt.fields.projectRepo,
				jobName:         tt.fields.jobName,
				outputPath:      tt.fields.outputPath,
				outputDir:       tt.fields.outputDir,
			}
			got, err := f.FetchStepOutputs(tt.args.jobNumber, tt.args.stepNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetcher.FetchStepOutputs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fetcher.FetchStepOutputs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newRootCmd(t *testing.T) {
	tests := []struct {
		name string
		want *cobra.Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newRootCmd(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newRootCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
