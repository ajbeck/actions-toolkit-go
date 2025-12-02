package actioncontext

import (
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
