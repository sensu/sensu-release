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
}

const (
	flagAuthToken             = "api-token"
	flagWorkflowID            = "workflow-id"
	flagProjectType           = "project-type"
	flagProjectUsername       = "project-username"
	flagProjectRepositoryName = "project-reponame"
	flagJobName               = "job-name"
	flagOutputPath            = "output-path"

	defaultProjectType = "github"
	defaultOutputPath  = "log.txt"

	jobStepSeparator = "========================================\n"
)

var (
	rootCmd = newRootCmd()
)

func er(msg interface{}) {
	fmt.Println("error:", msg)
	os.Exit(1)
}

func appendNewline(s string) string {
	return fmt.Sprintf("%s\n", s)
}

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

func newFetchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fetch",
		Short: "fetch logs for a given job name",
		Run:   executeFetchLogs,
	}

	cmd.Flags().StringP(flagWorkflowID, "w", "", "id of the workflow that the job belongs to")
	cmd.Flags().StringP(flagProjectType, "p", defaultProjectType, "type of the project that the organization exists within (github or bitbucket)")
	cmd.Flags().StringP(flagProjectUsername, "u", "", "github or bitbucket username of the project")
	cmd.Flags().StringP(flagProjectRepositoryName, "r", "", "name of the repository that the job belongs to")
	cmd.Flags().StringP(flagJobName, "j", "", "name of the job to find logs for")
	cmd.Flags().StringP(flagOutputPath, "o", defaultOutputPath, "path to a file to store the retrieved logs in")

	cmd.MarkFlagRequired(flagAuthToken)
	cmd.MarkFlagRequired(flagWorkflowID)
	cmd.MarkFlagRequired(flagProjectUsername)
	cmd.MarkFlagRequired(flagProjectRepositoryName)
	cmd.MarkFlagRequired(flagJobName)

	viper.BindEnv(flagAuthToken, "CIRCLE_API_TOKEN")
	viper.BindEnv(flagWorkflowID, "CIRCLE_WORKFLOW_ID")
	viper.BindEnv(flagProjectType, "CIRCLE_PROJECT_TYPE")
	viper.BindEnv(flagProjectUsername, "CIRCLE_PROJECT_USERNAME")
	viper.BindEnv(flagProjectRepositoryName, "CIRCLE_PROJECT_REPONAME")
	viper.BindEnv(flagJobName, "CIRCLE_JOB")
	viper.BindEnv(flagOutputPath, "LOG_OUTPUT_PATH")

	viper.BindPFlag(flagAuthToken, cmd.PersistentFlags().Lookup(flagAuthToken))
	viper.BindPFlag(flagWorkflowID, cmd.Flags().Lookup(flagWorkflowID))
	viper.BindPFlag(flagProjectType, cmd.Flags().Lookup(flagProjectType))
	viper.BindPFlag(flagProjectUsername, cmd.Flags().Lookup(flagProjectUsername))
	viper.BindPFlag(flagProjectRepositoryName, cmd.Flags().Lookup(flagProjectRepositoryName))
	viper.BindPFlag(flagJobName, cmd.Flags().Lookup(flagJobName))
	viper.BindPFlag(flagOutputPath, cmd.Flags().Lookup(flagOutputPath))

	return cmd
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

func executeFetchLogs(cmd *cobra.Command, args []string) {
	fetcher := Fetcher{
		apiToken:        viper.GetString(flagAuthToken),
		workflowID:      viper.GetString(flagWorkflowID),
		projectType:     viper.GetString(flagProjectType),
		projectUsername: viper.GetString(flagProjectUsername),
		projectRepo:     viper.GetString(flagProjectRepositoryName),
		jobName:         viper.GetString(flagJobName),
		outputPath:      viper.GetString(flagOutputPath),
	}

	logFile, err := os.Create(fetcher.outputPath)
	if err != nil {
		er(err)
	}
	defer logFile.Close()

	workflowJobs, err := fetcher.FetchWorkflowJobs()
	if err != nil {
		er(err)
	}

	var jobNumber int64
	for _, workflowJob := range workflowJobs.Items {
		if workflowJob.Name == fetcher.jobName {
			jobNumber = workflowJob.JobNumber
			break
		}
	}

	if jobNumber == 0 {
		er(fmt.Errorf("a job with the name \"%s\" was not found in workflow \"%s\"", fetcher.jobName, fetcher.workflowID))
	}

	job, err := fetcher.FetchJob(jobNumber)
	if err != nil {
		er(err)
	}

	for _, step := range job.Steps {
		for _, action := range step.Actions {
			stepOutputs, err := fetcher.FetchStepOutputs(jobNumber, action.StepNumber)
			if err != nil {
				// no-op for now
			}
			stepStr := fmt.Sprintf("  Step: %s\n", step.Name)
			logFile.WriteString(jobStepSeparator)
			logFile.WriteString(stepStr)
			logFile.WriteString(jobStepSeparator)
			for _, stepOutput := range stepOutputs {
				timeStr := fmt.Sprintf("--- %s ---\n", stepOutput.Time)
				logFile.WriteString(timeStr)
				logFile.WriteString(stepOutput.Message)
				logFile.WriteString("\n")
			}
			logFile.Sync()
		}
	}

	return
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "circleci-fetch-logs",
		Short: "fetch logs for circleci jobs",
	}

	cmd.PersistentFlags().StringP(flagAuthToken, "t", "", "circleci api token")
	cmd.AddCommand(newFetchCmd())

	return cmd
}

func main() {
	cobra.OnInitialize(initCobra)
	if err := rootCmd.Execute(); err != nil {
		er(err)
	}
}
