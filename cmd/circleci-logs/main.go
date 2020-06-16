package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// TODO(jk): make use of next page token to loop through pages
type WorkflowJobs struct {
	NextPageToken string        `json:"next_page_token"`
	Items         []WorkflowJob `json:"items"`
}

type WorkflowJob struct {
	JobNumber int64  `json:"job_number"`
	Name      string `json:"name"`
}

type Job struct {
	Steps []JobStep `json:"steps"`
}

type JobStep struct {
	Name    string          `json:"name"`
	Actions []JobStepAction `json:"actions"`
}

type JobStepAction struct {
	StepNumber int64 `json:"step"`
}

type JobStepOutput struct {
	Message    string `json:"message"`
	Time       string `json:"time"`
	OutputType string `json:"type"`
}

type Fetcher struct {
	apiToken        string
	workflowID      string
	projectType     string
	projectUsername string
	projectRepo     string
	jobName         string
	outputPath      string
	outputDir       string
}

const (
	flagAuthToken             = "api-token"
	flagWorkflowID            = "workflow-id"
	flagProjectType           = "project-type"
	flagProjectUsername       = "project-username"
	flagProjectRepositoryName = "project-reponame"
	flagJobName               = "job-name"
	flagOutputPath            = "output-path"
	flagOutputDir             = "output-dir"

	defaultProjectType = "github"
	defaultOutputPath  = "log.txt"

	jobStepSeparator = "========================================\n"
)

var (
	rootCmd = newRootCmd()
)

func initCobra() {
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	postInitCommands(rootCmd.Commands())
}

func postInitCommands(commands []*cobra.Command) {
	for _, cmd := range commands {
		presetRequiredFlags(cmd)
		if cmd.HasSubCommands() {
			postInitCommands(cmd.Commands())
		}
	}
}

func presetRequiredFlags(cmd *cobra.Command) {
	viper.BindPFlags(cmd.Flags())
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if viper.IsSet(f.Name) && viper.GetString(f.Name) != "" {
			cmd.Flags().Set(f.Name, viper.GetString(f.Name))
		}
	})
}

func (f Fetcher) Slug() string {
	return fmt.Sprintf("%s/%s/%s", f.projectType, f.projectUsername, f.projectRepo)
}

func (f Fetcher) fetchURL(url string) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Circle-Token", f.apiToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (f Fetcher) FetchWorkflowJobs() (*WorkflowJobs, error) {
	url := fmt.Sprintf("https://circleci.com/api/v2/workflow/%s/job", f.workflowID)
	body, err := f.fetchURL(url)
	if err != nil {
		return nil, err
	}

	var workflowJobs WorkflowJobs
	if err := json.Unmarshal(body, &workflowJobs); err != nil {
		return nil, err
	}

	return &workflowJobs, nil
}

func (f Fetcher) FetchJob(jobNumber int64) (*Job, error) {
	url := fmt.Sprintf("https://circleci.com/api/v1.1/project/%s/%d", f.Slug(), jobNumber)
	body, err := f.fetchURL(url)
	if err != nil {
		return nil, err
	}

	var job Job
	if err := json.Unmarshal(body, &job); err != nil {
		return nil, err
	}

	return &job, nil
}

func (f Fetcher) FetchStepOutputs(jobNumber, stepNumber int64) ([]JobStepOutput, error) {
	url := fmt.Sprintf("https://circleci.com/api/v1.1/project/%s/%d/output/%d/0", f.Slug(), jobNumber, stepNumber)
	body, err := f.fetchURL(url)
	if err != nil {
		return nil, err
	}

	var stepOutputs []JobStepOutput
	if err := json.Unmarshal(body, &stepOutputs); err != nil {
		return nil, err
	}

	return stepOutputs, nil
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "circleci-fetch-logs",
		Short: "fetch logs for circleci jobs",
	}

	cmd.PersistentFlags().StringP(flagAuthToken, "t", "", "circleci api token")
	cmd.AddCommand(newFetchCmd())
	cmd.AddCommand(newFetchAllCmd())

	return cmd
}

func main() {
	cobra.OnInitialize(initCobra)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
