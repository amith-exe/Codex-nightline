package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var activateStartDaemon = startActivateDaemon

type harnessList []string

func (h *harnessList) String() string {
	return strings.Join(*h, ",")
}

func (h *harnessList) Set(value string) error {
	for _, part := range strings.Split(value, ",") {
		name := strings.TrimSpace(part)
		if name == "" {
			continue
		}
		switch normalizeHarnessName(name) {
		case "claude-code", "codex":
			*h = append(*h, normalizeHarnessName(name))
		default:
			return fmt.Errorf("unknown harness %q", name)
		}
	}
	return nil
}

func activateCommand(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("activate", flag.ContinueOnError)
	flags.SetOutput(stderr)
	scope := flags.String("scope", "project", "hook scope: project or user")
	quiet := flags.Bool("quiet", false, "suppress non-error output")
	var harnesses harnessList
	flags.Var(&harnesses, "harness", "harness to activate; repeat or comma-separate (default claude-code,codex)")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintf(stderr, "flowsync activate: unexpected argument %q\n", flags.Arg(0))
		return 2
	}
	if len(harnesses) == 0 {
		harnesses = harnessList{"claude-code", "codex"}
	}
	if *scope != "project" && *scope != "user" {
		fmt.Fprintf(stderr, "flowsync activate: unknown scope %q\n", *scope)
		return 2
	}

	// If activating a project, create a local .flowsync directory so
	// journal and checkpoint storage can be placed inside the project.
	if *scope == "project" {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(stderr, "flowsync activate: get working dir: %v\n", err)
			return 1
		}
		local := filepath.Join(wd, ".flowsync")
		if err := os.MkdirAll(local, 0o700); err != nil {
			fmt.Fprintf(stderr, "flowsync activate: create %s: %v\n", local, err)
			return 1
		}
		if err := os.Chmod(local, 0o700); err != nil {
			fmt.Fprintf(stderr, "flowsync activate: chmod %s: %v\n", local, err)
			return 1
		}
		localJournal := filepath.Join(wd, "flowsync")
		if err := os.MkdirAll(localJournal, 0o700); err != nil {
			fmt.Fprintf(stderr, "flowsync activate: create %s: %v\n", localJournal, err)
			return 1
		}
		if err := os.Chmod(localJournal, 0o700); err != nil {
			fmt.Fprintf(stderr, "flowsync activate: chmod %s: %v\n", localJournal, err)
			return 1
		}
	}

	for _, harness := range harnesses {
		code := runActivateInit(harness, *scope, *quiet, stdout, stderr)
		if code != 0 {
			return code
		}
	}
	if code := activateStartDaemon(*quiet, stdout, stderr); code != 0 {
		return code
	}
	if !*quiet {
		fmt.Fprintln(stdout, "flowsync: project activation complete")
	}
	return 0
}

func startActivateDaemon(quiet bool, stdout, stderr io.Writer) int {
	root, socketPath, err := daemonPaths("", "")
	if err != nil {
		fmt.Fprintf(stderr, "flowsync activate: %v\n", err)
		return 1
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := ensureDaemon(ctx, root, socketPath, ""); err != nil {
		fmt.Fprintf(stderr, "flowsync activate: start daemon: %v\n", err)
		return 1
	}
	if !quiet {
		fmt.Fprintln(stdout, "flowsync: daemon ready")
	}
	return 0
}

func runActivateInit(harness, scope string, quiet bool, stdout, stderr io.Writer) int {
	initArgs := []string{harness, "--scope", scope}
	if !quiet {
		return initCommand(initArgs, stdout, stderr)
	}

	var initStdout, initStderr bytes.Buffer
	code := initCommand(initArgs, &initStdout, &initStderr)
	if code != 0 {
		_, _ = io.Copy(stderr, &initStderr)
		_, _ = io.Copy(stdout, &initStdout)
	}
	return code
}

func normalizeHarnessName(name string) string {
	switch name {
	case "claudecode":
		return "claude-code"
	default:
		return name
	}
}
