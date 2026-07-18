// Command flowsync is the user-facing CLI for inspecting and controlling
// FlowSync. Available subcommands:
//
//	init        install hook shims into the active harness's settings
//	install     install flowsync and flowsyncd onto PATH
//	activate    install default project hooks for supported harnesses
//	status      print daemon status and recent journal entries
//	journal     show / tail / search the current project's journal
//	enable      remove a project opt-out marker
//	disable     create a project opt-out marker
//	purge       remove local FlowSync storage
//
// See HOW_IT_WORKS.md for the full architecture.
package main

import (
	"flag"
	"fmt"
	"os"

	"flowsync/internal/buildinfo"
	tmjournal "flowsync/internal/journal"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "-h", "--help", "help":
		usage()
	case "-v", "--version", "version":
		fmt.Fprintln(os.Stdout, buildinfo.Version)
	case "activate":
		os.Exit(activateCommand(os.Args[2:], os.Stdout, os.Stderr))
	case "init":
		os.Exit(initCommand(os.Args[2:], os.Stdout, os.Stderr))
	case "install":
		os.Exit(installCommand(os.Args[2:], os.Stdout, os.Stderr))
	case "daemon":
		os.Exit(daemonCommand(os.Args[2:], os.Stdout, os.Stderr))
	case "doctor":
		os.Exit(doctorCommand(os.Args[2:], os.Stdout, os.Stderr))
	case "disable":
		os.Exit(disableCommand(os.Args[2:], os.Stdout, os.Stderr))
	case "enable":
		os.Exit(enableCommand(os.Args[2:], os.Stdout, os.Stderr))
	case "hook":
		os.Exit(hook(os.Args[2:], os.Stdin, os.Stdout, os.Stderr))
	case "purge":
		os.Exit(purgeCommand(os.Args[2:], os.Stdin, os.Stdout, os.Stderr))
	case "status":
		os.Exit(statusCommand(os.Args[2:], os.Stdout, os.Stderr))
	case "journal":
		journal(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "flowsync: unknown command %q\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
}

func journal(args []string) {
	flags := flag.NewFlagSet("journal", flag.ExitOnError)
	workingDir := flags.String("working-dir", ".", "project working directory")
	root := flags.String("root", "", "flowsync root directory (default ~/.flowsync)")
	last := flags.Int("last", 3, "number of recent entries to print")
	flags.Parse(args)

	store := tmjournal.NewStore(*root)
	entries, err := store.LastEntries(*workingDir, *last)
	if err != nil {
		fatal(err)
	}
	for i, entry := range entries {
		if i > 0 {
			fmt.Println()
		}
		fmt.Print(entry)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: flowsync <command> [options]")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(os.Stderr, "commands:")
	fmt.Fprintln(os.Stderr, "  activate  activate FlowSync for the current project")
	fmt.Fprintln(os.Stderr, "  init      install FlowSync adapter configuration")
	fmt.Fprintln(os.Stderr, "  install   install flowsync and flowsyncd onto PATH")
	fmt.Fprintln(os.Stderr, "  daemon    manage the local FlowSync daemon")
	fmt.Fprintln(os.Stderr, "  doctor    diagnose local FlowSync setup")
	fmt.Fprintln(os.Stderr, "  disable   disable FlowSync for a project")
	fmt.Fprintln(os.Stderr, "  enable    re-enable FlowSync for a project")
	fmt.Fprintln(os.Stderr, "  hook      internal hook bridge used by adapters")
	fmt.Fprintln(os.Stderr, "  purge     remove local FlowSync storage")
	fmt.Fprintln(os.Stderr, "  status    print local FlowSync paths")
	fmt.Fprintln(os.Stderr, "  version   print FlowSync version")
	fmt.Fprintln(os.Stderr, "  journal   print recent journal entries for a project")
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "flowsync: %v\n", err)
	os.Exit(1)
}
