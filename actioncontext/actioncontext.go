// Package actioncontext exposes strongly typed GitHub Actions runtime context
// derived from environment variables.
package actioncontext

import (
  "encoding/json"
  "fmt"
  "os"
  "strconv"
  "strings"

  "github.com/google/go-github/v79/github"
  "github.com/spf13/afero"

  "github.com/ajbeck/actions-toolkit-go/core"
)

type ActionContext struct {
  GithubCI             *bool       `json:"githubCI"`
  Action               *string     `json:"action"`
  ActionPath           *string     `json:"actionPath"`
  ActionRepository     *string     `json:"actionRepository"`
  Actions              *bool       `json:"actions"`
  Actor                *string     `json:"actor"`
  ActorId              *string     `json:"actorId"`
  ApiUrl               *string     `json:"apiUrl"`
  BaseRef              *string     `json:"baseRef"`
  EnvFilePath          *string     `json:"envFilePath"`
  EventName            *string     `json:"eventName"`
  EventPayloadFilePath *string     `json:"eventPayloadFilePath"`
  GraphqlUrl           *string     `json:"graphqlUrl"`
  HeadRef              *string     `json:"headRef"`
  JobId                *string     `json:"jobId"`
  OutputFilePath       *string     `json:"outputFilePath"`
  PathFilePath         *string     `json:"pathFilePath"`
  Ref                  *string     `json:"ref"`
  RefName              *string     `json:"refName"`
  RefProtected         *bool       `json:"refProtected"`
  RefType              *string     `json:"refType"`
  Repository           *string     `json:"repository"`
  RepositoryId         *string     `json:"repositoryId"`
  RepositoryOwner      *string     `json:"repositoryOwner"`
  RepositoryOwnerId    *string     `json:"repositoryOwnerId"`
  RetentionDays        *int        `json:"retentionDays"`
  RunAttempt           *int        `json:"runAttempt"`
  RunId                *int        `json:"runId"`
  RunNumber            *int        `json:"runNumber"`
  ServerUrl            *string     `json:"serverUrl"`
  Sha                  *string     `json:"sha"`
  StepSummaryFilePath  *string     `json:"stepSummaryFilePath"`
  TriggeringActor      *string     `json:"triggeringActor"`
  Workflow             *string     `json:"workflow"`
  WorkflowRef          *string     `json:"workflowRef"`
  WorkflowSha          *string     `json:"workflowSha"`
  WorkspacePath        *string     `json:"workspacePath"`
  RunnerArch           *string     `json:"runnerArch"`
  RunnerDebug          *bool       `json:"runnerDebug"`
  RunnerEnvironment    *string     `json:"runnerEnvironment"`
  RunnerName           *string     `json:"runnerName"`
  RunnerOs             *string     `json:"runnerOs"`
  RunnerTempPath       *string     `json:"runnerTempPath"`
  RunnerToolCachePath  *string     `json:"runnerToolCachePath"`
  EventPayload         interface{} `json:"eventPayload"`
}

// String returns a pretty-printed representation of the ActionContext with all fields.
// Fields that are nil are explicitly shown as "<nil>".
func (ac *ActionContext) String() string {
  if ac == nil {
    return "ActionContext: <nil>"
  }

  var sb strings.Builder
  sb.WriteString("ActionContext:\n")

  // Helper functions to format pointer values
  formatStr := func(name string, ptr *string) string {
    if ptr == nil {
      return fmt.Sprintf("  %-25s <nil>\n", name+":")
    }
    return fmt.Sprintf("  %-25s %s\n", name+":", *ptr)
  }

  formatBool := func(name string, ptr *bool) string {
    if ptr == nil {
      return fmt.Sprintf("  %-25s <nil>\n", name+":")
    }
    return fmt.Sprintf("  %-25s %t\n", name+":", *ptr)
  }

  formatInt := func(name string, ptr *int) string {
    if ptr == nil {
      return fmt.Sprintf("  %-25s <nil>\n", name+":")
    }
    return fmt.Sprintf("  %-25s %d\n", name+":", *ptr)
  }

  formatInterface := func(name string, val interface{}) string {
    if val == nil {
      return fmt.Sprintf("  %-25s <nil>\n", name+":")
    }
    return fmt.Sprintf("  %-25s %T\n", name+":", val)
  }

  sb.WriteString(formatBool("GithubCI", ac.GithubCI))
  sb.WriteString(formatStr("Action", ac.Action))
  sb.WriteString(formatStr("ActionPath", ac.ActionPath))
  sb.WriteString(formatStr("ActionRepository", ac.ActionRepository))
  sb.WriteString(formatBool("Actions", ac.Actions))
  sb.WriteString(formatStr("Actor", ac.Actor))
  sb.WriteString(formatStr("ActorId", ac.ActorId))
  sb.WriteString(formatStr("ApiUrl", ac.ApiUrl))
  sb.WriteString(formatStr("BaseRef", ac.BaseRef))
  sb.WriteString(formatStr("EnvFilePath", ac.EnvFilePath))
  sb.WriteString(formatStr("EventName", ac.EventName))
  sb.WriteString(formatStr("EventPayloadFilePath", ac.EventPayloadFilePath))
  sb.WriteString(formatStr("GraphqlUrl", ac.GraphqlUrl))
  sb.WriteString(formatStr("HeadRef", ac.HeadRef))
  sb.WriteString(formatStr("JobId", ac.JobId))
  sb.WriteString(formatStr("OutputFilePath", ac.OutputFilePath))
  sb.WriteString(formatStr("PathFilePath", ac.PathFilePath))
  sb.WriteString(formatStr("Ref", ac.Ref))
  sb.WriteString(formatStr("RefName", ac.RefName))
  sb.WriteString(formatBool("RefProtected", ac.RefProtected))
  sb.WriteString(formatStr("RefType", ac.RefType))
  sb.WriteString(formatStr("Repository", ac.Repository))
  sb.WriteString(formatStr("RepositoryId", ac.RepositoryId))
  sb.WriteString(formatStr("RepositoryOwner", ac.RepositoryOwner))
  sb.WriteString(formatStr("RepositoryOwnerId", ac.RepositoryOwnerId))
  sb.WriteString(formatInt("RetentionDays", ac.RetentionDays))
  sb.WriteString(formatInt("RunAttempt", ac.RunAttempt))
  sb.WriteString(formatInt("RunId", ac.RunId))
  sb.WriteString(formatInt("RunNumber", ac.RunNumber))
  sb.WriteString(formatStr("ServerUrl", ac.ServerUrl))
  sb.WriteString(formatStr("Sha", ac.Sha))
  sb.WriteString(formatStr("StepSummaryFilePath", ac.StepSummaryFilePath))
  sb.WriteString(formatStr("TriggeringActor", ac.TriggeringActor))
  sb.WriteString(formatStr("Workflow", ac.Workflow))
  sb.WriteString(formatStr("WorkflowRef", ac.WorkflowRef))
  sb.WriteString(formatStr("WorkflowSha", ac.WorkflowSha))
  sb.WriteString(formatStr("WorkspacePath", ac.WorkspacePath))
  sb.WriteString(formatStr("RunnerArch", ac.RunnerArch))
  sb.WriteString(formatBool("RunnerDebug", ac.RunnerDebug))
  sb.WriteString(formatStr("RunnerEnvironment", ac.RunnerEnvironment))
  sb.WriteString(formatStr("RunnerName", ac.RunnerName))
  sb.WriteString(formatStr("RunnerOs", ac.RunnerOs))
  sb.WriteString(formatStr("RunnerTempPath", ac.RunnerTempPath))
  sb.WriteString(formatStr("RunnerToolCachePath", ac.RunnerToolCachePath))
  sb.WriteString(formatInterface("EventPayload", ac.EventPayload))

  return sb.String()
}

// The exported variables are mapped to a GitHub Actions default environment variable.
// The details for each variable can be found here:
// https://docs.github.com/en/actions/reference/workflows-and-actions/variables#default-environment-variables
var (
  // CurrentActionContext
  CurrentActionContext *ActionContext

  // GithubCI indicates whether the runner is executing under GitHub Actions (CI).
  // During a GitHub Actions workflow run, this is always true.
  // Environment Variable: CI
  GithubCI *bool

  // Action is the name of the action currently running, or the id of a step. For example, for an action, __repo-owner_name-of-action-repo.
  // Environment Variable: GITHUB_ACTION
  Action *string

  // ActionPath is the path where an action is located. This property is only supported in composite actions.
  // Environment Variable: GITHUB_ACTION_PATH
  ActionPath *string

  // ActionRepository is, for a step executing an action, the owner and repository name of the action.
  // Environment Variable: GITHUB_ACTION_REPOSITORY
  ActionRepository *string

  // Actions is always set to true when GitHub Actions is running the workflow.
  // You can use this variable to differentiate when tests are being run locally or by GitHub Actions.
  // Environment Variable: GITHUB_ACTIONS
  Actions *bool

  // Actor is the name of the person or app that initiated the workflow.
  // Environment Variable: GITHUB_ACTOR
  Actor *string

  // ActorId is the account ID of the person or app that triggered the initial workflow run.
  // Environment Variable: GITHUB_ACTOR_ID
  ActorId *string

  // ApiUrl is the GitHub REST API base URL
  // Environment Variable: GITHUB_API_URL
  ApiUrl *string

  // BaseRef is the name of the base ref or target branch of the pull request in a workflow run.
  // This is only set when the event that triggers a workflow run is either pull_request or pull_request_target.
  // Environment Variable: GITHUB_BASE_REF
  BaseRef *string

  // EnvFilePath is the path on the runner to the file that sets variables from workflow commands.
  // The path to this file is unique to the current step and changes for each step in a job.
  // Environment Variable: GITHUB_ENV
  EnvFilePath *string

  // EventName is the name of the event that triggered the workflow.
  // Environment Variable: GITHUB_EVENT_NAME
  EventName *string

  // EventPayloadFilePath is the path to the file on the runner that contains the full event webhook payload.
  // Environment Variable: GITHUB_EVENT_PATH
  EventPayloadFilePath *string

  // GraphqlUrl is the GitHub GraphQL API base URL.
  // Environment Variable: GITHUB_GRAPHQL_URL
  GraphqlUrl *string

  // HeadRef is the head ref or source branch of the pull request in a workflow run.
  // This property is only set when the event that triggers a workflow run is either pull_request or pull_request_target.
  // Environment Variable: GITHUB_HEAD_REF
  HeadRef *string

  // JobId is the Job ID of the current job.
  // Environment Variable: GITHUB_JOB
  JobId *string

  // OutputFilePath is the path on the runner to the file that sets the current step's outputs from workflow commands.
  // The path to this file is unique to the current step and changes for each step in a job.
  // Environment Variable: GITHUB_OUTPUT
  OutputFilePath *string

  // PathFilePath is the path on the runner to the file that sets system PATH variables from workflow commands.
  // The path to this file is unique to the current step and changes for each step in a job.
  // Environment Variable: GITHUB_PATH
  PathFilePath *string

  // Ref is the fully-formed ref of the branch or tag that triggered the workflow run.
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

  // RefName is the short ref name of the branch or tag that triggered the workflow run.
  // This value matches the branch or tag name shown on GitHub.
  // Environment Variable: GITHUB_REF_NAME
  RefName *string

  // RefProtected is true if branch protections or rulesets are configured for the ref that triggered the workflow run.
  // Environment Variable: GITHUB_REF_PROTECTED
  RefProtected *bool

  // RefType is the type of ref that triggered the workflow run. Valid values are branch or tag.
  // Environment Variable: GITHUB_REF_TYPE
  RefType *string

  // Repository is the owner and repository name.
  // Environment Variable: GITHUB_REPOSITORY
  Repository *string

  // RepositoryId is the ID of the repository.
  // Environment Variable: GITHUB_REPOSITORY_ID
  RepositoryId *string

  // RepositoryOwner is the repository owner's name.
  // Environment Variable: GITHUB_REPOSITORY_OWNER
  RepositoryOwner *string

  // RepositoryOwnerId is the repository owner's account ID.
  // Environment Variable: GITHUB_REPOSITORY_OWNER_ID
  RepositoryOwnerId *string

  // RetentionDays is the number of days that workflow run logs and artifacts are kept.
  // Environment Variable: GITHUB_RETENTION_DAYS
  RetentionDays *int

  // RunAttempt is a unique number for each attempt of a particular workflow run in a repository.
  // This number begins at 1 for the workflow run's first attempt, and increments with each re-run.
  // Environment Variable: GITHUB_RUN_ATTEMPT
  RunAttempt *int

  // RunId is a unique number for each workflow run within a repository.
  // This number does not change if you re-run the workflow run.
  // Environment Variable: GITHUB_RUN_ID
  RunId *int

  // RunNumber is a unique number for each run of a particular workflow in a repository.
  // This number begins at 1 for the workflow's first run, and increments with each new run.
  // This number does not change if you re-run the workflow run.
  // Environment Variable: GITHUB_RUN_NUMBER
  RunNumber *int

  // ServerUrl is the URL of the GitHub server.
  // Environment Variable: GITHUB_SERVER_URL
  ServerUrl *string

  // Sha is the commit SHA that triggered the workflow.
  // The value of this commit SHA depends on the event that triggered the workflow.
  // Environment Variable: GITHUB_SHA
  Sha *string

  // StepSummaryFilePath is the path on the runner to the file that contains job summaries from workflow commands.
  // The path to this file is unique to the current step and changes for each step in a job.
  // Environment Variable: GITHUB_STEP_SUMMARY
  StepSummaryFilePath *string

  // TriggeringActor is the username of the user that initiated the workflow run.
  // If the workflow run is a re-run, this value may differ from github.actor.
  // Any workflow re-runs will use the privileges of github.actor,
  // even if the actor initiating the re-run (github.triggering_actor) has different privileges.
  // Environment Variable: GITHUB_TRIGGERING_ACTOR
  TriggeringActor *string

  // Workflow is the name of the workflow.
  // Environment Variable: GITHUB_WORKFLOW
  Workflow *string

  // WorkflowRef is the ref path to the workflow.
  // Environment Variable: GITHUB_WORKFLOW_REF
  WorkflowRef *string

  // WorkflowSha is the commit SHA for the workflow file.
  // Environment Variable: GITHUB_WORKFLOW_SHA
  WorkflowSha *string

  // WorkspacePath is the default working directory on the runner for steps,
  // and the default location of your repository when using the checkout action.
  // Environment Variable: GITHUB_WORKSPACE
  WorkspacePath *string

  // RunnerArch is the architecture of the runner executing the job. Possible values are X86, X64, ARM, or ARM64.
  // Environment Variable: RUNNER_ARCH
  RunnerArch *string

  // RunnerDebug is this is set only if debug logging is enabled, and always has the value of 1.
  // It can be useful as an indicator to enable additional debugging or verbose logging in your own job steps.
  // Environment Variable: RUNNER_DEBUG
  RunnerDebug *bool

  // RunnerEnvironment is the environment of the runner executing the job.
  // Possible values are: github-hosted for GitHub-hosted runners provided by GitHub,
  // and self-hosted for self-hosted runners configured by the repository owner.
  // Environment Variable: RUNNER_ENVIRONMENT
  RunnerEnvironment *string

  // RunnerName is the name of the runner executing the job.
  // This name may not be unique in a workflow run as runners at the repository
  // and organization levels could use the same name.
  // Environment Variable: RUNNER_NAME
  RunnerName *string

  // RunnerOs is the operating system of the runner executing the job.
  // Possible values are Linux, Windows, or macOS.
  // Environment Variable: RUNNER_OS
  RunnerOs *string

  // RunnerTempPath is the path to a temporary directory on the runner.
  // This directory is emptied at the beginning and end of each job.
  // Note that files will not be removed if the runner's user account does not have permission to delete them.
  // Environment Variable: RUNNER_TEMP
  RunnerTempPath *string

  // RunnerToolCachePath is the path to the directory containing preinstalled tools for GitHub-hosted runners.
  // Environment Variable: RUNNER_TOOL_CACHE
  RunnerToolCachePath *string

  // EventPayload is the typed event payload object parsed from the event file.
  // The type depends on the event name.
  EventPayload interface{}
)

func init() {
  err := loadActionContext(afero.NewOsFs(), os.LookupEnv)
  if err != nil {
    core.Debug("Error loading action context: " + err.Error())
  }

  if CurrentActionContext != nil {
    core.Debug(CurrentActionContext.String())
  }
}

type lookupEnvFunc func(key string) (string, bool)

// loadActionContext loads action context from environment variables and is provided for testing purposes.
// The provided lookupEnv function is used to retrieve environment variable values.
// The provided afero.Fs is used to read the event payload file.
// Returns an error if loading the event payload fails.
func loadActionContext(fs afero.Fs, lookupEnv lookupEnvFunc) error {
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
      return fmt.Errorf("failed to read event payload file: %w", err)
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
      return fmt.Errorf("failed to unmarshal event payload: %w", err)
    }
  }

  CurrentActionContext = &ActionContext{
    GithubCI:             GithubCI,
    Action:               Action,
    ActionPath:           ActionPath,
    ActionRepository:     ActionRepository,
    Actions:              Actions,
    Actor:                Actor,
    ActorId:              ActorId,
    ApiUrl:               ApiUrl,
    BaseRef:              BaseRef,
    EnvFilePath:          EnvFilePath,
    EventName:            EventName,
    EventPayloadFilePath: EventPayloadFilePath,
    GraphqlUrl:           GraphqlUrl,
    HeadRef:              HeadRef,
    JobId:                JobId,
    OutputFilePath:       OutputFilePath,
    PathFilePath:         PathFilePath,
    Ref:                  Ref,
    RefName:              RefName,
    RefProtected:         RefProtected,
    RefType:              RefType,
    Repository:           Repository,
    RepositoryId:         RepositoryId,
    RepositoryOwner:      RepositoryOwner,
    RepositoryOwnerId:    RepositoryOwnerId,
    RetentionDays:        RetentionDays,
    RunAttempt:           RunAttempt,
    RunId:                RunId,
    RunNumber:            RunNumber,
    ServerUrl:            ServerUrl,
    Sha:                  Sha,
    StepSummaryFilePath:  StepSummaryFilePath,
    TriggeringActor:      TriggeringActor,
    Workflow:             Workflow,
    WorkflowRef:          WorkflowRef,
    WorkflowSha:          WorkflowSha,
    WorkspacePath:        WorkspacePath,
    RunnerArch:           RunnerArch,
    RunnerDebug:          RunnerDebug,
    RunnerEnvironment:    RunnerEnvironment,
    RunnerName:           RunnerName,
    RunnerOs:             RunnerOs,
    RunnerTempPath:       RunnerTempPath,
    RunnerToolCachePath:  RunnerToolCachePath,
    EventPayload:         EventPayload,
  }

  return nil
}

// lookupToStrPtr is a helper function to convert lookupEnv results to a string pointer
func lookupToStrPtr(v string, isSet bool) *string {
  if isSet {
    return &v
  }
  return nil
}

// lookupToIntPtr is a helper function to convert lookupEnv results to an int pointer
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

// lookupToBoolPtr is a helper functions to convert lookupEnv results to a bool pointer
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
