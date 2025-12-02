// Package actioncontext exposes strongly typed GitHub Actions runtime context
// derived from environment variables.
package actioncontext

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/google/go-github/v79/github"
	"github.com/spf13/afero"

	"github.com/ajbeck/actions-toolkit-go/core"
)

// The exported variables are mapped to a GitHub Actions default environment variable.
// The details for each variable can be found here:
// https://docs.github.com/en/actions/reference/workflows-and-actions/variables#default-environment-variables
var (
	// GithubCI indicates whether the runner is executing under GitHub Actions (CI).
	// During a GitHub Actions workflow run, this is always true.
	// Environment Variable: CI
	GithubCI *bool

	// The name of the action currently running, or the id of a step. For example, for an action, __repo-owner_name-of-action-repo.
	// Environment Variable: GITHUB_ACTION
	Action *string

	// The path where an action is located. This property is only supported in composite actions.
	// Environment Variable: GITHUB_ACTION_PATH
	ActionPath *string

	// For a step executing an action, this is the owner and repository name of the action.
	// Environment Variable: GITHUB_ACTION_REPOSITORY
	ActionRepository *string

	// Always set to true when GitHub Actions is running the workflow.
	// You can use this variable to differentiate when tests are being run locally or by GitHub Actions.
	// Environment Variable: GITHUB_ACTIONS
	Actions *bool

	// The name of the person or app that initiated the workflow.
	// Environment Variable: GITHUB_ACTOR
	Actor *string

	// The account ID of the person or app that triggered the initial workflow run.
	// Environment Variable: GITHUB_ACTOR_ID
	ActorId *string

	// The GitHub REST API base URL
	// Environment Variable: GITHUB_API_URL
	ApiUrl *string

	// The name of the base ref or target branch of the pull request in a workflow run.
	// This is only set when the event that triggers a workflow run is either pull_request or pull_request_target.
	// Environment Variable: GITHUB_BASE_REF
	BaseRef *string

	// The path on the runner to the file that sets variables from workflow commands.
	// The path to this file is unique to the current step and changes for each step in a job.
	// Environment Variable: GITHUB_ENV
	EnvFilePath *string

	// The name of the event that triggered the workflow.
	// Environment Variable: GITHUB_EVENT_NAME
	EventName *string

	// The path to the file on the runner that contains the full event webhook payload.
	// Environment Variable: GITHUB_EVENT_PATH
	EventPayloadFilePath *string

	// The GitHub GraphQL API base URL.
	// Environment Variable: GITHUB_GRAPHQL_URL
	GraphqlUrl *string

	// The head ref or source branch of the pull request in a workflow run.
	// This property is only set when the event that triggers a workflow run is either pull_request or pull_request_target.
	// Environment Variable: GITHUB_HEAD_REF
	HeadRef *string

	// The Job ID of the current job.
	// Environment Variable: GITHUB_JOB
	JobId *string

	// The path on the runner to the file that sets the current step's outputs from workflow commands.
	// The path to this file is unique to the current step and changes for each step in a job.
	// Environment Variable: GITHUB_OUTPUT
	OutputFilePath *string

	// The path on the runner to the file that sets system PATH variables from workflow commands.
	// The path to this file is unique to the current step and changes for each step in a job.
	// Environment Variable: GITHUB_PATH
	PathFilePath *string

	// The fully-formed ref of the branch or tag that triggered the workflow run.
	// For workflows triggered by push, this is the branch or tag ref that was pushed.
	// For workflows triggered by pull_request, this is the pull request merge branch.
	// For workflows triggered by release, this is the release tag created.
	// For other triggers, this is the branch or tag ref that triggered the workflow run.
	// This is only set if a branch or tag is available for the event type.
	// The ref given is fully-formed, meaning that for branches the format is refs/heads/<branch_name>.
	// For pull requests events except pull_request_target, it is refs/pull/<pr_number>/merge
	// pull_request_target events have the ref from the base branch. For tags it is refs/tags/<tag_name>.
	// Environment Variable: GITHUB_REF
	Ref *string

	// The short ref name of the branch or tag that triggered the workflow run.
	// This value matches the branch or tag name shown on GitHub.
	// Environment Variable: GITHUB_REF_NAME
	RefName *string

	// true if branch protections or rulesets are configured for the ref that triggered the workflow run.
	// Environment Variable: GITHUB_REF_PROTECTED
	RefProtected *bool

	// The type of ref that triggered the workflow run. Valid values are branch or tag.
	// Environment Variable: GITHUB_REF_TYPE
	RefType *string

	// The owner and repository name.
	// Environment Variable: GITHUB_REPOSITORY
	Repository *string

	// The ID of the repository.
	// Environment Variable: GITHUB_REPOSITORY_ID
	RepositoryId *string

	// The repository owner's name.
	// Environment Variable: GITHUB_REPOSITORY_OWNER
	RepositoryOwner *string

	// The repository owner's account ID.
	// Environment Variable: GITHUB_REPOSITORY_OWNER_ID
	RepositoryOwnerId *string

	// The number of days that workflow run logs and artifacts are kept.
	// Environment Variable: GITHUB_RETENTION_DAYS
	RetentionDays *int

	// A unique number for each attempt of a particular workflow run in a repository.
	// This number begins at 1 for the workflow run's first attempt, and increments with each re-run.
	// Environment Variable: GITHUB_RUN_ATTEMPT
	RunAttempt *int

	// A unique number for each workflow run within a repository.
	// This number does not change if you re-run the workflow run.
	// Environment Variable: GITHUB_RUN_ID
	RunId *int

	// A unique number for each run of a particular workflow in a repository.
	// This number begins at 1 for the workflow's first run, and increments with each new run.
	// This number does not change if you re-run the workflow run.
	// Environment Variable: GITHUB_RUN_NUMBER
	RunNumber *int

	// The URL of the GitHub server.
	// Environment Variable: GITHUB_SERVER_URL
	ServerUrl *string

	// The commit SHA that triggered the workflow.
	// The value of this commit SHA depends on the event that triggered the workflow.
	// Environment Variable: GITHUB_SHA
	Sha *string

	// The path on the runner to the file that contains job summaries from workflow commands.
	// The path to this file is unique to the current step and changes for each step in a job.
	// Environment Variable: GITHUB_STEP_SUMMARY
	StepSummaryFilePath *string

	// The username of the user that initiated the workflow run.
	// If the workflow run is a re-run, this value may differ from github.actor.
	// Any workflow re-runs will use the privileges of github.actor,
	// even if the actor initiating the re-run (github.triggering_actor) has different privileges.
	// Environment Variable: GITHUB_TRIGGERING_ACTOR
	TriggeringActor *string

	// The name of the workflow.
	// Environment Variable: GITHUB_WORKFLOW
	Workflow *string

	// The ref path to the workflow.
	// Environment Variable: GITHUB_WORKFLOW_REF
	WorkflowRef *string

	// The commit SHA for the workflow file.
	// Environment Variable: GITHUB_WORKFLOW_SHA
	WorkflowSha *string

	// The default working directory on the runner for steps,
	// and the default location of your repository when using the checkout action.
	// Environment Variable: GITHUB_WORKSPACE
	WorkspacePath *string

	// The architecture of the runner executing the job. Possible values are X86, X64, ARM, or ARM64.
	// Environment Variable: RUNNER_ARCH
	RunnerArch *string

	// This is set only if debug logging is enabled, and always has the value of 1.
	// It can be useful as an indicator to enable additional debugging or verbose logging in your own job steps.
	// Environment Variable: RUNNER_DEBUG
	RunnerDebug *bool

	// The environment of the runner executing the job.
	// Possible values are: github-hosted for GitHub-hosted runners provided by GitHub,
	// and self-hosted for self-hosted runners configured by the repository owner.
	// Environment Variable: RUNNER_ENVIRONMENT
	RunnerEnvironment *string

	// The name of the runner executing the job.
	// This name may not be unique in a workflow run as runners at the repository
	// and organization levels could use the same name.
	// Environment Variable: RUNNER_NAME
	RunnerName *string

	// The operating system of the runner executing the job.
	// Possible values are Linux, Windows, or macOS.
	// Environment Variable: RUNNER_OS
	RunnerOs *string

	// The path to a temporary directory on the runner.
	// This directory is emptied at the beginning and end of each job.
	// Note that files will not be removed if the runner's user account does not have permission to delete them.
	// Environment Variable: RUNNER_TEMP
	RunnerTempPath *string

	// The path to the directory containing preinstalled tools for GitHub-hosted runners.
	// Environment Variable: RUNNER_TOOL_CACHE
	RunnerToolCachePath *string

	// The typed event payload object parsed from the event file.
	// The type depends on the event name.
	EventPayload interface{}
)

func init() {
	loadActionContext(afero.NewOsFs(), os.LookupEnv)
}

type lookupEnvFunc func(key string) (string, bool)

// loadActionContext loads action context from environment variables and is provided for testing purposes.
// The provided lookupEnv function is used to retrieve environment variable values.
// The provided afero.Fs is used to read the event payload file.
func loadActionContext(fs afero.Fs, lookupEnv lookupEnvFunc) {
	EventPayload = nil

	Action = lookupToStrPtr(lookupEnv("GITHUB_ACTION"))
	ActionPath = lookupToStrPtr(lookupEnv("GITHUB_ACTION_PATH"))
	ActionRepository = lookupToStrPtr(lookupEnv("GITHUB_ACTION_REPOSITORY"))
	Actions = lookupToBoolPtr(lookupEnv("GITHUB_ACTIONS"))
	Actor = lookupToStrPtr(lookupEnv("GITHUB_ACTOR"))
	ActorId = lookupToStrPtr(lookupEnv("GITHUB_ACTOR_ID"))
	ApiUrl = lookupToStrPtr(lookupEnv("GITHUB_API_URL"))
	BaseRef = lookupToStrPtr(lookupEnv("GITHUB_BASE_REF"))
	EnvFilePath = lookupToStrPtr(lookupEnv("GITHUB_ENV"))
	EventName = lookupToStrPtr(lookupEnv("GITHUB_EVENT_NAME"))
	EventPayloadFilePath = lookupToStrPtr(lookupEnv("GITHUB_EVENT_PATH"))
	GithubCI = lookupToBoolPtr(lookupEnv("CI"))
	GraphqlUrl = lookupToStrPtr(lookupEnv("GITHUB_GRAPHQL_URL"))
	HeadRef = lookupToStrPtr(lookupEnv("GITHUB_HEAD_REF"))
	JobId = lookupToStrPtr(lookupEnv("GITHUB_JOB"))
	OutputFilePath = lookupToStrPtr(lookupEnv("GITHUB_OUTPUT"))
	PathFilePath = lookupToStrPtr(lookupEnv("GITHUB_PATH"))
	Ref = lookupToStrPtr(lookupEnv("GITHUB_REF"))
	RefName = lookupToStrPtr(lookupEnv("GITHUB_REF_NAME"))
	RefProtected = lookupToBoolPtr(lookupEnv("GITHUB_REF_PROTECTED"))
	RefType = lookupToStrPtr(lookupEnv("GITHUB_REF_TYPE"))
	Repository = lookupToStrPtr(lookupEnv("GITHUB_REPOSITORY"))
	RepositoryId = lookupToStrPtr(lookupEnv("GITHUB_REPOSITORY_ID"))
	RepositoryOwner = lookupToStrPtr(lookupEnv("GITHUB_REPOSITORY_OWNER"))
	RepositoryOwnerId = lookupToStrPtr(lookupEnv("GITHUB_REPOSITORY_OWNER_ID"))
	RetentionDays = lookupToIntPtr(lookupEnv("GITHUB_RETENTION_DAYS"))
	RunAttempt = lookupToIntPtr(lookupEnv("GITHUB_RUN_ATTEMPT"))
	RunId = lookupToIntPtr(lookupEnv("GITHUB_RUN_ID"))
	RunNumber = lookupToIntPtr(lookupEnv("GITHUB_RUN_NUMBER"))
	ServerUrl = lookupToStrPtr(lookupEnv("GITHUB_SERVER_URL"))
	Sha = lookupToStrPtr(lookupEnv("GITHUB_SHA"))
	StepSummaryFilePath = lookupToStrPtr(lookupEnv("GITHUB_STEP_SUMMARY"))
	TriggeringActor = lookupToStrPtr(lookupEnv("GITHUB_TRIGGERING_ACTOR"))
	Workflow = lookupToStrPtr(lookupEnv("GITHUB_WORKFLOW"))
	WorkflowRef = lookupToStrPtr(lookupEnv("GITHUB_WORKFLOW_REF"))
	WorkflowSha = lookupToStrPtr(lookupEnv("GITHUB_WORKFLOW_SHA"))
	WorkspacePath = lookupToStrPtr(lookupEnv("GITHUB_WORKSPACE"))

	RunnerArch = lookupToStrPtr(lookupEnv("RUNNER_ARCH"))
	RunnerDebug = lookupToBoolPtr(lookupEnv("RUNNER_DEBUG"))
	RunnerEnvironment = lookupToStrPtr(lookupEnv("RUNNER_ENVIRONMENT"))
	RunnerName = lookupToStrPtr(lookupEnv("RUNNER_NAME"))
	RunnerOs = lookupToStrPtr(lookupEnv("RUNNER_OS"))
	RunnerTempPath = lookupToStrPtr(lookupEnv("RUNNER_TEMP"))
	RunnerToolCachePath = lookupToStrPtr(lookupEnv("RUNNER_TOOL_CACHE"))

	if EventName != nil && EventPayloadFilePath != nil {
		payloadBytes, err := afero.ReadFile(fs, *EventPayloadFilePath)

		if err != nil {
			EventPayload = nil

			core.Debug(err.Error())

			return
		}

		switch *EventName {
		case "branch_protection_rule":
			EventPayload = &github.BranchProtectionRuleEvent{}
		case "check_run":
			EventPayload = &github.CheckRunEvent{}
		case "check_suite":
			EventPayload = &github.CheckSuiteEvent{}
		case "create":
			EventPayload = &github.CreateEvent{}
		case "delete":
			EventPayload = &github.DeleteEvent{}
		case "deployment":
			EventPayload = &github.DeploymentEvent{}
		case "deployment_status":
			EventPayload = &github.DeploymentStatusEvent{}
		case "discussion":
			EventPayload = &github.DiscussionEvent{}
		case "discussion_comment":
			EventPayload = &github.DiscussionCommentEvent{}
		case "fork":
			EventPayload = &github.ForkEvent{}
		case "gollum":
			EventPayload = &github.GollumEvent{}
		case "package":
			EventPayload = &github.PackageEvent{}
		case "issue_comment":
			EventPayload = &github.IssueCommentEvent{}
		case "issues":
			EventPayload = &github.IssuesEvent{}
		case "label":
			EventPayload = &github.LabelEvent{}
		case "merge_group":
			EventPayload = &github.MergeGroupEvent{}
		case "milestone":
			EventPayload = &github.MilestoneEvent{}
		case "page_build":
			EventPayload = &github.PageBuildEvent{}
		case "public":
			EventPayload = &github.PublicEvent{}
		case "pull_request":
			EventPayload = &github.PullRequestEvent{}
		case "pull_request_comment":
			EventPayload = &github.IssueCommentEvent{}
		case "pull_request_review":
			EventPayload = &github.PullRequestReviewEvent{}
		case "pull_request_review_comment":
			EventPayload = &github.PullRequestReviewCommentEvent{}
		case "pull_request_target":
			EventPayload = &github.PullRequestTargetEvent{}
		case "push":
			EventPayload = &github.PushEvent{}
		case "registry_package":
			EventPayload = &github.RegistryPackageEvent{}
		case "release":
			EventPayload = &github.ReleaseEvent{}
		case "repository":
			EventPayload = &github.RepositoryEvent{}
		case "repository_dispatch":
			EventPayload = &github.RepositoryDispatchEvent{}
		case "status":
			EventPayload = &github.StatusEvent{}
		case "watch":
			EventPayload = &github.WatchEvent{}
		case "workflow_dispatch":
			EventPayload = &github.WorkflowDispatchEvent{}
		case "workflow_run":
			EventPayload = &github.WorkflowRunEvent{}
		default:
			EventPayload = &map[string]interface{}{}
		}

		err = json.Unmarshal(payloadBytes, EventPayload)

		if err != nil {
			EventPayload = nil

			core.Debug(err.Error())

			return
		}

	}
}

// lookupToStrPtr a helper functions to convert lookupEnv results to a string pointer
func lookupToStrPtr(v string, isSet bool) *string {
	if isSet {
		return &v
	}
	return nil
}

// lookupToIntPtr a helper functions to convert lookupEnv results to an int pointer
func lookupToIntPtr(v string, isSet bool) *int {
	if isSet {
		vI, err := strconv.Atoi(v)

		if err != nil {
			return nil
		}

		return &vI
	}
	return nil
}

// lookupToBoolPtr a helper functions to convert lookupEnv results to a bool pointer
func lookupToBoolPtr(v string, isSet bool) *bool {
	if isSet {
		vB, err := strconv.ParseBool(v)

		if err != nil {
			return nil
		}

		return &vB
	}
	return nil
}
