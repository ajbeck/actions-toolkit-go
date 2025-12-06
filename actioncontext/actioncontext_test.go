package actioncontext

import (
  "strings"
  "testing"

  "github.com/google/go-cmp/cmp"
  "github.com/google/go-github/v79/github"
  "github.com/spf13/afero"
)

func TestLookupToStrPtr(t *testing.T) {
  t.Parallel()

  tests := []struct {
    name   string
    input  string
    isSet  bool
    expect *string
  }{
    {name: "value set", input: "foo", isSet: true, expect: strPtr("foo")},
    {name: "value not set", input: "", isSet: false, expect: nil},
  }

  for _, tt := range tests {
    tt := tt
    t.Run(tt.name, func(t *testing.T) {
      t.Parallel()

      got := lookupToStrPtr(tt.input, tt.isSet)
      if diff := cmp.Diff(tt.expect, got); diff != "" {
        t.Fatalf("lookupToStrPtr() mismatch (-want +got):\n%s", diff)
      }
    })
  }
}

func TestLookupToIntPtr(t *testing.T) {
  t.Parallel()

  tests := []struct {
    name   string
    input  string
    isSet  bool
    expect *int
  }{
    {name: "valid int", input: "42", isSet: true, expect: intPtr(42)},
    {name: "invalid int", input: "x", isSet: true, expect: nil},
    {name: "value not set", input: "", isSet: false, expect: nil},
  }

  for _, tt := range tests {
    tt := tt
    t.Run(tt.name, func(t *testing.T) {
      t.Parallel()

      got := lookupToIntPtr(tt.input, tt.isSet)
      if diff := cmp.Diff(tt.expect, got); diff != "" {
        t.Fatalf("lookupToIntPtr() mismatch (-want +got):\n%s", diff)
      }
    })
  }
}

func TestLookupToBoolPtr(t *testing.T) {
  t.Parallel()

  tests := []struct {
    name   string
    input  string
    isSet  bool
    expect *bool
  }{
    {name: "true", input: "true", isSet: true, expect: boolPtr(true)},
    {name: "false", input: "false", isSet: true, expect: boolPtr(false)},
    {name: "invalid bool", input: "yes", isSet: true, expect: nil},
    {name: "value not set", input: "", isSet: false, expect: nil},
  }

  for _, tt := range tests {
    tt := tt
    t.Run(tt.name, func(t *testing.T) {
      t.Parallel()

      got := lookupToBoolPtr(tt.input, tt.isSet)
      if diff := cmp.Diff(tt.expect, got); diff != "" {
        t.Fatalf("lookupToBoolPtr() mismatch (-want +got):\n%s", diff)
      }
    })
  }
}

func TestLoadActionContextPopulatesValuesAndEventPayload(t *testing.T) {
  defer func() { EventPayload = nil }()

  fs := afero.NewMemMapFs()

  const payloadPath = "/event.json"
  err := afero.WriteFile(fs, payloadPath, []byte(`{"ref":"refs/heads/main"}`), 0o644)
  if err != nil {
    t.Fatalf("failed to write payload file: %v", err)
  }

  env := map[string]string{
    "GITHUB_ACTION":            "build",
    "GITHUB_ACTION_REPOSITORY": "owner/action",
    "GITHUB_ACTOR":             "octocat",
    "GITHUB_ACTIONS":           "true",
    "GITHUB_RUN_ID":            "123",
    "GITHUB_RUN_ATTEMPT":       "2",
    "GITHUB_REF_PROTECTED":     "false",
    "GITHUB_EVENT_NAME":        "push",
    "GITHUB_EVENT_PATH":        payloadPath,
    "RUNNER_DEBUG":             "1",
  }

  loadActionContext(fs, mockLookup(env))

  if diff := cmp.Diff(strPtr("build"), Action); diff != "" {
    t.Fatalf("Action mismatch (-want +got):\n%s", diff)
  }

  if diff := cmp.Diff(strPtr("owner/action"), ActionRepository); diff != "" {
    t.Fatalf("ActionRepository mismatch (-want +got):\n%s", diff)
  }

  if diff := cmp.Diff(strPtr("octocat"), Actor); diff != "" {
    t.Fatalf("Actor mismatch (-want +got):\n%s", diff)
  }

  if diff := cmp.Diff(boolPtr(true), Actions); diff != "" {
    t.Fatalf("Actions mismatch (-want +got):\n%s", diff)
  }

  if diff := cmp.Diff(intPtr(123), RunId); diff != "" {
    t.Fatalf("RunId mismatch (-want +got):\n%s", diff)
  }

  if diff := cmp.Diff(intPtr(2), RunAttempt); diff != "" {
    t.Fatalf("RunAttempt mismatch (-want +got):\n%s", diff)
  }

  if diff := cmp.Diff(boolPtr(false), RefProtected); diff != "" {
    t.Fatalf("RefProtected mismatch (-want +got):\n%s", diff)
  }

  if diff := cmp.Diff(boolPtr(true), RunnerDebug); diff != "" {
    t.Fatalf("RunnerDebug mismatch (-want +got):\n%s", diff)
  }

  pushPayload, ok := EventPayload.(*github.PushEvent)
  if !ok {
    t.Fatalf("expected PushEvent payload, got %T", EventPayload)
  }

  if got := pushPayload.GetRef(); got != "refs/heads/main" {
    t.Fatalf("expected ref refs/heads/main, got %q", got)
  }
}

func TestLoadActionContextInvalidPayload(t *testing.T) {
  defer func() { EventPayload = nil }()

  fs := afero.NewMemMapFs()

  const payloadPath = "/event.json"
  err := afero.WriteFile(fs, payloadPath, []byte("{not-json"), 0o644)
  if err != nil {
    t.Fatalf("failed to write payload file: %v", err)
  }

  env := map[string]string{
    "GITHUB_EVENT_NAME": "push",
    "GITHUB_EVENT_PATH": payloadPath,
  }

  loadActionContext(fs, mockLookup(env))

  if EventPayload != nil {
    t.Fatalf("expected nil payload for invalid JSON, got %T", EventPayload)
  }
}

func TestLoadActionContextMissingPayloadFile(t *testing.T) {
  defer func() { EventPayload = nil }()

  fs := afero.NewMemMapFs()

  env := map[string]string{
    "GITHUB_EVENT_NAME": "push",
    "GITHUB_EVENT_PATH": "/missing.json",
  }

  loadActionContext(fs, mockLookup(env))

  if EventPayload != nil {
    t.Fatalf("expected nil payload for missing file, got %T", EventPayload)
  }
}

func mockLookup(values map[string]string) lookupEnvFunc {
  return func(key string) (string, bool) {
    v, ok := values[key]
    return v, ok
  }
}

func strPtr(v string) *string { return &v }

func intPtr(v int) *int { return &v }

func boolPtr(v bool) *bool { return &v }

func TestActionContextString_NilContext(t *testing.T) {
  var ac *ActionContext
  got := ac.String()

  if got != "ActionContext: <nil>" {
    t.Fatalf("expected \"ActionContext: <nil>\", got %q", got)
  }
}

func TestActionContextString_EmptyContext(t *testing.T) {
  ac := &ActionContext{}
  got := ac.String()

  // Should contain the header
  if !strings.Contains(got, "ActionContext:") {
    t.Fatalf("expected output to contain \"ActionContext:\", got %q", got)
  }

  // Should contain field names with <nil> values
  expectedFields := []string{
    "GithubCI:",
    "Action:",
    "Actor:",
    "Repository:",
    "RunId:",
    "EventPayload:",
    "<nil>",
  }

  for _, field := range expectedFields {
    if !strings.Contains(got, field) {
      t.Fatalf("expected output to contain %q, got %q", field, got)
    }
  }
}

func TestActionContextString_PartiallyPopulated(t *testing.T) {
  ac := &ActionContext{
    Actor:      strPtr("octocat"),
    Repository: strPtr("owner/repo"),
    RunId:      intPtr(12345),
    Actions:    boolPtr(true),
    RunnerOs:   strPtr("Linux"),
  }

  got := ac.String()

  // Check that populated fields show their values
  if !strings.Contains(got, "Actor:") || !strings.Contains(got, "octocat") {
    t.Fatalf("expected Actor to show \"octocat\", got %q", got)
  }

  if !strings.Contains(got, "Repository:") || !strings.Contains(got, "owner/repo") {
    t.Fatalf("expected Repository to show \"owner/repo\", got %q", got)
  }

  if !strings.Contains(got, "RunId:") || !strings.Contains(got, "12345") {
    t.Fatalf("expected RunId to show \"12345\", got %q", got)
  }

  if !strings.Contains(got, "Actions:") || !strings.Contains(got, "true") {
    t.Fatalf("expected Actions to show \"true\", got %q", got)
  }

  if !strings.Contains(got, "RunnerOs:") || !strings.Contains(got, "Linux") {
    t.Fatalf("expected RunnerOs to show \"Linux\", got %q", got)
  }

  // Check that unpopulated fields show <nil>
  if !strings.Contains(got, "Action:") {
    t.Fatalf("expected Action field to be present, got %q", got)
  }

  // Should have <nil> values for unpopulated fields
  if !strings.Contains(got, "<nil>") {
    t.Fatalf("expected output to contain \"<nil>\" for unpopulated fields, got %q", got)
  }
}

func TestActionContextString_FullyPopulated(t *testing.T) {
  ac := &ActionContext{
    GithubCI:             boolPtr(true),
    Action:               strPtr("build"),
    ActionPath:           strPtr("/path/to/action"),
    ActionRepository:     strPtr("owner/action-repo"),
    Actions:              boolPtr(true),
    Actor:                strPtr("octocat"),
    ActorId:              strPtr("1234567"),
    ApiUrl:               strPtr("https://api.github.com"),
    BaseRef:              strPtr("main"),
    EnvFilePath:          strPtr("/tmp/env"),
    EventName:            strPtr("push"),
    EventPayloadFilePath: strPtr("/tmp/event.json"),
    GraphqlUrl:           strPtr("https://api.github.com/graphql"),
    HeadRef:              strPtr("feature-branch"),
    JobId:                strPtr("job-123"),
    OutputFilePath:       strPtr("/tmp/output"),
    PathFilePath:         strPtr("/tmp/path"),
    Ref:                  strPtr("refs/heads/main"),
    RefName:              strPtr("main"),
    RefProtected:         boolPtr(true),
    RefType:              strPtr("branch"),
    Repository:           strPtr("owner/repo"),
    RepositoryId:         strPtr("repo-123"),
    RepositoryOwner:      strPtr("owner"),
    RepositoryOwnerId:    strPtr("owner-123"),
    RetentionDays:        intPtr(90),
    RunAttempt:           intPtr(1),
    RunId:                intPtr(12345),
    RunNumber:            intPtr(42),
    ServerUrl:            strPtr("https://github.com"),
    Sha:                  strPtr("abc123def456"),
    StepSummaryFilePath:  strPtr("/tmp/summary"),
    TriggeringActor:      strPtr("octocat"),
    Workflow:             strPtr("CI"),
    WorkflowRef:          strPtr("owner/repo/.github/workflows/ci.yml@main"),
    WorkflowSha:          strPtr("workflow-sha"),
    WorkspacePath:        strPtr("/workspace"),
    RunnerArch:           strPtr("X64"),
    RunnerDebug:          boolPtr(false),
    RunnerEnvironment:    strPtr("github-hosted"),
    RunnerName:           strPtr("runner-1"),
    RunnerOs:             strPtr("Linux"),
    RunnerTempPath:       strPtr("/tmp"),
    RunnerToolCachePath:  strPtr("/tool-cache"),
    EventPayload:         &github.PushEvent{},
  }

  got := ac.String()

  // Check header
  if !strings.HasPrefix(got, "ActionContext:\n") {
    t.Fatalf("expected output to start with \"ActionContext:\\n\", got %q", got)
  }

  // Check that all string values are present
  expectedValues := []string{
    "build",
    "octocat",
    "owner/repo",
    "https://api.github.com",
    "refs/heads/main",
    "main",
    "branch",
    "abc123def456",
    "CI",
    "Linux",
    "X64",
    "github-hosted",
  }

  for _, val := range expectedValues {
    if !strings.Contains(got, val) {
      t.Fatalf("expected output to contain %q, got %q", val, got)
    }
  }

  // Check that integer values are present
  if !strings.Contains(got, "12345") {
    t.Fatalf("expected RunId 12345, got %q", got)
  }
  if !strings.Contains(got, "42") {
    t.Fatalf("expected RunNumber 42, got %q", got)
  }
  if !strings.Contains(got, "90") {
    t.Fatalf("expected RetentionDays 90, got %q", got)
  }

  // Check that boolean values are present
  if !strings.Contains(got, "true") {
    t.Fatalf("expected bool value true, got %q", got)
  }
  if !strings.Contains(got, "false") {
    t.Fatalf("expected bool value false, got %q", got)
  }

  // Check that EventPayload type is shown
  if !strings.Contains(got, "github.PushEvent") {
    t.Fatalf("expected EventPayload type *github.PushEvent, got %q", got)
  }
}

func TestActionContextString_WithNilEventPayload(t *testing.T) {
  ac := &ActionContext{
    Actor:        strPtr("octocat"),
    EventPayload: nil,
  }

  got := ac.String()

  // Check that EventPayload is shown as <nil>
  if !strings.Contains(got, "EventPayload:") {
    t.Fatalf("expected output to contain \"EventPayload:\", got %q", got)
  }

  // Should show <nil> for EventPayload
  lines := strings.Split(got, "\n")
  foundEventPayloadNil := false
  for _, line := range lines {
    if strings.Contains(line, "EventPayload:") && strings.Contains(line, "<nil>") {
      foundEventPayloadNil = true
      break
    }
  }

  if !foundEventPayloadNil {
    t.Fatalf("expected EventPayload to show <nil>, got %q", got)
  }
}
