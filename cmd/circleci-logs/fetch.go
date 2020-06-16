package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func executeFetchLogs(cmd *cobra.Command, args []string) error {
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
		return err
	}
	defer logFile.Close()

	workflowJobs, err := fetcher.FetchWorkflowJobs()
	if err != nil {
		return err
	}

	var jobNumber int64
	for _, workflowJob := range workflowJobs.Items {
		if workflowJob.Name == fetcher.jobName {
			jobNumber = workflowJob.JobNumber
			break
		}
	}

	if jobNumber == 0 {
		return fmt.Errorf("a job with the name \"%s\" was not found in workflow \"%s\"",
			fetcher.jobName, fetcher.workflowID)
	}

	job, err := fetcher.FetchJob(jobNumber)
	if err != nil {
		return err
	}

	for _, step := range job.Steps {
		for _, action := range step.Actions {
			stepOutputs, err := fetcher.FetchStepOutputs(jobNumber, action.StepNumber)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error fetching output for step \"%d\" in job \"%d\"", action.StepNumber, jobNumber)
				continue
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

	return nil
}

func newFetchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fetch",
		Short: "fetch logs for a given job name",
		RunE:  executeFetchLogs,
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
