package reflector

import (
	"bytes"
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"flowsync/internal/core"
)

const (
	DefaultCommand     = "ollama"
	DefaultModel       = "qwen2.5-coder:7b"
	ReflectorActiveEnv = "THREADMARK_REFLECTOR_ACTIVE"
)

const disableHooksSettings = `{"disableAllHooks":true}`

//go:embed default_prompt.md
var embeddedDefaultPrompt string

type Mode string

const (
	ModeConvenience Mode = "convenience"
	ModeBare        Mode = "bare"
)

func ParseMode(value string) (Mode, error) {
	switch Mode(strings.ToLower(strings.TrimSpace(value))) {
	case "", ModeConvenience:
		return ModeConvenience, nil
	case ModeBare:
		return ModeBare, nil
	default:
		return "", fmt.Errorf("unknown reflector mode %q", value)
	}
}

type Config struct {
	Command    string
	Mode       Mode
	Model      string
	PromptPath string
	Runner     Runner
}

type Runner interface {
	Run(context.Context, string, []string, io.Reader) (string, error)
}

type ExecRunner struct{}

func (ExecRunner) Run(ctx context.Context, command string, args []string, stdin io.Reader) (string, error) {
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Env = withEnv(os.Environ(), ReflectorActiveEnv, "1")
	cmd.Stdin = stdin
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		redactedStderr := truncate(Redact(stderr.String()), 2048)
		return "", fmt.Errorf("run reflector command: %w: %s", err, redactedStderr)
	}
	return stdout.String(), nil
}

func withEnv(base []string, key, value string) []string {
	prefix := key + "="
	out := make([]string, 0, len(base)+1)
	replaced := false
	for _, entry := range base {
		if strings.HasPrefix(entry, prefix) {
			if !replaced {
				out = append(out, prefix+value)
				replaced = true
			}
			continue
		}
		out = append(out, entry)
	}
	if !replaced {
		out = append(out, prefix+value)
	}
	return out
}

type Client struct {
	config Config
}

type Request struct {
	Checkpoint core.CheckpointRequest
	Transcript string
}

type Response struct {
	Body        string
	Backend     string
	CommandMode Mode
	Model       string
	PromptHash  string
	StartedAt   time.Time
	CompletedAt time.Time
}

func NewClient(config Config) (*Client, error) {
	mode, err := ParseMode(string(config.Mode))
	if err != nil {
		return nil, err
	}
	config.Mode = mode
	if strings.TrimSpace(config.Command) == "" {
		config.Command = DefaultCommand
	}
	if strings.TrimSpace(config.Model) == "" {
		config.Model = DefaultModel
	}
	if config.Runner == nil {
		config.Runner = ExecRunner{}
	}
	return &Client{config: config}, nil
}

func DefaultPromptPath() string {
	return filepath.Join("prompts", "reflector.md")
}

func (c *Client) Reflect(ctx context.Context, request Request) (Response, error) {
	promptPath := c.config.PromptPath
	if strings.TrimSpace(promptPath) == "" {
		promptPath = DefaultPromptPath()
	}

	systemPrompt, err := LoadPromptFile(promptPath)
	if err != nil {
		return Response{}, err
	}
	promptHash := PromptHash(systemPrompt)
	prompt := BuildPrompt(systemPrompt, request)

	startedAt := time.Now().UTC()
	var body string
	var runErr error
	// If the configured command is `ollama`, pass the prompt via stdin
	// and run: `ollama run <model>` so local Ollama models can be used as the reflector backend.
	if strings.TrimSpace(filepath.Base(c.config.Command)) == "ollama" || strings.TrimSpace(c.config.Command) == "ollama" {
		args := []string{"run", c.config.Model}
		stdin := strings.NewReader(prompt)
		body, runErr = c.config.Runner.Run(ctx, c.config.Command, args, stdin)
	} else {
		args := c.args(prompt)
		body, runErr = c.config.Runner.Run(ctx, c.config.Command, args, nil)
	}
	if runErr != nil {
		return Response{}, runErr
	}
	body = strings.TrimSpace(body)
	if body == "" {
		return Response{}, errors.New("reflector returned empty body")
	}

	backend := c.config.Command
	return Response{
		Body:        body,
		Backend:     backend,
		CommandMode: c.config.Mode,
		Model:       c.config.Model,
		PromptHash:  promptHash,
		StartedAt:   startedAt,
		CompletedAt: time.Now().UTC(),
	}, nil
}

func (c *Client) args(prompt string) []string {
	if c.config.Mode == ModeBare {
		return []string{"--bare", "-p", prompt}
	}
	return []string{"--settings", disableHooksSettings, "-p", prompt}
}

func LoadPromptFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if isDefaultPromptPath(path) {
			return ExtractPromptMarkdown(embeddedDefaultPrompt), nil
		}
		return "", fmt.Errorf("read reflector prompt: %w", err)
	}
	return ExtractPromptMarkdown(string(data)), nil
}

func isDefaultPromptPath(path string) bool {
	return filepath.Clean(path) == filepath.Clean(DefaultPromptPath())
}

func ExtractPromptMarkdown(markdown string) string {
	parts := strings.Split(markdown, "```")
	if len(parts) >= 3 {
		return strings.TrimSpace(parts[1])
	}
	return strings.TrimSpace(markdown)
}

func PromptHash(prompt string) string {
	sum := sha256.Sum256([]byte(prompt))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func BuildPrompt(systemPrompt string, request Request) string {
	checkpoint := request.Checkpoint
	transcript := Redact(request.Transcript)
	if strings.TrimSpace(transcript) == "" {
		transcript = "No in-memory event excerpt was available for this checkpoint."
	}

	var b strings.Builder
	b.WriteString("<!-- flowsync:cacheable-system-prompt ")
	b.WriteString(PromptHash(systemPrompt))
	b.WriteString(" -->\n")
	b.WriteString(strings.TrimSpace(systemPrompt))
	b.WriteString("\n<!-- /flowsync:cacheable-system-prompt -->\n\n")
	b.WriteString("FlowSync checkpoint context:\n")
	b.WriteString("- trigger: ")
	b.WriteString(checkpoint.Trigger.Reason)
	b.WriteByte('\n')
	b.WriteString("- thread_id: ")
	b.WriteString(checkpoint.ThreadID)
	b.WriteByte('\n')
	b.WriteString("- transcript_id: ")
	b.WriteString(checkpoint.TranscriptID)
	b.WriteByte('\n')
	b.WriteString("- harness: ")
	b.WriteString(checkpoint.Harness)
	b.WriteByte('\n')
	b.WriteString("- working_dir: ")
	b.WriteString(checkpoint.WorkingDir)
	b.WriteByte('\n')
	b.WriteString("- source_event_range: [")
	b.WriteString(checkpoint.SourceRange.First)
	b.WriteString(", ")
	b.WriteString(checkpoint.SourceRange.Last)
	b.WriteString("]\n\n")
	b.WriteString("The following event excerpt is redacted and may be incomplete. Use it as the source material for the journal body, not as a file-by-file diff.\n\n")
	b.WriteString("```flowsync-events\n")
	b.WriteString(transcript)
	if !strings.HasSuffix(transcript, "\n") {
		b.WriteByte('\n')
	}
	b.WriteString("```\n")
	return b.String()
}

func truncate(value string, max int) string {
	value = strings.TrimSpace(value)
	if len(value) <= max {
		return value
	}
	return value[:max] + "...<truncated>"
}
