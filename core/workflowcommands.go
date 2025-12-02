package core

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/google/uuid"
)

// WorkflowCommandString is the delimiter that wraps GitHub Actions workflow commands.
// https://docs.github.com/actions/using-workflows/workflow-commands-for-github-actions#about-workflow-commands
const WorkflowCommandString = "::"

const (
	DebugCommandStr        = "debug"
	NoticeCommandStr       = "notice"
	WarningCommandStr      = "warning"
	ErrorCommandStr        = "error"
	StartGroupCommandStr   = "group"
	EndGroupCommandStr     = "endgroup"
	MaskCommandStr         = "add-mask"
	StopCommandsCommandStr = "stop-commands"
)

// WorkflowCommand represents a formatted GitHub Actions workflow command.
// GitHub Workflow Command Reference: https://docs.github.com/actions/using-workflows/workflow-commands-for-github-actions#about-workflow-commands
type WorkflowCommand struct {
	w          io.Writer
	command    string
	properties map[string]string
	message    string
}

func NewWorkflowCommand(w io.Writer, command string, properties map[string]string, message string) *WorkflowCommand {
	return &WorkflowCommand{
		w:          w,
		command:    command,
		properties: properties,
		message:    message,
	}
}

// String renders the workflow command into the GitHub Actions command syntax.
// GitHub Workflow Command Reference: https://docs.github.com/actions/using-workflows/workflow-commands-for-github-actions#about-workflow-commands
func (wc *WorkflowCommand) String() string {
	var cmdStr strings.Builder
	fmt.Fprintf(&cmdStr, "%s%s", WorkflowCommandString, wc.command)

	if wc.properties != nil && len(wc.properties) > 0 {
		cmdStr.WriteString(" ")
		isFirst := true
		for k, v := range wc.properties {
			if !isFirst {
				cmdStr.WriteString(",")
			}

			fmt.Fprintf(&cmdStr, "%s=%s", k, escapeCommandParameter(v))

			isFirst = false
		}
	}

	fmt.Fprintf(&cmdStr, "%s%s", WorkflowCommandString, escapeCommandData(wc.message))

	return cmdStr.String()
}

// Send writes the workflow command to the provided writer.
// GitHub Workflow Command Reference: https://docs.github.com/actions/using-workflows/workflow-commands-for-github-actions#about-workflow-commands
func (wc *WorkflowCommand) Send() error {
	_, err := fmt.Fprintln(wc.w, wc)

	return err
}

// Debug writes a debug message to stdout in workflow-command format.
// GitHub Debug Command Reference: https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-commands#setting-a-debug-message
func Debug(message string) error {
	cmd := NewWorkflowCommand(os.Stdout, DebugCommandStr, nil, message)

	return cmd.Send()
}

// WorkflowCommandOptions describes optional metadata for notice/warning/error commands.
type WorkflowCommandOptions struct {
	FileName    *string
	StartLine   *int
	EndLine     *int
	Title       *string
	StartColumn *int
	EndColumn   *int
}

func (wco *WorkflowCommandOptions) CmdProperties() map[string]string {
	if wco == nil {
		return nil
	}

	params := make(map[string]string)

	if wco.FileName != nil && *wco.FileName != "" {
		params["file"] = *wco.FileName
	}
	if wco.StartLine != nil {
		params["line"] = fmt.Sprintf("%d", *wco.StartLine)
	}
	if wco.EndLine != nil {
		params["endLine"] = fmt.Sprintf("%d", *wco.EndLine)
	}
	if wco.Title != nil {
		params["title"] = *wco.Title
	}
	if wco.StartColumn != nil {
		params["col"] = fmt.Sprintf("%d", *wco.StartColumn)
	}
	if wco.EndColumn != nil {
		params["endColumn"] = fmt.Sprintf("%d", *wco.EndColumn)
	}

	return params
}

// Notice writes a notice message to the supplied writer in workflow-command format.
// GitHub Notice Command Reference: https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-commands#setting-a-notice-message
func Notice(message string, options *WorkflowCommandOptions) error {
	cmd := NewWorkflowCommand(os.Stdout, NoticeCommandStr, options.CmdProperties(), message)

	return cmd.Send()
}

// Warning writes a warning message to stdout in workflow-command format.
// GitHub Warning Command Reference: https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-commands#setting-a-warning-message
func Warning(message string, options *WorkflowCommandOptions) error {
	cmd := NewWorkflowCommand(os.Stdout, WarningCommandStr, options.CmdProperties(), message)

	return cmd.Send()
}

// Error writes an error message to stdout in workflow-command format.
// GitHub Error Command Reference: https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-commands#setting-an-error-message
func Error(message string, options *WorkflowCommandOptions) error {
	cmd := NewWorkflowCommand(os.Stdout, ErrorCommandStr, options.CmdProperties(), message)

	return cmd.Send()
}

// StartGroup begins a collapsible log group.
// GitHub Group Command Reference: https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-commands#grouping-log-lines
func StartGroup(name string) error {
	cmd := NewWorkflowCommand(os.Stdout, StartGroupCommandStr, nil, name)

	return cmd.Send()
}

// EndGroup closes the most recent collapsible log group.
// GitHub Group Command Reference: https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-commands#grouping-log-lines
func EndGroup() error {
	cmd := NewWorkflowCommand(os.Stdout, EndGroupCommandStr, nil, "")

	return cmd.Send()
}

// MaskValue instructs GitHub Actions to redact the provided value from logs.
// GitHub Mask Command Reference: https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-commands#masking-a-value-in-a-log
func MaskValue(value string) error {
	cmd := NewWorkflowCommand(os.Stdout, MaskCommandStr, nil, value)

	return cmd.Send()
}

type ResumeCommandsFunc func() error

// StopCommands suspends workflow command processing until ResumeCommands is called.
// GitHub Stop Commands Reference: https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-commands#stopping-and-starting-workflow-commands
func StopCommands() (ResumeCommandsFunc, error) {
	delimUuid, err := uuid.NewV7()

	if err != nil {
		return nil, err
	}

	cmd := NewWorkflowCommand(os.Stdout, StopCommandsCommandStr, nil, delimUuid.String())

	err = cmd.Send()

	if err != nil {
		return nil, err
	}

	resumeCommands := func() error {
		resumeCmd := NewWorkflowCommand(os.Stdout, delimUuid.String(), nil, "")

		return resumeCmd.Send()
	}

	return resumeCommands, nil
}

// escapeCommandData escapes workflow command data according to GitHub Actions rules.
func escapeCommandData(s string) string {
	if s == "" {
		return s
	}

	s = strings.ReplaceAll(s, "%", "%25")
	s = strings.ReplaceAll(s, "\r", "%0D")
	s = strings.ReplaceAll(s, "\n", "%0A")
	return s
}

// escapeCommandParameter escapes workflow command parameters according to GitHub Actions rules.
func escapeCommandParameter(s string) string {
	if s == "" {
		return s
	}

	s = strings.ReplaceAll(s, "%", "%25")
	s = strings.ReplaceAll(s, "\r", "%0D")
	s = strings.ReplaceAll(s, "\n", "%0A")
	s = strings.ReplaceAll(s, ":", "%3A")
	s = strings.ReplaceAll(s, ",", "%2C")
	return s
}
